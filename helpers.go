package main

import (
	"github.com/aprosvetova/tg-to-slack-shim/slackdown"
	"github.com/russross/blackfriday/v2"
)

func convertMarkdown(tg string) string {
	return string(blackfriday.Run([]byte(tg),
		blackfriday.WithRenderer(&slackdown.Renderer{})))
}
