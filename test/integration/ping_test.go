package integration

import (
	"github.com/ppal31/gochat/cli/client"
	"testing"
)

const MAX = 1000

func TestChatService_Ping(t *testing.T) {
	c := client.NewClient([]string{"127.0.0.1:2181"})
	err := c.Ping()
	if err != nil {
		t.Errorf("Failed to ping %v", err.Error())
	}
}
