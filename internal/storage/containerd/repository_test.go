package containerd

import (
	"context"
	"testing"

	"github.com/distribution/reference"
)

func TestContainerdImageName(t *testing.T) {
	name, err := reference.WithName("something")
	if err != nil {
		t.Fatalf("create reference: %v", err)
	}

	t.Run("preserves request host", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "http.request.host", "localhost:5051")

		got, err := containerdImageName(ctx, name)
		if err != nil {
			t.Fatalf("containerdImageName: %v", err)
		}

		if got.String() != "localhost:5051/something" {
			t.Fatalf("got %q, want %q", got.String(), "localhost:5051/something")
		}
	})

	t.Run("keeps name when host is missing", func(t *testing.T) {
		got, err := containerdImageName(context.Background(), name)
		if err != nil {
			t.Fatalf("containerdImageName: %v", err)
		}

		if got.String() != "something" {
			t.Fatalf("got %q, want %q", got.String(), "something")
		}
	})
}
