package smtbot

import (
	"errors"
	"fmt"
	"log"
	"smt/internal/pkg/ans"
	"smt/internal/pkg/db"
	"smt/internal/pkg/record"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type SmtBot struct {
	TgBotAPI *tgbotapi.BotAPI
	UpdatesChan tgbotapi.UpdatesChannel
	record *record.Record
	db *db.DB
	ans *ans.Ans
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
	var err error

	// if no permission
	log.Println("registed admin:", smtbot.record.AdminUsersID, "chat from", update.Message.From.ID)
	if update.Message.From.ID != smtbot.record.AdminUsersID && update.Message.Chat.ID != smtbot.record.RegistedGroupID {
		smtbot.Send(update, "No permission, you are not admin or message is no in a registed chat")
		return
	}

	// regist group
	if update.Message.Command() == "registerchat" && update.Message.From.ID == smtbot.record.AdminUsersID {
		smtbot.record.RegistedGroupID = update.Message.Chat.ID
		smtbot.Send(update, "【Chat注册成功】大家好，我是本群的吃饭睡觉提醒小助手，希望看见这条消息的群友可以和我一样，做个一天4顿饭+睡20个小时的five吧~\n【使用方法】在群内发送 /register 命令进行注册，然后把 #朴素一餐 的照片通过私聊发给我就可以啦~")
		return
	}

	if update.Message.Command() == "register" && update.Message.Chat.ID == smtbot.record.RegistedGroupID {
		if smtbot.record.IsRegistedUser(update.Message.From.ID) {
			smtbot.Send(update, "歪？" + update.Message.From.FirstName + "你注册过了诶？")
			return
		}
		smtbot.record.AddRegistedUser(update.Message.From.ID)
		smtbot.Send(update, "可爱的【" + update.Message.From.FirstName + "】你已经成功注册")
		return
	}

	if update.Message.Command() == "save" && smtbot.record.IsRegistedUser(update.Message.From.ID) {
		err = smtbot.record.Save()
		if err != nil {
			smtbot.Send(update, "保存操作执行失败：" + err.Error())
			return
		}
		smtbot.Send(update, "保存配置完成")
		return
	}

	if update.Message.Command() == "words" && smtbot.record.IsRegistedUser(update.Message.From.ID) {
		allMessages, err := smtbot.db.GetAllMessages()
		numOfAllMessages := len(allMessages)
		if err != nil {
			smtbot.Send(update, "获取词频错误：" + err.Error())
			return
		}
		wordCounts := smtbot.ans.CalcWordCounts(allMessages)
		numofAllWords := 0
		for _, v := range wordCounts {
			numofAllWords += v
		}
		words := smtbot.ans.CalcDailyWordsTrend(wordCounts)
		wordlist := strings.Join(words, "，")
		msgText := fmt.Sprintf(
			"今日统计：\n消息数量：%d，词条数量：%d\n今日关键词：%s",
			numOfAllMessages,
			numofAllWords,
			wordlist,
		)
		smtbot.Send(update, msgText)
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
	var err error

	// ignore all not registed chat message
	if update.Message.Chat.ID != smtbot.record.RegistedGroupID {
		return
	}

	if update.Message.Text != "" {
		err = smtbot.db.RecordMessage(update.Message.Text)
		if err != nil {
			log.Println("Failed at record messageText to db at " + err.Error())
		}
	}
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

		if update.Message.Chat.IsSuperGroup() {
			smtbot.ProcessGroupMessage(update)
			continue
		}

		if smtbot.TgBotAPI.Debug {
			smtbot.Send(update, "invaild message")
		}
	}
	return nil
}

func NewSmtBot(token string, adminUserId int, debug bool, timeout int, recordFile string, databaseName string) (*SmtBot, error) {
	var err error
	smtbot := new(SmtBot)

	// load record
	smtbot.record = record.NewRecord(recordFile, adminUserId)

	// load database
	smtbot.db, err = db.NewDB(databaseName)
	if err != nil {
		log.Fatal("Can not create new db at " + err.Error())
	}

	// load analysis
	smtbot.ans = ans.NewAns()

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
