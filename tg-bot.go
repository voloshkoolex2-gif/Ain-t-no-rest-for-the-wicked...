package main

import (
 "fmt"
 "log"
 "sync"

 tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
 // Map to store user scores.
 //in know i can use heidiSQL data base but i'm to lazy
 scores = make(map[int64]int)
 
 // Mutex for safe concurrent updates of scores to prevent data races.
 mu sync.Mutex
)

func main() {
 // Connect to the bot using the token provided by BotFather.
 bot, err := tgbotapi.NewBotAPI("YOUR_BOTFATHER_TOKEN_HERE")
 if err != nil {
  log.Panic("Connection error: ", err)
 }

 bot.Debug = true // Set to false in production to reduce log spam
 log.Printf("Authorized on account %s", bot.Self.UserName)

 u := tgbotapi.NewUpdate(0)
 u.Timeout = 60

 // Channel to receive updates from Telegram (messages, button clicks, etc.)
 updates := bot.GetUpdatesChan(u)

 for update := range updates {
  // 1. Handle regular messages and commands (e.g., /start)
  if update.Message != nil && update.Message.IsCommand() {
   if update.Message.Command() == "start" {
    // Send the initial welcome message with the keyboard
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Это бот-тарелка 🍽\nТвои очки: 0")
    msg.ReplyMarkup = createKeyboard()
    
    if _, err := bot.Send(msg); err != nil {
     log.Println("Error sending start message:", err)
    }
   }
  }

  // 2. Handle Inline button clicks (Callback queries)
  if update.CallbackQuery != nil {
   // Always answer the callback query to stop the loading spinner on the user's client
   callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
   if _, err := bot.Request(callback); err != nil {
    log.Println("Error answering callback:", err)
   }

   // Check if the specific "tap" button was clicked
   if update.CallbackQuery.Data == "tap" {
    userID := update.CallbackQuery.From.ID

    // Safely increment the score by locking the mutex
    mu.Lock()
    scores[userID]++
    currentScore := scores[userID]
    mu.Unlock()

    // Format the new message text with the updated score
    text := fmt.Sprintf("Бот-тарелка 🍽\nТвои очки: %d", currentScore)

    // Update the existing message inline instead of sending a new one
    editMsg := tgbotapi.NewEditMessageTextAndMarkup(
     update.CallbackQuery.Message.Chat.ID,
     update.CallbackQuery.Message.MessageID,
     text,
     createKeyboard(),
    )
    
    if _, err := bot.Send(editMsg); err != nil {
     log.Println("Error editing message:", err)
    }
   }
  }
 }
}

// createKeyboard creates an inline keyboard with a single "plate" button.
func createKeyboard() tgbotapi.InlineKeyboardMarkup {
 return tgbotapi.NewInlineKeyboardMarkup(
  tgbotapi.NewInlineKeyboardRow(
   tgbotapi.NewInlineKeyboardButtonData("🍽 Тапнуть!", "tap"),
  ),
 )
}