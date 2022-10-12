package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/livekit/protocol/auth"
)

func createToken(apiKey string, apiSecret string, p string, name string, room string, metadata string, validFor string) string {
	grant := &auth.VideoGrant{}
	grant.RoomJoin = true
	grant.Room = room

	at := accessToken(apiKey, apiSecret, grant, p)

	if metadata != "" {
		at.SetMetadata(metadata)
	}
	if name == "" {
		name = p
	}
	at.SetName(name)
	if validFor != "" {
		if dur, err := time.ParseDuration(validFor); err == nil {
			fmt.Println("valid for (mins): ", int(dur/time.Minute))
			at.SetValidFor(dur)
		}
	}

	token, _ := at.ToJWT()

	return token
}

func accessToken(apiKey string, apiSecret string, grant *auth.VideoGrant, identity string) *auth.AccessToken {
	if apiKey == "" && apiSecret == "" {
		// not provided, don't sign request
		return nil
	}
	at := auth.NewAccessToken(apiKey, apiSecret).
		AddGrant(grant).
		SetIdentity(identity)
	return at
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message == nil { // ignore any non-Message updates
				continue
			}

			if !update.Message.IsCommand() { // ignore any non-command Messages
				continue
			}

			// Create a new MessageConfig. We don't have text yet,
			// so we leave it empty.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			// msg.ReplyToMessageID = update.Message.MessageID

			// Extract the command from the Message.
			switch update.Message.Command() {
			case "help":
				msg.Text = "I understand /soy nombre and /status."
			case "soy (nombre)":
				name := update.Message.CommandArguments()
				msg.Text = "https://www.accreativos.com/livekit/index.html#/room?url=https%3A%2F%2Flivekit.accreativos.com%2F&token=" +
					createToken(
						os.Getenv("LIVEKIT_API"), os.Getenv("LIVEKIT_SECRET"),
						String(10), name, "room", "metadata", "validFor") +
					"&videoEnabled=1&audioEnabled=1&simulcast=1&dynacast=1&adaptiveStream=1"
				// "&videoEnabled=1&audioEnabled=1&simulcast=1&dynacast=1&adaptiveStream=1&videoDeviceId=9aa58fb1841f23afb837ba11afd74a7d213bc4343e46dda734452ce558681fe9"
			case "status":
				msg.Text = "I'm ok."
			default:
				msg.Text = "I don't know that command"
			}

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}
