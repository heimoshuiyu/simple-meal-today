package smtbot

import (
	"sync"
	"time"

)

var sleepSticker string = "CAACAgUAAxkBAAK5H2CVQcBquJ9LqKGCr9rAYWj4wjCfAAI7AAMtPtckaxkpI1cDyakfBA"
var learningSticker string = "CAACAgUAAxkBAAK5I2CVRReCvK-ZAoDtjmKFUfkjbenHAAI9AAMtPtcksc1z1vy5AAFuHwQ"
var shuayaSticker string = "CAACAgUAAxkBAAK5J2CVRcLzkEpjyHLjKElAowvSTLqEAAI8AAMtPtckPoNCwGCtGzEfBA"
var guguVideo string = "BAACAgUAAxkBAAK5K2CVSoNnugu991BZ4gK4ld3u3Dj_AAL9AQACy-nxVzyxRBUpFzDQHwQ"

func (smtbot *SmtBot) StartTaskWorkers() {
	smtbot.wg = new(sync.WaitGroup)
	smtbot.wg.Add(1)
	go smtbot.MainWorker()
}

func (smtbot *SmtBot) MainWorker() {
	for {
		time.Sleep(time.Second)
		hour, minute, _ := time.Now().Clock()
		if hour == 8-8 && minute == 0 {
			smtbot.GoodMorning()
			time.Sleep(time.Second * 61)
		} else if hour == 23-8 && minute == 0 {
			smtbot.GoodEvening()
			time.Sleep(time.Second * 61)
		} else if hour == 22-8 && minute == 30 {
			smtbot.RemindToLearn()
			time.Sleep(time.Second * 61)
		} else if hour == 22-8 && minute == 45 {
			smtbot.PrepareToSleep()
			time.Sleep(time.Second * 61)
		}
	}
}

func (smtbot *SmtBot) RemindToLearn() {
	smtbot.SendStickerToID(smtbot.record.RegistedGroupID, learningSticker)
}

func (smtbot *SmtBot) PrepareToSleep() {
	smtbot.SendStickerToID(smtbot.record.RegistedGroupID, shuayaSticker)
	smtbot.SendVideoToID(smtbot.record.RegistedGroupID, guguVideo)
}

func (smtbot *SmtBot) GoodEvening() {
	smtbot.SendStickerToID(smtbot.record.RegistedGroupID, sleepSticker)
}

func (smtbot *SmtBot) GoodMorning() {
	var err error
	var databaseResetStatus string

	smtbot.ReportWords(smtbot.record.RegistedGroupID)

	err = smtbot.db.ResetMessages()
	if err != nil {
		databaseResetStatus = "数据库：重置失败：" + err.Error()
	} else {
		databaseResetStatus = "数据库：重置成功哒"
	}

	text := `起～～床！！！！
	早上要认真刷牙洗脸，来跟[亚托莉]一起咕噜咕噜~呸♡`
	text = text + "\n" + databaseResetStatus
	smtbot.SendToID(smtbot.record.RegistedGroupID, text)

	smtbot.SendStickerToID(smtbot.record.RegistedGroupID, shuayaSticker)
	smtbot.SendVideoToID(smtbot.record.RegistedGroupID, guguVideo)
}
