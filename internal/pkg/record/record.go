package record

import (
	"encoding/json"
	"os"
)

type Record struct {
	AdminUsersID int
	RegistedGroupID int64
	RegistedUsersID []int
}

func NewRecord(recordFile string) (*Record) {
	new_record := new(Record)
	new_record.Load(recordFile)
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

func (r *Record) Save(filename string) (error) {
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
	r.RegistedUsersID = append(r.RegistedUsersID, userId)
}

func (r *Record) IsRegistedUser(userId int) (bool) {
	for _, i := range r.RegistedUsersID {
		if userId == i {
			return true
		}
	}
	return false
}
