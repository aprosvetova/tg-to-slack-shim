package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"net/http"
	"strings"
)

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}

		code := c.Writer.Status()

		c.JSON(code, responseWrapper{
			Ok:          false,
			ErrorCode:   &code,
			Description: err.Error(),
		})
	}
}

type responseWrapper struct {
	Ok          bool        `json:"ok"`
	ErrorCode   *int        `json:"error_code,omitempty"`
	Description string      `json:"description,omitempty"`
	Result      interface{} `json:"result,omitempty"`
}

func apiHandler(c *gin.Context) {
	token, err := extractSlackToken(c.Param("token"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	method, ok := methods[c.Param("method")]
	if !ok {
		c.AbortWithError(http.StatusNotImplemented, errors.New("method not implemented"))
		return
	}

	result := method(slack.New(token), c)
	if c.IsAborted() {
		return
	}

	c.JSON(200, responseWrapper{
		Ok:     true,
		Result: result,
	})
}

func extractSlackToken(token string) (string, error) {
	p := strings.Index(token, "xoxb-")
	if p < 0 {
		return "", errors.New("wrong token format")
	}
	return token[p:], nil
}
