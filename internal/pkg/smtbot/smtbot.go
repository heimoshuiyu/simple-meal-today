package smtbot

import (
	"errors"
	"log"
	"smt/internal/pkg/record"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type SmtBot struct {
	TgBotAPI *tgbotapi.BotAPI
	UpdatesChan tgbotapi.UpdatesChannel
	record *record.Record
}


func (smtbot *SmtBot) RecordPhoto(update tgbotapi.Update) (error) {
	var err error
	photoList := update.Message.Photo
	msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID,
		(*photoList)[0].FileID)
	msg.Caption = update.Message.Caption

	_, err = smtbot.TgBotAPI.Send(msg)
	return err
}

func (smtbot *SmtBot) Send(update tgbotapi.Update, text string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	smtbot.SendMessage(msg)
}

func (smtbot *SmtBot) SendMessage(msg tgbotapi.Chattable) {
	_, err := smtbot.TgBotAPI.Send(msg)
	if err != nil {
		log.Println("Send message failed ", err.Error())
	}
}

func (smtbot *SmtBot) ProcessCommand(update tgbotapi.Update) {

	// if no permission
	log.Println("registed admin:", smtbot.record.AdminUsersID, "chat from", update.Message.From.ID)
	if update.Message.From.ID != smtbot.record.AdminUsersID && update.Message.Chat.ID != smtbot.record.RegistedGroupID {
		smtbot.Send(update, "No permission, you are not admin or message is no in a registed chat")
		return
	}

	// regist group
	if update.Message.Command() == "registchat" && update.Message.From.ID == smtbot.record.AdminUsersID {
		smtbot.record.RegistedGroupID = update.Message.Chat.ID
		smtbot.Send(update, "【Chat注册成功】大家好，我是本群的吃饭睡觉提醒小助手，希望看见这条消息的群友可以和我一样，做个一天4顿饭+睡20个小时的five吧~\n【使用方法】在群内发送 /register 命令进行注册，然后把 #朴素一餐 的照片通过私聊发给我就可以啦~")
		return
	}

	if update.Message.Command() == "register" && update.Message.Chat.ID == smtbot.record.RegistedGroupID {
		smtbot.record.AddRegistedUser(update.Message.From.ID)
		smtbot.Send(update, "可爱的【" + update.Message.From.FirstName + "】你已经成功注册")
		return
	}

	if update.Message.Command() == "save" && update.Message.From.ID == smtbot.record.AdminUsersID {
		smtbot.record.Save()
		return
	}

}


func (smtbot *SmtBot) ProcessPrivateMessage(update tgbotapi.Update) {
	// check registed group
	if smtbot.record.RegistedGroupID == 0 {
		if smtbot.TgBotAPI.Debug {
			smtbot.Send(update, "Do not have registed chat currently")
		}
		log.Println("Not have registed chat but recived private message")
		return
	}
	// not a registed user
	if !smtbot.record.IsRegistedUser(update.Message.From.ID) {
		if smtbot.TgBotAPI.Debug {
			smtbot.Send(update, "permission deny, you are not registed")
		}
		log.Println("Not registed user message from", update.Message.From.ID, update.Message.From.FirstName, update.Message.Text)
		return
	}

	// no photo
	photoList := update.Message.Photo
	if photoList == nil {
		smtbot.Send(update, "no photo in message!")
		return
	}

	// record photo

	// forward photo to chat
	msg := tgbotapi.NewPhotoShare(smtbot.record.RegistedGroupID,
		(*photoList)[0].FileID)
	msg.Caption = update.Message.Caption
	smtbot.SendMessage(msg)
}

func (smtbot *SmtBot) ProcessGroupMessage(update tgbotapi.Update) {
}

func (smtbot *SmtBot) Run() (error) {
	for update := range smtbot.UpdatesChan{
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			smtbot.ProcessCommand(update)
			continue
		}

		if update.Message.Chat.IsPrivate() {
			smtbot.ProcessPrivateMessage(update)
			continue
		}

		if update.Message.Chat.IsGroup() {
			smtbot.ProcessGroupMessage(update)
			continue
		}

		if smtbot.TgBotAPI.Debug {
			smtbot.Send(update, "invaild message")
		}
	}
	return nil
}

func NewSmtBot(token string, debug bool, timeout int, recordFile string) (*SmtBot, error) {
	var err error
	smtbot := new(SmtBot)

	// load record
	smtbot.record = record.NewRecord(recordFile)

	smtbot.TgBotAPI, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.New("Failed at creating tgbotapi " + err.Error())
	}

	smtbot.TgBotAPI.Debug = debug

	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout
	smtbot.UpdatesChan, err = smtbot.TgBotAPI.GetUpdatesChan(u)
	if err != nil {
		return nil, errors.New("Failed at creating updates channel " + err.Error())
	}

	return smtbot, nil
}
