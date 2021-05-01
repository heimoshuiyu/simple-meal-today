package smtbot

import "testing"

func TestNewSmtBot(t *testing.T) {
	_, err := NewSmtBot("585292276:AAH1Dlj65t3HU4qXGQ_6oqbnTUa1cXH7rj8")
	if err != nil {
		t.Fatal("Create bot failed")
	}
}
