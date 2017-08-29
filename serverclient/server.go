package main

import (
	pb "github.com/yuro8/grpctodolist/proto"
	"net"
	"github.com/labstack/gommon/log"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"golang.org/x/net/context"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

const (
	port = ":50050"
)

type server struct {
	db *gorm.DB
}

func NewToDoServer(db *gorm.DB) *server {
	return &server{db:db}
}


func (s *server) GetItems(context context.Context, in *pb.Nothing) (*pb.ToDoList, error) {
	list := []*pb.ToDoItem{}
	s.db.Find(&list)
	return &pb.ToDoList{ToDoList: list}, nil
}

func ( s *server) AddItem(context context.Context, item *pb.ToDoItem) (*pb.ToDoItem,error) {
	if item.Name == "" {
		return &pb.ToDoItem{}, nil
	}
	s.db.Create(&item)
	return item, nil
}

func ( s *server) DeleteItem(context context.Context, id *pb.ID) (*pb.Nothing, error) {
	if s.db.Delete(&pb.ToDoItem{}, "id = ?", id.Id).RecordNotFound() {
		return &pb.Nothing{}, errors.New("Could not find item")
	}
	return &pb.Nothing{}, nil
}


func ( s *server) GetItem(context context.Context, id *pb.ID) (*pb.ToDoItem, error) {
	todo := pb.ToDoItem{}
	if s.db.Where("id = ?", id.Id).First(&todo).RecordNotFound() {
		return &pb.ToDoItem{}, errors.New("Item not found")
	}
	return &todo, nil
}

func ( s *server) UpdateItem(context context.Context, item *pb.UpdatedItem) (*pb.Nothing, error) {
	todo := pb.ToDoItem{}
	if s.db.Where("id = ?", item.UpdateID.Id).First(&todo).RecordNotFound() {
		return &pb.Nothing{}, errors.New("Item not found")
	} else {
		todo = *item.UpdateTodo
		todo.Id = item.UpdateID.Id
		err := s.db.Save(todo)
		if err != nil {
			return &pb.Nothing{}, errors.New("Could not update item")
		}
	}
	return &pb.Nothing{}, nil
}

func main() {
	database, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		panic("Can't connect to database: ")
	}
	database.AutoMigrate(&pb.ToDoItem{})
	// database.DropTable(&pb.ToDoItem{})
	defer database.Close()

	server := NewToDoServer(database)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		//log.Fatalf("Cannot listen to server")
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterToDoActionsServer(s, server)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
