package confutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//execute gotest by "go test -args -c=any -test.run=GoAny"
func TestGoAny(t *testing.T) {
	assert.Equal(t, strings.Join(GetConfs("Redis", "redis"), ","), "127.0.0.1:6379,127.0.0.1:7379")
}
