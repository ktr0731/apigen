# apigen

[![PkgGoDev](https://pkg.go.dev/badge/github.com/ktr0731/apigen)](https://pkg.go.dev/github.com/ktr0731/apigen)
[![GitHub Actions](https://github.com/ktr0731/apigen/workflows/main/badge.svg)](https://github.com/ktr0731/apigen/actions)  

`apigen` generates API client via execution environment such as `curl`.

## Installation
``` bash
$ go get github.com/ktr0731/apigen
```

## Usage

**This example is located under [here](./_example)**  

### Generator
`apigen` requires `*Definition` which describes methods the service has.  
Following definition defines `CreatePost`, `ListPosts`, `GetPost`, `UpdatePost` and `DeletePost` which belong to `Dummy` service.
`Request` specify execution environment, `apigen` generates the API client and request/response types based on the execution result.
`ParamHint` is only required when its method uses path parameters such as `"/post/{postID}"`. `apigen` generates the request type by using it.

The artifact will be written to `client_gen.go` which is specified by `apigen.WithWriter`. Default output is stdout.

``` go
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
```

The artifact is [here](./_example/client_gen.go).  

### Client

We can invoke the API server using the generated API client.  

``` go
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
```

The output is:

```json
[
  {
    "body": "quo et expedita modi cum officia vel magni\ndoloribus qui repudiandae\nvero nisi sit\nquos veniam quod sed accusamus veritatis error",
    "id": 10,
    "title": "optio molestias id quia eum",
    "userId": 1
  }
]
```
