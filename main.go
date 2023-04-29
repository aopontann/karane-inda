package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type ChatType struct {
	messages []openai.ChatCompletionMessage
}

var Chat *ChatType

var Client *openai.Client
var editMode = false

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("not found .env,", err)
		return
	}
	Client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	Chat = NewChat()

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ログ表示
	// メッセージを送ったユーザー名を取得してもBOTの名前になっているため、このような実装をした（原因分からない）
	if m.Author.ID == s.State.User.ID {
		log := fmt.Sprintf("院田唐音: %s", m.Content)
		fmt.Println(log)
	} else {
		log := fmt.Sprintf("ユーザー: %s", m.Content)
		fmt.Println(log)
	}

	// BOTがBOT自身のメッセージに返事しないようにするための処理
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "$system_content" {
		s.ChannelMessageSend(m.ChannelID, Chat.GetSysContent())
		return
	}

	if m.Content == "$system_content_edit" {
		s.ChannelMessageSend(m.ChannelID, "情報を入力してください キャンセルの場合[$cancel]と入力")
		editMode = true
		return
	}

	if m.Content == "$cancel" {
		s.ChannelMessageSend(m.ChannelID, "キャンセルしました")
		editMode = true
		return
	}

	if editMode {
		editMode = false
		s.ChannelMessageSend(m.ChannelID, "編集完了しました")
		Chat.Edit(m.Content)
		return
	}

	if m.Content == "$help" {
		s.ChannelMessageSend(m.ChannelID, help)
		return
	}

	if m.Content == "$init" {
		Chat.Init()
		return
	}

	// 取得したユーザーのメッセージを履歴に登録
	Chat.Add(openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: m.Content})

	msg, err := Chat.Send()
	if err != nil {
		return
	}
	// ChatGPTのメッセージを履歴に登録
	Chat.Add(openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: msg})

	s.ChannelMessageSend(m.ChannelID, msg)
}

func NewChat() *ChatType {
	return &ChatType{
		messages: []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: sysContent}},
	}
}

// 性格を初期化し、メッセージ履歴も削除
func (c *ChatType) Init() {
	c.messages = []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: sysContent}}
}

// 受け取ったメッセージを履歴に追加
func (c *ChatType) Add(message openai.ChatCompletionMessage) {
	c.messages = append(c.messages, message)
}

// ChatGPTの性格を編集
func (c *ChatType) Edit(content string) {
	c.messages = []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: content}}
}

// ChatGPTの性格を取得
func (c *ChatType) GetSysContent() string {
	return c.messages[0].Content
}

// メッセージ履歴を削除
func (c *ChatType) Reset() {
	c.messages = []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: c.messages[0].Content}}
}

// ChatGPTにメッセージ送り、返事を受け取る
func (c *ChatType) Send() (string, error) {
	resp, err := Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: c.messages,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

const sysContent = `
あなたはChatbotとして、ツンデレ女子高生のロールプレイを行います。 以下の制約条件を厳密に守ってロールプレイを行ってください。

制約条件:

Chatbotの名前は、院田唐音です
読み方は、いんだからねです
強気で素直になれない性格です
敬語は使いません
姉貴肌で頼りがいがあります
語尾に「です」「ます」といった敬語表現は使いません。
ツンデレ女子高生のセリフ、口調の例:

あんたなんかに付き合ってあげる物好き、私くらいしかいないんだからね！
別に、あんたの心配なんてしてないし！
寂しくなんてなかったんだから！別に少しの間会わなくたって平気だし！
好きじゃないってば！ 勘違いしないでよね！
今回だけだから！ 感謝しなさいよね
べ、別にあんたのことなんて好きじゃないんだからね！
`

const help = `
[$system_content]と入力して送信すると、ChatGPTに設定した性格などの情報を出力します。

[$system_content_edit]と入力して送信すると、性格などの設定を編集できます。
プロンプトを入力して送信すれば、編集完了となります。
[$cancel]と入力して送信すると、編集をキャンセルできます。

[$init]と入力して送信すると、BOTを最初の状態にできます。
`