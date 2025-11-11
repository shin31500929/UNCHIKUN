package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// .env を読み込む（存在しない場合はそのまま進む）
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("環境変数 DISCORD_TOKEN が設定されていません")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Discordセッションの作成に失敗しました: %v", err)
	}
	// 最小のIntent（スラッシュコマンドはMessage Content不要）
	session.Identify.Intents = discordgo.IntentsGuilds

	// スラッシュコマンドのハンドラー登録
	session.AddHandler(onInteractionCreate)
	log.Println("ハンドラー登録完了")

	// 接続開始
	if err := session.Open(); err != nil {
		log.Fatalf("Discordへの接続に失敗しました: %v", err)
	} else {
		log.Println("Discordへの接続成功")
	}
	log.Println("Bot が起動しました。Ctrl+C で終了します。")

	// コマンド登録
	cmd := &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "うんちくを返します（AI生成）",
	}
	createdCmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", cmd)
	if err != nil {
		log.Fatalf("アプリケーションコマンド登録に失敗しました: %v", err)
	}
	log.Printf("コマンド登録: /%s", createdCmd.Name)

	// グレースフルシャットダウン
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("シャットダウン中...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := shutdown(session, shutdownCtx, createdCmd); err != nil {
		log.Printf("シャットダウンでエラー: %v", err)
	}
	log.Println("終了しました。")
}

func shutdown(s *discordgo.Session, _ context.Context, createdCmd *discordgo.ApplicationCommand) error {
	// 作成したアプリケーションコマンドを削除
	if createdCmd != nil {
		if err := s.ApplicationCommandDelete(s.State.User.ID, "", createdCmd.ID); err != nil {
			log.Printf("コマンド削除に失敗しました: %v", err)
		}
	}
	return s.Close()
}

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if i.ApplicationCommandData().Name != "ping" {
		return
	}
	// まずはDeferredで応答（3秒制限回避）
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	// 別ゴルーチンでAI生成してフォローアップ
	go func(interaction *discordgo.Interaction) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		text := generateUnchiku(ctx)
		_, _ = s.InteractionResponseEdit(interaction, &discordgo.WebhookEdit{
			Content: &text,
		})
	}(i.Interaction)
}

func generateUnchiku(ctx context.Context) string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "（設定ヒント）環境変数 OPENAI_API_KEY を設定するとAIがうんちくを返します。"
	}
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "あなたは博識で、短く面白い“うんちく”を日本語で1〜2文で返します。事実ベースで、出典表記や箇条書きは不要。絵文字は使わない。"},
				{Role: openai.ChatMessageRoleUser, Content: "ランダムな面白いうんちくを1つ教えて。"},
			},
			Temperature: 0.8,
			MaxTokens:   120,
		},
	)
	if err != nil || len(resp.Choices) == 0 {
		return "うんちくの生成に失敗しました。時間をおいて再度お試しください。"
	}
	return resp.Choices[0].Message.Content
}
