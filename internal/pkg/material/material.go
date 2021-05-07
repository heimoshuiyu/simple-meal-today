package material

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"time"
)

type Material struct {
	FoodJsonFilePath string
	Food []string
}

func NewMaterial(foodFilePath string) (*Material, error) {
	rand.Seed(time.Now().Unix())
	material := new(Material)
	material.FoodJsonFilePath = foodFilePath
	file, err := os.Open(material.FoodJsonFilePath)
	if err != nil {
		return nil, errors.New("Can not open food material file " + err.Error())
	}
	err = json.NewDecoder(file).Decode(&material.Food)
	if err != nil {
		return nil, errors.New("Can not decode food json " + err.Error())
	}
	return material, nil
}

func (material *Material) GetFood() (string) {
	return material.Food[rand.Int() % (len(material.Food) - 1)]
}
