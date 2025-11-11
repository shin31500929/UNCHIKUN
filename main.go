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
	// 必要なIntentを設定（メッセージ受信）
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	// ハンドラー登録
	session.AddHandler(onMessageCreate)

	// 接続開始
	if err := session.Open(); err != nil {
		log.Fatalf("Discordへの接続に失敗しました: %v", err)
	}
	log.Println("Bot が起動しました。Ctrl+C で終了します。")

	// グレースフルシャットダウン
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("シャットダウン中...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := shutdown(session, shutdownCtx); err != nil {
		log.Printf("シャットダウンでエラー: %v", err)
	}
	log.Println("終了しました。")
}

func shutdown(s *discordgo.Session, _ context.Context) error {
	return s.Close()
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 自分のメッセージは無視
	if m.Author == nil || m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "!ping":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "!help":
		_, _ = s.ChannelMessageSend(m.ChannelID, "使い方: `!ping` で応答します。")
	}
}


