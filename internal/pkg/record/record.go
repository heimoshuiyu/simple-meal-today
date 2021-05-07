package record

import (
	"encoding/json"
	"os"
)

type Record struct {
	AdminUsersID int
	RegistedGroupID int64
	RegistedUsersID map[int]bool
	recordFile string
}

func NewRecord(recordFile string, adminUserId int) (*Record) {
	new_record := new(Record)
	new_record.Load(recordFile)
	new_record.recordFile = recordFile
	if adminUserId != 0 {
		new_record.AdminUsersID = adminUserId
	}
	return new_record
}

func (r *Record) Load(filename string) (error) {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(r)
	if err != nil {
		return err
	}

	return nil
}

func (r *Record) Save() (error) {
	return r.SaveAs(r.recordFile)
}

func (r *Record) SaveAs(filename string) (error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(r)
	if err != nil {
		return err
	}

	return nil
}

func (r *Record) AddRegistedUser(userId int) {
	r.RegistedUsersID[userId] = true
}

func (r *Record) IsRegistedUser(userId int) (bool) {
	_, ok := r.RegistedUsersID[userId]
	return ok
}
