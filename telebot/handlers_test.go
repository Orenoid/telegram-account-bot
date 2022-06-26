package telebot

import (
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	teleDAL "github.com/orenoid/telegram-account-bot/dal/telegram"
	"github.com/orenoid/telegram-account-bot/mock/telebotmock"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/orenoid/telegram-account-bot/service/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/telebot.v3"
	"reflect"
	"testing"
)

func TestNewHandlerHub(t *testing.T) {
	teleService := &telegram.Service{}
	hub := NewHandlerHub(teleService)
	assert.IsType(t, &HandlersHub{}, hub)
	assert.NotNil(t, hub)
	assert.True(t, hub.teleService == teleService)
}

type HandlersHubTestSuite struct {
	suite.Suite

	teleMockCtrl *gomock.Controller
	teleRepo     *teleDAL.MockRepository
	teleService  *telegram.Service
	hub          *HandlersHub
}

func (suite *HandlersHubTestSuite) SetupTest(t *testing.T) func() {
	suite.teleMockCtrl = gomock.NewController(t)
	suite.teleRepo = teleDAL.NewMockRepository(suite.teleMockCtrl)
	suite.teleService = telegram.NewService(suite.teleRepo)
	suite.hub = NewHandlerHub(suite.teleService)

	return func() {
		suite.teleMockCtrl = nil
		suite.teleRepo = nil
		suite.teleService = nil
		suite.hub = nil
	}
}

func (suite *HandlersHubTestSuite) TestHandleStartCommand() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()

	ctx := telebotmock.NewMockContext(gomock.NewController(suite.T()))

	ctx.EXPECT().Chat().Return(&telebot.Chat{ID: 592371906012}).Times(1)
	ctx.EXPECT().Sender().Return(&telebot.User{ID: 417714530102, Username: "JustARandomName"}).Times(1)
	ctx.EXPECT().Send("hello")

	var methodParams = struct {
		userID   int64
		userName string
		chatID   int64
	}{}
	gomonkey.ApplyMethod(reflect.TypeOf(suite.teleService), "CreateOrUpdateTelegramUser",
		func(service *telegram.Service, userID int64, userName string, chatID int64) (*models.TelegramUser, error) {
			methodParams.userID = userID
			methodParams.userName = userName
			methodParams.chatID = chatID
			return nil, nil
		},
	)
	err := suite.hub.HandleStartCommand(ctx)
	suite.NoError(err)
	suite.Equal(
		struct {
			userID   int64
			userName string
			chatID   int64
		}{
			417714530102, "JustARandomName", 592371906012,
		}, methodParams)
}

func (suite *HandlersHubTestSuite) TestHandleStartCommandIfNilChat() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()

	ctx := telebotmock.NewMockContext(gomock.NewController(suite.T()))
	ctx.EXPECT().Chat().Return(nil)
	ctx.EXPECT().Sender().Times(0)

	err := suite.hub.HandleStartCommand(ctx)
	suite.ErrorContains(err, "nil chat of context")
}

func TestHandlersHub(t *testing.T) {
	suite.Run(t, new(HandlersHubTestSuite))
}
