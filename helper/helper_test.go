package helper

import (
	"fmt"
	"testing"
)

func TestGetExternalIP(t *testing.T) {
	a, err := GetExternalIP()
	fmt.Printf("aaaaaa, %+v, err: %+v\n", a, err)
}
