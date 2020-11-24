package curl

import (
	"context"
	"testing"
)

func TestParseCommand(t *testing.T) {
	cmd, err := ParseCommand(in)
	if err != nil {
		t.Fatalf("should not return an error, but got '%+v'", err)
	}
	if cmd.url == nil {
		t.Error("url should not be empty")
	}

	cmd.Request(context.Background())
}
