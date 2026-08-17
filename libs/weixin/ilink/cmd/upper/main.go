package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fasionchan/goutils/libs/logging"
	"github.com/fasionchan/goutils/libs/qrcode"
	"github.com/fasionchan/goutils/libs/weixin/ilink"
)

const botType = 3

func main() {
	log.SetFlags(log.LstdFlags)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	client, err := ilink.NewClient()
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}
	client.WithLogger(logging.GetNopLogger())

	bot, err := login(ctx, client)
	if err != nil {
		return err
	}
	bot.WithLogger(logging.GetNopLogger())

	if err := bot.NotifyStart(ctx); err != nil {
		log.Printf("notify start: %v", err)
	}
	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := bot.NotifyStop(stopCtx); err != nil {
			log.Printf("notify stop: %v", err)
		}
	}()

	log.Printf("logged in, waiting for messages")
	return serve(ctx, bot)
}

func login(ctx context.Context, client *ilink.Client) (*ilink.BotClient, error) {
	for {
		qr, err := client.GetBotQrcode(ctx, botType)
		if err != nil {
			return nil, fmt.Errorf("get bot qrcode: %w", err)
		}

		fmt.Println(qr.QrcodeImgContent)
		if err := qrcode.EncodeToTerminal(qr.QrcodeImgContent, os.Stdout); err != nil {
			return nil, fmt.Errorf("print qrcode: %w", err)
		}
		log.Printf("scan the qrcode to login")

		bot, err := waitBotClient(ctx, client, qr.Qrcode)
		if err == nil {
			return bot, nil
		}
		if ctx.Err() != nil {
			return nil, err
		}
		if strings.Contains(err.Error(), "expired") {
			log.Printf("qrcode expired, fetching a new one")
			continue
		}
		return nil, err
	}
}

func waitBotClient(ctx context.Context, client *ilink.Client, qrcode string) (*ilink.BotClient, error) {
	for {
		result, err := client.GetQrcodeStatus(ctx, qrcode)
		if err != nil {
			return nil, fmt.Errorf("get qrcode status: %w", err)
		}

		switch result.Status {
		case ilink.GetQrCodeStatusConfirmed:
			return ilink.NewBotClient(result.BaseUrl, result.BotToken)
		case ilink.GetQrCodeStatusExpired:
			return nil, fmt.Errorf("qrcode expired")
		case ilink.GetQrCodeStatusWait:
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Second):
			}
		default:
			return nil, fmt.Errorf("unknown qrcode status: %s", result.Status)
		}
	}
}

func serve(ctx context.Context, bot *ilink.BotClient) error {
	getUpdatesBuf := ""
	for {
		result, err := bot.GetUpdates(ctx, getUpdatesBuf)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("get updates: %v", err)
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(2 * time.Second):
			}
			continue
		}
		if result == nil {
			continue
		}
		if result.GetUpdatesBuf != "" {
			getUpdatesBuf = result.GetUpdatesBuf
		}

		for _, msg := range result.Msgs {
			if err := replyUpper(ctx, bot, msg); err != nil {
				log.Printf("reply: %v", err)
			}
		}
	}
}

func replyUpper(ctx context.Context, bot *ilink.BotClient, msg *ilink.WeixinMessage) error {
	if msg == nil || msg.MessageType != ilink.MessageTypeUser {
		return nil
	}

	text := messageText(msg)
	if text == "" {
		return nil
	}

	upper := strings.ToUpper(text)
	log.Printf("from=%s text=%q reply=%q", msg.FromUserId, text, upper)

	return bot.SendMessage(ctx, &ilink.WeixinMessage{
		ToUserId:     msg.FromUserId,
		ClientId:     fmt.Sprintf("%d", time.Now().UnixNano()),
		MessageType:  ilink.MessageTypeBot,
		MessageState: ilink.MessageStateFinish,
		ContextToken: msg.ContextToken,
		ItemList: []*ilink.MessageItem{
			{
				Type:     ilink.MessageItemTypeText,
				TextItem: &ilink.TextItem{Text: upper},
			},
		},
	})
}

func messageText(msg *ilink.WeixinMessage) string {
	var parts []string
	for _, item := range msg.ItemList {
		if item == nil {
			continue
		}
		if item.TextItem != nil && item.TextItem.Text != "" {
			parts = append(parts, item.TextItem.Text)
			continue
		}
		if item.VoiceItem != nil && item.VoiceItem.Text != "" {
			parts = append(parts, item.VoiceItem.Text)
		}
	}
	return strings.Join(parts, "\n")
}
