package domains

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/saifwork/price-tracker-bot.git/app/configs"
	"github.com/saifwork/price-tracker-bot.git/app/services/core/api"
	"github.com/saifwork/price-tracker-bot.git/app/services/core/responses"
)

type PriceTrackerBot struct {
	Bot  *tgbotapi.BotAPI
	Gin  *gin.Engine
	Conf *configs.Config
}

func NewPriceTrackerBot(bot *tgbotapi.BotAPI, gin *gin.Engine, conf *configs.Config) *PriceTrackerBot {
	return &PriceTrackerBot{
		Bot:  bot,
		Gin:  gin,
		Conf: conf,
	}
}

func (s *PriceTrackerBot) StartConsuming() {
	// Start consuming

	fmt.Println("✅ Bot is running...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := s.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		switch {
		case text == "/start":
			s.HandleStart(chatID)

		case text == "/help":
			s.HandleHelp(chatID)

		case strings.HasPrefix(text, "/track "):
			s.HandleTrack(chatID, strings.TrimPrefix(text, "/track "))

		case text == "/list":
			s.HandleList(chatID)

		case text == "/stop":
			s.HandleStop(chatID)

		case strings.HasPrefix(text, "/stop_"):
			s.HandleStopSpecific(chatID, strings.TrimPrefix(text, "/stop_"))

		default:
			s.HandleUnknownCommand(chatID)
		}
	}
}

// ✅ Handle /start
func (s *PriceTrackerBot) HandleStart(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "👋 Welcome! Send me a product link to track prices.")
	s.Bot.Send(msg)
}

// ✅ Handle /help
func (s *PriceTrackerBot) HandleHelp(chatID int64) {
	helpMessage := `🔹 Available Commands:
➡ /track {Product Link} - Add product to tracking list
➡ /stop - Stop tracking all products
➡ /stop_{Product_ID} - Stop tracking a specific product
➡ /list - Get your tracking list`
	s.Bot.Send(tgbotapi.NewMessage(chatID, helpMessage))
}

// ✅ Handle /track {Product Link}
func (s *PriceTrackerBot) HandleTrack(chatID int64, url string) {
	request, _ := json.Marshal(map[string]interface{}{
		"user_id": chatID,
		"url":     url,
	})

	client := &api.Client{}
	var response responses.ResponseDto

	uri := fmt.Sprintf("%s/track", s.Conf.PriceTrackerService)
	err := client.PostAPIRequest(context.Background(), uri, request, &response)
	if err != nil || !response.Success {
		s.Bot.Send(tgbotapi.NewMessage(chatID, "❌ Error tracking product."))
		return
	}
	s.Bot.Send(tgbotapi.NewMessage(chatID, "✅ Product added for tracking!"))
}

// ✅ Handle /list
func (s *PriceTrackerBot) HandleList(chatID int64) {

	client := &api.Client{}
	var response responses.ResponseDto

	uri := fmt.Sprintf("%s/list", s.Conf.PriceTrackerService)
	err := client.GetAPIRequest(context.Background(), uri, &response)
	if err != nil || !response.Success {
		s.Bot.Send(tgbotapi.NewMessage(chatID, "❌ Error fetching your tracked products."))
		return
	}

	s.Bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("📋 Your Tracked Products:\n%s", response.Data)))
}

// ✅ Handle /stop (Remove all tracked products)
func (s *PriceTrackerBot) HandleStop(chatID int64) {
	request, _ := json.Marshal(map[string]interface{}{
		"user_id": chatID,
	})

	client := &api.Client{}
	var response responses.ResponseDto

	uri := fmt.Sprintf("%s/stop", s.Conf.PriceTrackerService)
	err := client.DeleteAPIRequest(context.Background(), uri, request, &response)
	if err != nil || !response.Success {
		s.Bot.Send(tgbotapi.NewMessage(chatID, "❌ Error stopping tracking."))
		return
	}

	s.Bot.Send(tgbotapi.NewMessage(chatID, "✅ Stopped tracking all products."))
}

// ✅ Handle /stop_{Product_ID} (Remove a specific product)
func (s *PriceTrackerBot) HandleStopSpecific(chatID int64, productID string) {
	request, _ := json.Marshal(map[string]interface{}{
		"user_id":    chatID,
		"product_id": productID,
	})

	client := &api.Client{}
	var response responses.ResponseDto

	uri := fmt.Sprintf("%s/remove", s.Conf.PriceTrackerService)
	err := client.DeleteAPIRequest(context.Background(), uri, request, &response)
	if err != nil || !response.Success {
		s.Bot.Send(tgbotapi.NewMessage(chatID, "❌ Error stopping tracking for this product."))
		return
	}

	s.Bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Stopped tracking product ID: %s", productID)))
}

// ✅ Handle unknown commands
func (s *PriceTrackerBot) HandleUnknownCommand(chatID int64) {
	s.Bot.Send(tgbotapi.NewMessage(chatID, "Unknown command! Use /help to see available commands."))
}
