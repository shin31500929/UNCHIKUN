## Go Discord Bot（UNCHIKUN）

このリポジトリは Go で動作するシンプルな Discord Bot のスターターです。`!ping` に `Pong!` で応答します。

### 必要要件
- Go 1.21+
- Discord Bot トークン（Botを作成しTOKENを取得）

### 環境変数
- `DISCORD_TOKEN`: DiscordのBotトークン（`.env` に設定）
- `OPENAI_API_KEY`: OpenAIのAPIキー（AIうんちく生成に使用）

リポジトリ直下に `.env` を作成し、以下のように設定してください。

```
DISCORD_TOKEN=YOUR_BOT_TOKEN_HERE
OPENAI_API_KEY=YOUR_OPENAI_API_KEY
```

### セットアップ
```bash
go mod tidy
go run ./...
```

またはビルドして実行:
```bash
go build -o bot
./bot
```

### 使い方
Botをサーバーに追加した上で、任意のテキストチャンネルにて:
```
!ping
```
と送ると、`Pong!` と返答します。

### 開発メモ
- メッセージ閲覧には `Message Content Intent` が必要です。DiscordのDeveloper Portalで有効化してください。
- ハンドラーは `onMessageCreate` に実装しています。コマンドを増やしたい場合は `switch` に追加してください。


