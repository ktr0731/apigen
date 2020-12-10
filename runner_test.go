package apigen_test

import (
	"bytes"
	"context"
	"flag"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ktr0731/apigen"
	"github.com/ktr0731/apigen/curl"
)

var update = flag.Bool("update", false, "update golden files")

func TestGenerate(t *testing.T) {
	t.Parallel()

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
					Name:      "ListComments",
					Request:   curl.ParseCommand(`curl https://jsonplaceholder.typicode.com/posts/1/comments`),
					ParamHint: "/posts/{postID}/comments",
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

	var w bytes.Buffer
	if err := apigen.Generate(context.Background(), def, apigen.WithWriter(&w)); err != nil {
		t.Fatalf("should not return an error, but got '%s'", err)
	}

	assertWithGolden(t, w.String())
}

func assertWithGolden(t *testing.T, actual string) {
	t.Helper()

	name := t.Name()
	r := strings.NewReplacer(
		"/", "-",
		" ", "_",
		"=", "-",
		"'", "",
		`"`, "",
		",", "",
	)
	normalizeFilename := func(name string) string {
		fname := r.Replace(strings.ToLower(name)) + ".golden"

		return filepath.Join("testdata", fname)
	}

	fname := normalizeFilename(name)

	if *update {
		if err := ioutil.WriteFile(fname, []byte(actual), 0600); err != nil {
			t.Fatalf("failed to update the golden file: %s", err)
		}

		return
	}

	// Load the golden file.
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("failed to load a golden file: %s", err)
	}

	expected := string(b)
	if runtime.GOOS == "windows" {
		expected = strings.ReplaceAll(expected, "\r\n", "\n")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
}
