package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"net/http"
	"strings"
)

var methods = map[string]func(s *slack.Client, c *gin.Context) interface{}{
	"getMe":       getMe,
	"sendMessage": sendMessage,
}

func getMe(s *slack.Client, c *gin.Context) interface{} {
	bot, err := s.AuthTest()
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return nil
	}

	return getMeResponse{
		ID:                      1,
		IsBot:                   true,
		FirstName:               bot.User,
		Username:                bot.BotID,
		CanJoinGroups:           false,
		CanReadAllGroupMessages: false,
		SupportsInlineQueries:   false,
	}
}

type getMeResponse struct {
	ID                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	Username                string `json:"username"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
}

func sendMessage(s *slack.Client, c *gin.Context) interface{} {
	var req sendMessageRequest
	c.ShouldBindQuery(&req)
	c.ShouldBind(&req)

	if req.Text == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("text must be non-empty"))
		return nil
	}

	if !strings.HasPrefix(req.ChatID, "@") {
		c.AbortWithError(http.StatusBadRequest, errors.New("only string starting with @ allowed in chat_id"))
		return nil
	}

	if req.ParseMode != "" && req.ParseMode != "Markdown" {
		c.AbortWithError(http.StatusNotImplemented, errors.New("only Markdown parse mode is supported"))
		return nil
	}

	_, _, err := s.PostMessage(req.ChatID[1:], slack.MsgOptionBlocks(
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, convertMarkdown(req.Text), false, false),
			nil, nil,
		),
	))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	return nil
}

type sendMessageRequest struct {
	ChatID    string `json:"chat_id" form:"chat_id"`
	Text      string `json:"text" form:"text"`
	ParseMode string `json:"parse_mode" form:"parse_mode"`
}
