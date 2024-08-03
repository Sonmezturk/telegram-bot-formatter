package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/Sonmezturk/telegram-bot-formatter/db"
	"github.com/Sonmezturk/telegram-bot-formatter/order"
)

func main() {
	err := godotenv.Load()
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	db.MongoInit()
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)
	var awaitingDateRange = make(map[int64]bool)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		chatID := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		fmt.Println(chatID)
		if update.Message.IsCommand() {
			fmt.Println("here")
			switch update.Message.Command() {
			case "order":
				msg.Text = "Please send me the PDF file."
			case "filter":
				msg.Text = "Please send me the 'from' and 'to' dates in format YYYY-MM-DD, followed by the caption. Example: 2006-01-02 2006-01-03 mert"
				msg.ReplyMarkup = tgbotapi.ForceReply{
					ForceReply: true,
					Selective:  true,
				}
				awaitingDateRange[chatID] = true
			case "trend":
				trend(bot)
				msg.Text = "Please send me the PDF file."
			default:
				msg.Text = "I don't know that command"
			}

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} else if update.Message.Document != nil || update.Message.Photo != nil {
			chatID := update.Message.Chat.ID
			userName := update.Message.From.UserName
			fmt.Println(chatID)
			fmt.Println(userName)

			fileID := update.Message.Document.FileID
			fileURL, err := bot.GetFileDirectURL(fileID)
			if err != nil {
				log.Panic(err)
			}

			originalFilename := strings.ReplaceAll(update.Message.Document.FileName, ".pdf", "")
			result := parser(fileURL)
			caption := update.Message.Caption
			if caption != "" {
				originalFilename = caption
			}
			fmt.Println(originalFilename)

			order.SaveOrder(result, time.Now(), userName, originalFilename)
			formattedDate := time.Now().Format("2006-01-02_15-04-05")
			fileName := fmt.Sprintf("%s_%s_orders.csv", formattedDate, originalFilename)
			prepareCsv(result, fileName)

			//-1001810967398
			msg := tgbotapi.NewDocumentUpload(-1002099498788, fileName)
			_, err = bot.Send(msg)
			if err != nil {
				log.Panic(err)
			}
		} else if awaitingDateRange[chatID] {
			// If we're awaiting a date range from this user, process their message accordingly.
			messages := strings.Split(update.Message.Text, " ")
			if len(messages) == 3 {
				fromDate, err := time.Parse("2006-01-02", messages[0])
				if err != nil {
					log.Panic(err) // Handle the error appropriately
				}
				toDate, err := time.Parse("2006-01-02", messages[1])
				if err != nil {
					log.Panic(err) // Handle the error appropriately
				}
				msg.Text = "Thank you. I will now filter the data from: " + fromDate.Format("2006-01-02") + " to: " + toDate.Format("2006-01-02") + " for " + messages[2] + "."
				aggregatedItems, _ := order.GetFilteredAggregatedData(fromDate, toDate, messages[2])
		
				fileName := fmt.Sprintf("%s_aggreagatedOrders.csv", messages[2])
				if err := prepareCsvForAggregatedItems(aggregatedItems, fileName); err != nil {
					log.Fatalf("Failed to prepare CSV: %v", err)
				}
				msg := tgbotapi.NewDocumentUpload(-1002099498788, fileName)
				_, err = bot.Send(msg)
				if err != nil {
					log.Panic(err)
				}

				err = os.Remove(fileName)
				if err != nil {
					log.Panic(err)
				}
				log.Printf("%s file removed", fileName)
				awaitingDateRange[chatID] = false
			} else {
				// If the user's message doesn't match the expected format, prompt them again.
				msg.Text = "Please send the dates in the correct format: YYYY-MM-DD YYYY-MM-DD."
				msg.ReplyMarkup = tgbotapi.ForceReply{
					ForceReply: true,
					Selective:  true,
				}
			}
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} 

	}
}

func sendMessageToChannel(bot *tgbotapi.BotAPI, channel string, message string) {
	msg := tgbotapi.NewMessageToChannel(channel, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

func sendImageToChannel(bot *tgbotapi.BotAPI, channel string, imageURL string) {
	chatID, err := strconv.ParseInt(channel, 10, 64)
if err != nil {
	log.Panic(err)
}
	msg := tgbotapi.NewPhotoUpload(chatID, imageURL)
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

func sendImageToChannel(bot *tgbotapi.BotAPI, channel int64, imageURL string) {
	tmpFile, err := os.CreateTemp("", "image-*.png")
	if err != nil {
		log.Panic(err)
	}
	defer os.Remove(tmpFile.Name())

	resp, err := http.Get(imageURL)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		log.Panic(err)
	}

	// Send the image
	// msg := tgbotapi.NewPhotoUpload(channel, tmpFile.Name())
	// _, err = bot.Send(msg)
	// if err != nil {
	// 	log.Panic(err)
	// }
}
