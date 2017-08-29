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
	r, err := c.AddItem(context.Background(), &pb.ToDoItem{Name:"Pizza", Description:"eating"})
	if err != nil {
		log.Fatalf("could add item: %v", err)
	}
	log.Printf("Added item: %v", r)
	list, err := c.GetItems(context.Background(), &pb.Nothing{})
	if err != nil {
		log.Fatalf("could get items: %v", err)
	}
	log.Printf("ToDoItems: %v", list)
}
