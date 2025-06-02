package redis

import (
	"context"
	"testing"
)

func TestGetRedis(t *testing.T) {
	if err := GetRedis().Set(context.Background(), "test_key", "test_value", 0).Err(); err != nil {
		t.Fatal(err)
	}
}
