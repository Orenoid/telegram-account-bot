package telebot

import (
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"gopkg.in/telebot.v3"
	"math/rand"
	"reflect"
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
		"/help":             hub.HandleHelpCommand,
		"/start":            hub.HandleStartCommand,
		"/day":              hub.HandleDayCommand,
		"/month":            hub.HandleMonthCommand,
		"/cancel":           hub.HandleCancelCommand,
		"/set_keyboard":     hub.HandleSetKeyboardCommand,
		"/set_balance":      hub.HandleSetBalanceCommand,
		"/balance":          hub.HandleBalanceCommand,
		"/create_token":     hub.HandleCreateTokenCommand,
		telebot.OnText:      hub.HandleText,
		&prevDayBillBtnTmpl: hub.HandleDayBillSelectionCallback,
		&nextDayBillBtnTmpl: hub.HandleDayBillSelectionCallback,
		&prevMonthBtnTmpl:   hub.HandleMonthBillSelectionCallback,
		&nextMonthBtnTmpl:   hub.HandleMonthBillSelectionCallback,
		&cancelBillBtnTmpl:  hub.HandleCancelBillCallback,
	}
	defer patches.Reset()
	RegisterHandlers(bot, hub)

	for endpoint := range expectedRegistered {
		_, found := realRegistered[endpoint]
		assert.Truef(t, found, "expected endpoint: %v not registered", endpoint)
	}
}
