package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"os"

	"github.com/TianqiuHuang/grpc-client-app/pd/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare(client *Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		bearerToken := c.GetHeader("Authorization")
		IDToken := strings.Split(bearerToken, " ")
		if len(IDToken) != 2 {
			if os.Getenv("DEBUG") != "" {
				fmt.Println(bearerToken)
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
				"error": "failed to auth, expected header: 'Authorization: bearer <token>'",
			})
			return
		}
		//ctx, _ := c.Get("SpanContext")
		//email, err := client.Validate(IDToken[1], ctx.(context.Context))
		email, err := client.Validate(IDToken[1], c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
			return
		}
		c.Set("id", email)
		c.Next()
	}
}

func AdminAuthMiddleWare(client *Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		bearerToken := c.GetHeader("Authorization")
		IDToken := strings.Split(bearerToken, " ")
		if len(IDToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
				"error": "failed to auth, expected header: 'Authorization: bearer <token>'",
			})
			return
		}
		ctx, _ := c.Get("SpanContext")
		email, err := client.ValidateAdmin(IDToken[1], ctx.(context.Context))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
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

func (c *Client) Validate(token string, ctx context.Context) (string, error) {
	resp, err := c.authClient.Validate(ctx, &auth.ValidateRequest{
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

func (c *Client) ValidateAdmin(token string, ctx context.Context) (string, error) {
	resp, err := c.authClient.Validate(ctx, &auth.ValidateRequest{
		RawIdToken: token,
		ClaimNames: []string{"email", "groups"},
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
