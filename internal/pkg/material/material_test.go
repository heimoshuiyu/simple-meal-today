package material

import "testing"

func TestMaterial(t *testing.T) {
	var err error
	material, err := NewMaterial("../../../cmd/simple-meal-today/food.json")
	if err != nil {
		t.Fatal(err)
	}
	for i:=0; i<10; i++ {
		t.Log(material.GetFood())
	}
}
