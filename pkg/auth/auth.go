package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/new-adventure-areolite/grpc-app-server/pd/auth"
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

		email, isAdmin, err := client.ValidateAdmin(c.Request.Context(), IDToken[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
			return
		}
		c.Set("id", email)
		if isAdmin {
			c.Set("user-type", "admin")
		} else {
			c.Set("user-type", "normal")
		}
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
		ctxInterface, _ := c.Get("SpanContext")
		ctx := ctxInterface.(context.Context)
		email, isAdmin, err := client.ValidateAdmin(ctx, IDToken[1])
		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "this request needs admin access",
			})
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
			return
		}
		c.Set("id", email)
		if isAdmin {
			c.Set("user-type", "admin")
		} else {
			c.Set("user-type", "normal")
		}
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

func (c *Client) ValidateAdmin(ctx context.Context, token string) (string, bool, error) {
	resp, err := c.authClient.Validate(ctx, &auth.ValidateRequest{
		RawIdToken: token,
		ClaimNames: []string{"email", "groups"},
	})
	if err != nil {
		return "", false, err
	}
	if resp.Email == "" {
		return "", false, fmt.Errorf("email must not be empty")
	}
	var isAdmin = false
	for i := range resp.Groups {
		if strings.Contains(resp.Groups[i], "admin") {
			isAdmin = true
		}
	}
	return resp.Email, isAdmin, nil
}
