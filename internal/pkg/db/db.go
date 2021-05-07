package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/heimoshuiyu/gocc"
	_ "github.com/mattn/go-sqlite3"
)

var sqlCreateTable string = "CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY, messageText TEXT NOT NULL, time INTEGER NOT NULL);"
var sqlRecordMessage string = "INSERT INTO messages(messageText, time) VALUES (?, ?);"
var sqlDropMessageTable string = "DROP TABLE messages;"
var sqlGetAllMessages string = "SELECT (messageText) FROM messages;"

type DB struct {
	DatabaseName string
	SqlConn *sql.DB
	InitStmt *sql.Stmt
	RecordStmt *sql.Stmt
	DropStmt *sql.Stmt
	GetAllMessagesStmt *sql.Stmt
	T2s *gocc.OpenCC
}

func (db *DB) GetAllMessages() ([]string, error) {
	var err error
	var messages string
	resultSQL, err := db.GetAllMessagesStmt.Query()
	if err != nil {
		return nil, errors.New("Can not get all messages from DB at " + err.Error())
	}
	resultSlice := make([]string, 0)
	for resultSQL.Next() {
		err = resultSQL.Scan(&messages)
		if err != nil {
			return nil, errors.New("Can not read row message at " + err.Error())
		}
		resultSlice = append(resultSlice, messages)
	}

	return resultSlice, nil
}

func (db *DB) ResetMessages() (error) {
	var err error
	_, err = db.DropStmt.Exec()
	if err != nil {
		return errors.New("Can not reset messages at " + err.Error())
	}
	_, err = db.InitStmt.Exec()
	if err != nil {
		return errors.New("Can not re-init tables after reset at " + err.Error())
	}
	return nil
}

func (db *DB) RecordMessage(messageText string) (error) {
	translatedMessageText, err := db.T2s.Convert(messageText)
	if err != nil {
		log.Println("Can not translate", messageText, err.Error())
		translatedMessageText = messageText
	}
	db.RecordStmt.Exec(translatedMessageText, time.Now().Unix())
	return nil
}

func NewDB(DatabaseName string) (*DB, error) {
	var err error

	new_db := new(DB)

	new_db.DatabaseName = DatabaseName

	new_db.SqlConn, err = sql.Open("sqlite3", DatabaseName)
	if err != nil {
		return nil, errors.New("Can not open database at " + err.Error())
	}

	new_db.InitStmt, err = new_db.SqlConn.Prepare(sqlCreateTable)
	if err != nil {
		return nil, errors.New("Can not init InitStmt at " + err.Error())
	}

	_, err = new_db.InitStmt.Exec()
	if err != nil {
		return nil, errors.New("Faild exec InitStmt at " + err.Error())
	}

	new_db.RecordStmt, err = new_db.SqlConn.Prepare(sqlRecordMessage)
	if err != nil {
		return nil, errors.New("Can not init RecordStmt at " + err.Error())
	}

	new_db.DropStmt, err = new_db.SqlConn.Prepare(sqlDropMessageTable)
	if err != nil {
		return nil, errors.New("Can not init DropStmt at " + err.Error())
	}

	new_db.GetAllMessagesStmt, err = new_db.SqlConn.Prepare(sqlGetAllMessages)
	if err != nil {
		return nil, errors.New("Can not init GetAllMessagesStmt at " + err.Error())
	}

	new_db.T2s, err = gocc.New("t2s")
	if err != nil {
		return nil, errors.New("Can not init t2s at " + err.Error())
	}

	return new_db, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
