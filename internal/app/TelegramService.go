package app

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	Consts "password-storage/internal/pkg/consts"
	"reflect"
	"strings"
)

const telegramToken = "5556005985:AAEkYXaFm40VbyiObxfNP25dAS0v0qtLXWA"

const (
	start              = "/start"
	login              = "/login"
	getAllPasswordData = "/getallpassworddata"
	empty              = ""
)

func TelegramBot() {

	//Создаем бота
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		panic(err)
	}

	//Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//Получаем обновления от бота
	updates, err := bot.GetUpdatesChan(u)

	//Последняя операция по которой был задан вопрос и требуется продолжение
	lastOperation := empty
	gotLogin := ""
	gotPassword := ""
	for update := range updates {
		if update.Message == nil {
			continue
		}

		//Проверяем что от пользователья пришло именно текстовое сообщение
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			operationToProcess := update.Message.Text
			if lastOperation != empty && update.Message.Text != start {
				operationToProcess = lastOperation
			}

			switch operationToProcess {
			case start:

				lastOperation = empty

				sendMessage(update, bot, "Hi am KeyValue Storage, please /login")

			case login:

				//ищем сессию, если нет то заставляем логиниттся
				if lastOperation == empty && IsSessionValidByTelegramId(update.Message.Chat.ID) {
					sendMessage(update, bot, "Hi, I know you!")
				} else {
					// проверяем логин пароль пользователя
					if lastOperation == empty && gotLogin == "" {
						sendMessage(update, bot, "You have to send login")
						lastOperation = login
						continue
					}
					if lastOperation == login && gotLogin == "" {
						gotLogin = update.Message.Text
						sendMessage(update, bot, "You have to send password")
						continue
					}
					if lastOperation == login && gotLogin != "" {
						gotPassword = update.Message.Text
					}

					if gotLogin == Consts.DefaultUser && gotPassword == Consts.DefaultPassword {
						CreateNewSessionWithTelegramId(update.Message.Chat.ID)
						sendMessage(update, bot, "Hi I know you!")
						lastOperation = empty
					}
				}

			case getAllPasswordData:

				if IsSessionValidByTelegramId(update.Message.Chat.ID) {
					data := GetAllPasswordData()
					builder := strings.Builder{}
					for i := range data {
						builder.Write([]byte(" Логин: "))
						builder.Write([]byte(*data[i][0].Key))
						builder.Write([]byte(" Пароль: "))
						builder.Write([]byte(*data[i][0].Value))
					}
					sendMessage(update, bot, builder.String())
				} else {
					sendMessage(update, bot, "May be your session expired, use /login to get access")
				}

			default:
				sendMessage(update, bot, "I didn't understand your command BRO!")
			}
		} else {
			sendMessage(update, bot, "I didn't understand your command BRO!")
		}
	}
}

func sendMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI, text string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Printf("Не удалось отправить сообщение в chatId = %v\n", update.Message.Chat.ID)
	}
}
