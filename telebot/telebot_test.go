package telebot

import (
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"gopkg.in/telebot.v3"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"testing"
)

func TestNewBot(t *testing.T) {
	inputSettings := telebot.Settings{Token: strconv.Itoa(rand.Int()), Offline: true}
	inputHub := &HandlersHub{}

	var telebotNewBotPref telebot.Settings
	telebotNewBotReturnedBot := &telebot.Bot{Token: strconv.Itoa(rand.Int())}
	patches := gomonkey.ApplyFunc(telebot.NewBot, func(pref telebot.Settings) (*telebot.Bot, error) {
		telebotNewBotPref = pref
		return telebotNewBotReturnedBot, nil
	})
	defer patches.Reset()

	var registerHandlersBot *telebot.Bot
	var registerHandlersHub *HandlersHub
	patches = gomonkey.ApplyFunc(RegisterHandlers, func(bot *telebot.Bot, hub *HandlersHub) {
		registerHandlersBot = bot
		registerHandlersHub = hub
	})
	defer patches.Reset()

	newBot, err := NewBot(inputSettings, inputHub)
	assert.Equal(t, inputSettings, telebotNewBotPref)
	assert.Equal(t, telebotNewBotReturnedBot, registerHandlersBot)
	assert.True(t, inputHub == registerHandlersHub)
	assert.Equal(t, telebotNewBotReturnedBot, newBot)
	assert.NoError(t, err)
}

func TestRegisterHandlers(t *testing.T) {
	var bot *telebot.Bot
	hub := &HandlersHub{}
	realRegistered := map[interface{}]telebot.HandlerFunc{}
	patches := gomonkey.ApplyMethod(reflect.TypeOf(bot), "Handle", func(_ *telebot.Bot, endpoint interface{}, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc) {
		realRegistered[endpoint] = h
	})
	expectedRegistered := map[interface{}]telebot.HandlerFunc{
		"/start":       hub.HandleStartCommand,
		"/day":         hub.HandleDayCommand,
		telebot.OnText: hub.HandleText,
	}
	defer patches.Reset()
	RegisterHandlers(bot, hub)

	for endpoint := range expectedRegistered {
		_, found := realRegistered[endpoint]
		assert.Truef(t, found, "expected endpoint: %v not registered", endpoint)
	}

	for endpoint, handler := range realRegistered {
		expectedHandler, found := expectedRegistered[endpoint]
		assert.Truef(t, found, "unknown endpoint registered: %v", endpoint)
		expectedName := runtime.FuncForPC(reflect.ValueOf(expectedHandler).Pointer()).Name()
		realName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		assert.Equal(t, expectedName, realName)
	}
}
