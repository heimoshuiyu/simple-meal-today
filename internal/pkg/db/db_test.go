package db

import (
	"testing"
)

func TestDB(t *testing.T) {
	new_db, err := NewDB("test.sqlite")
	if err != nil {
		t.Fatal("Can not create new_db at " + err.Error())
	}

	err = new_db.RecordMessage("Hello")
	if err != nil {
		t.Fatal("Can not record message at " + err.Error())
	}

	err = new_db.ResetMessages()
	if err != nil {
		t.Fatal("Can not reset table at " + err.Error())
	}

	err = new_db.RecordMessage("Hello !!!")
	if err != nil {
		t.Fatal("Can not record message at " + err.Error())
	}

	err = new_db.RecordMessage("Hi !!!!!!!")
	if err != nil {
		t.Fatal("Can not record message at " + err.Error())
	}

	result, err := new_db.GetAllMessages()
	if err != nil {
		t.Fatal("Can not get all messages at " + err.Error())
	}
	t.Log(result)
}
