package proto2schema

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProto2schema(t *testing.T) {
	str := Proto2schema("./test.proto")
	b, _ := os.ReadFile("./test.k")

	assert.Equal(t, string(b), str)
}
