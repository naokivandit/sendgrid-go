package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("コンフィグの読み込みに失敗しました: ", err)
	}

	if err := sendEmail(config); err != nil {
		log.Println("メールの送信に失敗しました: ", err)
	} else {
		fmt.Println("メールが送信されました。")
	}
}

type Config struct {
	SendgridAPIKey  string
	UnsubscribeLink string
	FromEmail       string
	ToEmail         string
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		SendgridAPIKey:  os.Getenv("SENDGRID_API_KEY"),
		UnsubscribeLink: os.Getenv("UNSUBSCRIBE_LINK"),
		FromEmail:       os.Getenv("FROM"),
		ToEmail:         os.Getenv("TO"),
	}, nil
}

func sendEmail(config *Config) error {
	client := sendgrid.NewSendClient(config.SendgridAPIKey)

	from := mail.NewEmail("Example From User", config.FromEmail)
	to := mail.NewEmail("Example To User", config.ToEmail)
	subject := "Sending with SendGrid is Fun"
	content := mail.NewContent("text/plain", "メールの本文（テキスト形式）")

	sg := mail.NewV3MailInit(from, subject, to, content)
	sg.SetHeader("List-Unsubscribe", fmt.Sprintf("<mailto:%s>, <%s>", config.FromEmail, config.UnsubscribeLink))

	if _, err := client.Send(sg); err != nil {
		return fmt.Errorf("メール送信エラー: %w", err)
	}
	return nil
}
