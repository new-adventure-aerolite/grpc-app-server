package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/TianqiuHuang/grpc-client-app/pd/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare(client *Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		IDToken := strings.Split(bearerToken, " ")
		if len(IDToken) != 2 {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf(
				"failed to auth, expected header: 'Authorization: bearer <token>'",
			))
			return
		}
		email, err := client.Validate(IDToken[1])
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		c.Set("id", email)
		c.Next()
	}
}

func AdminAuthMiddleWare(client *Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		IDToken := strings.Split(bearerToken, " ")
		if len(IDToken) != 2 {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf(
				"failed to auth, expected header: 'Authorization: bearer <token>'",
			))
			return
		}
		email, err := client.ValidateAdmin(IDToken[1])
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		c.Set("id", email)
		c.Next()
	}
}

type Client struct {
	authClient auth.AuthServiceClient
}

func New(client auth.AuthServiceClient) *Client {
	return &Client{
		authClient: client,
	}
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
