package services

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/nickname76/telegrambot"
)

type InstanceTelegram struct {
	TelegramBotRunning bool              // Indica si hay una instancia de telegram corriendo.
	InstanceBot        *telegrambot.API  // Apunta a una instancia de la API Telegram.
	Bot                *telegrambot.User // Apunta a user de la instancia de Telegram.
}
type SendMessageTelegramRequest struct {
	ChatID      int64  `json:"chat_id"`
	Message     string `json:"message"`
	ChatIDs     []int  `json:"chat_ids"`
	ProcessName string `json:"process_name"`
	SistemaID   int64  `json:"sistema_id"`
}

type TelegramService struct {
}

func NewTelegramService() *TelegramService {
	return &TelegramService{}
}

func (service *TelegramService) SendMessageTelegram(SendMessageTelegramRequest SendMessageTelegramRequest, instanceTelegram *InstanceTelegram) {
	if !instanceTelegram.TelegramBotRunning {
		fmt.Println("¡ No hay una instancia de Telegram ejecutandose !")
		fmt.Println("EJECUTANDO NUEVA INSTANCIA...")
		// Si por alguna razon no hay una instancia corriendo, se crea nuevamente la instancia.
		api, me, err := telegrambot.NewAPI(os.Getenv("TOKEN_BOT_TELEGRAM"))
		if err != nil {
			fmt.Println(err.Error())
		}
		instanceTelegram.TelegramBotRunning = true
		instanceTelegram.InstanceBot = api
		instanceTelegram.Bot = me
	}

	go func() {
		// Enviar la mensaje
		err := service.uploadConfigAndsendMessage(instanceTelegram.InstanceBot, instanceTelegram.Bot, SendMessageTelegramRequest)
		if err != nil {
			fmt.Println("Error:" + err.Error())
		}
	}()

}

func (service *TelegramService) uploadConfigAndsendMessage(api *telegrambot.API, me *telegrambot.User, requestSendMessage SendMessageTelegramRequest) (err error) {
	dtoIstanceTelegram := InstanceTelegram{
		TelegramBotRunning: true,
		InstanceBot:        api,
		Bot:                me,
	}

	if requestSendMessage.ChatID > 0 {
		chatId := telegrambot.ChatID(requestSendMessage.ChatID)
		err = service.sendMessage(requestSendMessage.Message, &chatId, &dtoIstanceTelegram)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Mensaje enviado con éxito a: ", chatId)
	}

	return
}

func (service *TelegramService) sendMessage(messsage string, chatId *telegrambot.ChatID, instancetelegram *InstanceTelegram) (err error) {
	_, err = instancetelegram.InstanceBot.SendMessage(&telegrambot.SendMessageParams{
		ChatID: chatId,
		Text:   messsage,
		ReplyMarkup: &telegrambot.ReplyKeyboardMarkup{
			Keyboard: [][]*telegrambot.KeyboardButton{{
				{
					Text: "Hello",
				},
			}},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		},
	})

	return
}

func (service *TelegramService) RunTelegramGoRoutine(instancetelegram *InstanceTelegram) {
	if instancetelegram.TelegramBotRunning {
		fmt.Println("Ya esta corriendo una instancia de telegram. (InstanceBot ejecutado)")
		fmt.Println("Ya esta corriendo una instancia de telegram.")
		return
	}

	api, me, err := telegrambot.NewAPI(os.Getenv("TOKEN_BOT_TELEGRAM"))
	if err != nil {
		fmt.Printf("\n\nError: %v\n\n", err)
	}
	instancetelegram.TelegramBotRunning = true
	instancetelegram.InstanceBot = api
	instancetelegram.Bot = me

	// Se reciben todos los mensajes que se envian al bot
	stop := telegrambot.StartReceivingUpdates(instancetelegram.InstanceBot, func(update *telegrambot.Update, err error) {
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		msg := update.Message
		if msg == nil {
			return
		}

		// Obtener la última parte del slice resultante. Éste indica el nombre del archivo que se quiere obtener los chats_id.

		// Si no está registrado y manda "start" se le envía el mensaje de bienvenida.
		if msg.Text == "/start" {

			message := "Es un placer saludarte desde BotCore 🤖. \n\n Tu ID es: " + strconv.Itoa(int(msg.Chat.ID))
			SendMessagesTelegram(message, &msg.Chat.ID, instancetelegram)

			return
		}

		// Si NO está registrado y no hay fallos, se procede a registrar a este usuario y asociar a dicho canal

		//SendMessagesTelegram("Api botcore corriendo correctamente.  🔔 ✅", &msg.Chat.ID, instancetelegram)

	})

	log.Printf("Started on %v", me.Username)

	message := "API REST BotCore corriendo 🤖🦾. "
	chatIDString := os.Getenv("CHAT_ID_TO_NOTIFY")
	var chatID telegrambot.ChatID
	chatIDToInt, err := strconv.Atoi(chatIDString)
	if err != nil {
		fmt.Print("error al convertir chatID: " + err.Error())
	}

	chatID = telegrambot.ChatID(chatIDToInt)
	SendMessagesTelegram(message, &chatID, instancetelegram)

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt)

	<-exitCh

	// Waits for all updates handling to complete
	stop()

}

func SendMessagesTelegram(messsage string, chatId *telegrambot.ChatID, instancetelegram *InstanceTelegram) (err error) {
	_, err = instancetelegram.InstanceBot.SendMessage(&telegrambot.SendMessageParams{
		ChatID: chatId,
		Text:   messsage,
		ReplyMarkup: &telegrambot.ReplyKeyboardMarkup{
			Keyboard: [][]*telegrambot.KeyboardButton{{
				{
					Text: "Hello",
				},
			}},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		},
	})

	return
}
