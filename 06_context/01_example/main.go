// Minimal example where create a context and add value inside and retreive
package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	ctx = AddValue(ctx)
	ReadValue(ctx)
}

// This method add a key inside context
func AddValue(ctx context.Context) context.Context {
	return context.WithValue(ctx, "key", "test-value")
}

// Read value from one key inside context
func ReadValue(ctx context.Context) {
	val := ctx.Value("key")
	fmt.Println(val)
}
