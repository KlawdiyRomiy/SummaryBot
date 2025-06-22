package bot

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	ai "summary_bot/ai"
	"summary_bot/files"
)

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	if update.Message.IsCommand() {
		handleCommand(bot, update.Message)
		return
	}
	if update.Message.Text != "" {
		handleText(bot, update.Message)
		return
	}
	if update.Message.Document != nil {
		handleDocument(bot, update.Message)
		return
	}
}
func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(message.Chat.ID, "Здравствуй! Я бот, который делает выжимку из текста или файла, который ты мне пришлешь. Напиши /help для пояснительной бригады. ")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
	case "help":
		msg := tgbotapi.NewMessage(message.Chat.ID, "Просто отправь мне текст до 4000 символов или файл формата DOCX весом до 20 мегабайт, и я сделаю выжимку. ")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Зачем ты лишнюю команду отправил? Либо файл или текст отправь, либо трусы надень.")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
	}
}

func handleText(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	text := message.Text
	if len(text) > 4000 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Текст должен быть короче 4000 символов.")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}
	summary, err := ai.SummarizeText(text)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при выжимке текста")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, summary)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func handleDocument(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	const maxFileSize = 20 * 1024 * 1024 // 20 мегабайт
	doc := message.Document

	if doc.FileSize > maxFileSize {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Брат, мне нужен файл поменьше. У меня тоже есть лимиты и это 20 мегабайт.")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	ext := filepath.Ext(doc.FileName)
	if ext != ".pdf" && ext != ".docx" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неподдерживаемый формат документа. Я хочу только docx.")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	fileConfig := tgbotapi.FileConfig{FileID: doc.FileID}
	file, err := bot.GetFile(fileConfig)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при получении файла")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	url := file.Link(bot.Token)
	tempPath := "./temp" + ext

	err = downloadFile(url, tempPath)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при загрузке файла")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}
	defer os.Remove(tempPath)

	var text string
	if ext == ".pdf" {
		text, err = files.ExtractPDF(tempPath)
	} else {
		text, err = files.ExtractDocx(tempPath)
	}

	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при извлечении текста")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	summary, err := ai.SummarizeText(text)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при выжимке текста")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, summary)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
