package record

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TesNewtRecord(t *testing.T) {
	new_record := new(Record)
	new_record.RegistedUsersID = make(map[int]bool)
}

func TestSaveAndLoadFile(t *testing.T) {
	r := Record{
		AdminUsersID: 1,
		RegistedGroupID: -100,
		RegistedUsersID: make(map[int]bool),
	}
	f, err := os.Create("config_before.json")
	if err != nil {
		t.Fatal("Can not create original config file")
	}
	json.NewEncoder(f).Encode(&r)
	record := NewRecord("config_before")
	record.Save("config_before.json")
	record.AdminUsersID = 39
	record.Save("config_after.json")
	record.Load("config_before.json")
	if record.AdminUsersID != 0 {
		t.Fatal("AdminUsersID after load not match!")
	}
}

func TestRegistedUsers(t *testing.T) {
	r := Record{
		AdminUsersID: 1,
		RegistedGroupID: -100,
		RegistedUsersID: make(map[int]bool),
	}
	f, err := os.Create("config_before.json")
	if err != nil {
		t.Fatal("Can not create original config file")
	}
	json.NewEncoder(f).Encode(&r)
	record := NewRecord("config_before.json")
	if record.IsRegistedUser(0) {
		t.Fatal("0 should not be registed user")
	}
	record.AddRegistedUser(39)
	if !record.IsRegistedUser(39) {
		fmt.Println(record.RegistedUsersID)
		t.Fatal("39 should be registed user")
	}
	record.Save("config_after.json")
}
