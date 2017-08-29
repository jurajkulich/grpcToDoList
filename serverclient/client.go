package main

import (
	"google.golang.org/grpc"
	"log"
	pb "github.com/yuro8/grpctodolist/proto"
	"golang.org/x/net/context"

)

const (
	address = "localhost:50050"
)
func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewToDoActionsClient(conn)

	r, err := c.AddItem(context.Background(), &pb.ToDoItem{Name:"Dog", Description:"eating", IsDone:false})
	if err != nil {
		log.Fatalf("could not add item: %v", err)
	}
	log.Printf("Added item: %v", r)

	list, err := c.GetItems(context.Background(), &pb.Nothing{})
	if err != nil {
		log.Fatalf("could get items: %v", err)
	}
	log.Printf("ToDoItems: %v\n", list)

	_, err = c.DeleteItem(context.Background(), &pb.ID{5})
	if err != nil {
		log.Fatalf(err.Error())
	}

	item, err := c.GetItem(context.Background(), &pb.ID{7})
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("ToDoItem: %v\n", item)

	updated := pb.ToDoItem{Name:"Mirinda", Description:"let her go", IsDone:true}
	_, err = c.UpdateItem(context.Background(), &pb.UpdatedItem{&updated, &pb.ID{2}})
	if err != nil {
		log.Fatalf("%v", err)
	}
	list, err = c.GetItems(context.Background(), &pb.Nothing{})
	if err != nil {
		log.Fatalf("could get items: %v", err)
	}
	log.Printf("ToDoItems: %v\n", list)


}

// /home/juraj/go/src/github.com/yuro8/grpctodolist/serverclient
