package record

import (
	"fmt"
	"testing"
)

func TesNewtRecord(t *testing.T) {
	new_record := new(Record)
	new_record.RegistedUsersID = make([]int, 0)
}

func TestSaveAndLoadFile(t *testing.T) {
	record := NewRecord()
	record.Save("config_before.json")
	record.AdminUsersID = 39
	record.Save("config_after.json")
	record.Load("config_before.json")
	if record.AdminUsersID != 0 {
		t.Fatal("AdminUsersID after load not match!")
	}
}

func TestRegistedUsers(t *testing.T) {
	record := NewRecord()
	if record.IsRegistedUser(0) {
		t.Fatal("0 should not be registed user")
	}
	record.AddRegistedUser(39)
	if !record.IsRegistedUser(39) {
		fmt.Println(record.RegistedUsersID)
		t.Fatal("39 should be registed user")
	}
}
