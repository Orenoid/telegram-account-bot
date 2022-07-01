package telebot

import (
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	billDAL "github.com/orenoid/telegram-account-bot/dal/bill"
	teleDAL "github.com/orenoid/telegram-account-bot/dal/telegram"
	"github.com/orenoid/telegram-account-bot/dal/user"
	"github.com/orenoid/telegram-account-bot/mock/telebotmock"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/orenoid/telegram-account-bot/service/bill"
	"github.com/orenoid/telegram-account-bot/service/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/telebot.v3"
	"reflect"
	"testing"
)

func TestNewHandlerHub(t *testing.T) {
	billService := &bill.Service{}
	teleService := &telegram.Service{}
	manager := NewInMemoryUserStateManager()
	hub := NewHandlerHub(billService, teleService, manager)
	assert.IsType(t, &HandlersHub{}, hub)
	assert.NotNil(t, hub)
	assert.True(t, hub.teleService == teleService)
	assert.True(t, hub.userStateManager == manager)
	assert.True(t, hub.billService == billService)
}

type HandlersHubTestSuite struct {
	suite.Suite

	teleMockCtrl *gomock.Controller
	teleRepo     *teleDAL.MockRepository
	billRepo     *billDAL.MockRepository
	userRepo     *user.MockRepository
	teleService  *telegram.Service
	billService  *bill.Service
	hub          *HandlersHub
}

func (suite *HandlersHubTestSuite) SetupTest(t *testing.T) func() {
	suite.teleMockCtrl = gomock.NewController(t)
	suite.teleRepo = teleDAL.NewMockRepository(suite.teleMockCtrl)
	suite.billRepo = billDAL.NewMockRepository(gomock.NewController(t))
	suite.userRepo = user.NewMockRepository(gomock.NewController(t))

	suite.teleService = telegram.NewService(suite.teleRepo)
	suite.billService = bill.NewService(suite.billRepo, suite.userRepo)

	userStateManager := NewInMemoryUserStateManager()
	suite.hub = NewHandlerHub(suite.billService, suite.teleService, userStateManager)

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

	type MethodParams struct {
		userID   int64
		userName string
		chatID   int64
	}
	var methodParams MethodParams
	gomonkey.ApplyMethod(reflect.TypeOf(suite.teleService), "CreateOrUpdateTelegramUser",
		func(service *telegram.Service, userID int64, userName string, chatID int64) (*models.TelegramUser, error) {
			methodParams = MethodParams{userID, userName, chatID}
			return nil, nil
		},
	)
	err := suite.hub.HandleStartCommand(ctx)
	suite.NoError(err)
	suite.Equal(MethodParams{417714530102, "JustARandomName", 592371906012}, methodParams)
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

func TestParseBill(t *testing.T) {
	category, name := ParseBill("饮食")
	assert.Equal(t, "饮食", category)
	assert.Nil(t, name)

	category, name = ParseBill("娱乐 欧卡2")
	assert.Equal(t, "娱乐", category)
	assert.NotNil(t, name)
	assert.Equal(t, "欧卡2", *name)

	category, name = ParseBill("娱乐 欧卡2 后面的字符串也当作 name 的一部分")
	assert.Equal(t, "娱乐", category)
	assert.NotNil(t, name)
	assert.Equal(t, "欧卡2 后面的字符串也当作 name 的一部分", *name)
}
