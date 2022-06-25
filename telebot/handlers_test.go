package telebot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHandlerHub(t *testing.T) {
	hub := NewHandlerHub()
	assert.NotNil(t, hub)
	assert.IsType(t, &HandlersHub{}, hub)
}
