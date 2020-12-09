package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	client := NewDummyClient()

	res, err := client.GetPost(context.Background(), &GetPostRequest{ID: "10"})
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(&res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
