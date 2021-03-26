package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/TianqiuHuang/grpc-client-app/pd/auth"
)

type Client struct {
	authClient auth.AuthServiceClient
}

func New() {

}

func (c *Client) Validate(token string) (string, error) {
	resp, err := c.authClient.Validate(context.Background(), &auth.ValidateRequest{
		RawIdToken: token,
		ClaimNames: []string{"email"},
	})
	if err != nil {
		return "", err
	}
	if resp.Email == "" {
		return "", fmt.Errorf("email must not be empty")
	}
	return resp.Email, nil
}

func (c *Client) ValidateAdmin(token string) (string, error) {
	resp, err := c.authClient.Validate(context.Background(), &auth.ValidateRequest{
		RawIdToken: token,
		ClaimNames: []string{"email"},
	})
	if err != nil {
		return "", err
	}
	if resp.Email == "" {
		return "", fmt.Errorf("email must not be empty")
	}
	var admin = false
	for i := range resp.Groups {
		if strings.Contains(resp.Groups[i], "admin") {
			admin = true
		}
	}
	if !admin {
		return "", fmt.Errorf("email is not an admin")
	}
	return resp.Email, nil
}
