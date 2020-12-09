package main

import (
	"context"
	"log"

	"github.com/k0kubun/pp"
)

func main() {
	client := NewDummyClient()

	res, err := client.GetPost(context.Background(), &GetPostRequest{ID: "10"})
	if err != nil {
		log.Fatal(err)
	}

	pp.Printf("%+v", res)
}
