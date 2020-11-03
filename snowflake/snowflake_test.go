package snowflake

import (
	"fmt"
	"testing"
)

func TestSnowflake(t *testing.T) {
	node, err := NewNode(1020)
	if err != nil {
		t.Fatal(err)
	}
	id := node.Generate()
	fmt.Printf("generate id success:%d\n", id.Int64())
}
