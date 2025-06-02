package utils

import (
	"fmt"
	"testing"
)

func TestChatMessage(t *testing.T) {
	message, err := ChatMessage("你是谁", 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(message)
}
