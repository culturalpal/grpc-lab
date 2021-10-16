package lb

import (
	"github.com/ppal31/grpc-lab/cli/lb/client"
	"testing"
)

func TestChatService_Ping(t *testing.T) {
	c := client.NewClient([]string{"127.0.0.1:2181"}, []string{"5001", "5002"})
	err := c.Ping()
	if err != nil {
		t.Errorf("Failed to ping %v", err.Error())
	}
}

func TestChatService_Chat(t *testing.T) {
	c := client.NewClient([]string{"127.0.0.1:2181"}, []string{"5001", "5002"})
	err := c.Chat("5001")
	if err != nil {
		t.Errorf("Failed to ping %v", err.Error())
	}

	err = c.Chat("5002")
	if err != nil {
		t.Errorf("Failed to ping %v", err.Error())
	}
}
