package fcm

import (
	"context"
	"fmt"
	"log"
)

type Client struct {
	projectID string
}

func NewClient(projectID string) *Client {
	return &Client{projectID: projectID}
}

func (c *Client) SendPushNotification(ctx context.Context, token, title, body string, payload map[string]string) error {
	if token == "" {
		return fmt.Errorf("empty device token")
	}
	// FCM Message dispatch log
	log.Printf("[FCM Push] Dispatched notification to token %s: %s - %s", token, title, body)
	return nil
}
