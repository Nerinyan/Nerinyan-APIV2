package webhook

import (
	"bytes"
	"github.com/Nerinyan/Nerinyan-APIV2/bodyStruct"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/pterm/pterm"
	"net/http"
	"time"
)

type discordWebhookBody struct {
	Content  string   `json:"content"`
	Embeds   []embeds `json:"embeds"`
	Username string   `json:"username"`
}
type embeds struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
	Footer      footer `json:"footer"`
	Timestamp   string `json:"timestamp"`
}
type footer struct {
	Text string `json:"text"`
}

//{
//  "content": null,
//  "embeds": [
//    {
//      "title": "What's this?",
//      "description": "da",
//      "color": 65535
//    },
//    {
//      "title": "What's this?",
//      "description": "da",
//      "color": 65280
//    },
//    {
//      "color": 16711680
//    }
//  ]
//}

func DiscordError(v *bodyStruct.ErrorStruct) {
	defer func() {
		if r := recover(); r != nil {
			pterm.Error.WithShowLineNumber().Printfln("%+v", r)
		}
	}()
	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	body := discordWebhookBody{
		Embeds: []embeds{
			{
				Title:       "Error Code: *" + v.Code + "*",
				Description: "```json\n" + string(*utils.ToJsonIndentString(v)) + "```",
				Color:       16711680,
				Footer: footer{
					Text: "RequestID: " + v.RequestId,
				},
				Timestamp: now,
			},
		},
		Username: "Nerinyan-APIv2",
	}

	resp, err := http.Post(config.Config.Discord.Webhook.Error, "application/json", bytes.NewReader(*utils.ToJsonString(body)))
	if err != nil {
		pterm.Error.WithShowLineNumber().Printfln("%+v", err)
	}
	if !resp.Close {
		_ = resp.Body.Close()
	}
}
func DiscordInfo(v *bodyStruct.ErrorStruct) {
	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	body := discordWebhookBody{
		Embeds: []embeds{
			{
				Title:       "Info",
				Description: "```json\n" + string(*utils.ToJsonIndentString(v)) + "```",
				Color:       65535,
				Footer: footer{
					Text: "RequestID: " + v.RequestId,
				},
				Timestamp: now,
			},
		},
		Username: "Nerinyan-APIv2",
	}

	resp, err := http.Post(config.Config.Discord.Webhook.Info, "application/json", bytes.NewReader(*utils.ToJsonString(body)))
	if err != nil {
		pterm.Error.Println(err)
	}
	if !resp.Close {
		_ = resp.Body.Close()
	}
}
func DiscordInfoStartUP() {
	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	body := discordWebhookBody{
		Embeds: []embeds{
			{
				Title:       "Nerinyan API Server has been started",
				Description: "Startup Time: " + string(now),
				Color:       65535,
				Footer: footer{
					Text: "Start UP",
				},
				Timestamp: now,
			},
		},
		Username: "Nerinyan-APIv2",
	}

	resp, err := http.Post(config.Config.Discord.Webhook.Startup, "application/json", bytes.NewReader(*utils.ToJsonString(body)))
	if err != nil {
		pterm.Error.Println(err)
	}
	if !resp.Close {
		_ = resp.Body.Close()
	}
}
