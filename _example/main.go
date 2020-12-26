package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ktr0731/apigen/client"
)

func main() {
	client := NewDummyClient(client.WithInterceptors(client.ConvertStatusCodeToErrorInterceptor()))

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
