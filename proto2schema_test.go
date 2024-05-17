package proto2schema

import (
	"fmt"
	"os"
	"testing"
)

func TestProto2schema(t *testing.T) {
	str := Proto2schema("./test.proto")
	b, _ := os.ReadFile("./test.k")

	fmt.Println(string(b) == str)
}
