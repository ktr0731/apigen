package main

import (
	"context"
	"log"
	"os"

	"github.com/ktr0731/apigen"
	"github.com/ktr0731/apigen/curl"
)

func main() {
	def := &apigen.Definition{
		Services: map[string][]*apigen.Method{
			"Dummy": {
				{
					Name:    "CreatePost",
					Request: curl.ParseCommand(`curl 'https://jsonplaceholder.typicode.com/posts' --data-binary '{"title":"foo","body":"bar","userId":1}'`),
				},
				{
					Name:    "ListPosts",
					Request: curl.ParseCommand(`curl https://jsonplaceholder.typicode.com/posts`),
				},
				{
					Name:    "GetPost",
					Request: curl.ParseCommand(`curl https://jsonplaceholder.typicode.com/posts?id=1`),
				},
				{
					Name:      "UpdatePost",
					Request:   curl.ParseCommand(`curl 'https://jsonplaceholder.typicode.com/posts/1' -X 'PUT' --data-binary '{"title":"foo","body":"bar","userId":1}'`),
					ParamHint: "/posts/{postID}",
				},
				{
					Name:      "DeletePost",
					Request:   curl.ParseCommand(`curl 'https://jsonplaceholder.typicode.com/posts/1' -X 'DELETE'`),
					ParamHint: "/posts/{postID}",
				},
			},
		},
	}

	f, err := os.Create("client_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := apigen.Generate(context.Background(), def, apigen.WithWriter(f)); err != nil {
		log.Fatal(err)
	}
}
