package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

// HandleTimenCommand は /timen コマンドを処理します
func HandleTimenCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// コマンドのオプションを取得
	options := i.ApplicationCommandData().Options
	var unchiku, location, time string

	for _, option := range options {
		switch option.Name {
		case "unchiku":
			unchiku = option.StringValue()
		case "location":
			location = option.StringValue()
		case "time":
			time = option.StringValue()
		}
	}

	// メッセージを生成
	message := fmt.Sprintf("@学生\nお疲れさまです！\n今日も対面あります！\n今日のうんちく：%s\n\n\n場所：%s\n\n時間：%s ~", unchiku, location, time)

	// 環境変数からチャンネルIDを取得
	channelID := os.Getenv("TAIMEN_CHANNEL")
	if channelID == "" {
		// チャンネルが設定されていない場合はエラーメッセージを返す
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "エラー: 環境変数 TAIMEN_CHANNEL が設定されていません。",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// 指定されたチャンネルにメッセージを送信
	_, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("エラー: メッセージの送信に失敗しました: %v", err),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// 成功レスポンスを返す
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "対面の通知を送信しました！",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

