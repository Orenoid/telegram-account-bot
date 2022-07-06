package telebot

import (
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	billDAL "github.com/orenoid/telegram-account-bot/dal/bill"
	teleDAL "github.com/orenoid/telegram-account-bot/dal/telegram"
	"github.com/orenoid/telegram-account-bot/dal/user"
	"github.com/orenoid/telegram-account-bot/mock/telebotmock"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/orenoid/telegram-account-bot/service/bill"
	"github.com/orenoid/telegram-account-bot/service/telegram"
	"github.com/orenoid/telegram-account-bot/utils/strings"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/telebot.v3"
	"reflect"
	"testing"
	"time"
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

	teleRepo *teleDAL.MockRepository
	billRepo *billDAL.MockRepository
	userRepo *user.MockRepository

	teleService *telegram.Service
	billService *bill.Service

	userStateManager *MockUserStateManager

	hub *HandlersHub
}

func (suite *HandlersHubTestSuite) SetupTest(t *testing.T) func() {
	suite.teleRepo = teleDAL.NewMockRepository(gomock.NewController(t))
	suite.billRepo = billDAL.NewMockRepository(gomock.NewController(t))
	suite.userRepo = user.NewMockRepository(gomock.NewController(t))

	suite.teleService = telegram.NewService(suite.teleRepo)
	suite.billService = bill.NewService(suite.billRepo, suite.userRepo)
	suite.userStateManager = NewMockUserStateManager(gomock.NewController(t))

	suite.hub = NewHandlerHub(suite.billService, suite.teleService, suite.userStateManager)

	return func() {
		suite.teleRepo = nil
		suite.billRepo = nil
		suite.userRepo = nil
		suite.teleService = nil
		suite.billService = nil
		suite.userStateManager = nil
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

func (suite *HandlersHubTestSuite) TestHandleText() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()

	userStates := []*UserState{
		{Type: Empty}, {Type: CreatingBill},
	}

	type OnEmptyPatchesHelper struct {
		called    bool
		paramCtx  telebot.Context
		returnErr error
	}
	onEmpty := &OnEmptyPatchesHelper{returnErr: errors.New("OnEmptyError")}
	patches := gomonkey.ApplyMethod(reflect.TypeOf(suite.hub), "OnEmpty", func(_ *HandlersHub, ctx telebot.Context) error {
		onEmpty.called = true
		onEmpty.paramCtx = ctx
		return onEmpty.returnErr
	})
	defer patches.Reset()

	type OnCreatingBill struct {
		called     bool
		paramCtx   telebot.Context
		paramState *UserState
		returnErr  error
	}

	onCreatingBill := &OnCreatingBill{returnErr: errors.New("OnCreatingBillError")}
	patches = gomonkey.ApplyMethod(reflect.TypeOf(suite.hub), "OnCreatingBill",
		func(_ *HandlersHub, ctx telebot.Context, state *UserState) error {
			onCreatingBill.called = true
			onCreatingBill.paramState = state
			onCreatingBill.paramCtx = ctx
			return onCreatingBill.returnErr
		})
	defer patches.Reset()

	for _, state := range userStates {
		ctx := telebotmock.NewMockContext(gomock.NewController(suite.T()))
		ctx.EXPECT().Sender().Return(&telebot.User{ID: 42}).Times(1)
		suite.userStateManager.EXPECT().GetUserState(int64(42)).Return(state, nil).Times(1)

		err := suite.hub.HandleText(ctx)

		switch state.Type {
		case Empty:
			suite.True(onEmpty.called)
			suite.False(onCreatingBill.called)
			suite.True(onEmpty.paramCtx == ctx)
			suite.Equal(onEmpty.returnErr, err)
		case CreatingBill:
			suite.True(onCreatingBill.called)
			suite.False(onEmpty.called)
			suite.True(state == onCreatingBill.paramState)
			suite.True(onCreatingBill.paramCtx == ctx)
			suite.Equal(onCreatingBill.returnErr, err)
		default:
			suite.NoError(err)
		}
		onEmpty.called = false
		onEmpty.paramCtx = nil
		onCreatingBill.called = false
		onCreatingBill.paramCtx = nil
		onCreatingBill.paramState = nil
	}

}

func (suite *HandlersHubTestSuite) TestOnEmpty() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()

	type ParseBillPatchesHelper struct {
		called bool
		// params
		text string
		// return
		category string
		name     *string
	}
	parseBillPatchesHelper := &ParseBillPatchesHelper{category: "饮食"}
	patches := gomonkey.ApplyFunc(ParseBill, func(text string) (string, *string) {
		parseBillPatchesHelper.called = true
		parseBillPatchesHelper.text = text
		return parseBillPatchesHelper.category, parseBillPatchesHelper.name
	})
	defer patches.Reset()

	ctx := telebotmock.NewMockContext(gomock.NewController(suite.T()))
	ctx.EXPECT().Text().Return("饮食").Times(1)
	ctx.EXPECT().Sender().Return(&telebot.User{ID: 6379})
	suite.userStateManager.EXPECT().SetUserState(int64(6379),
		&UserState{CreatingBill, &parseBillPatchesHelper.category, parseBillPatchesHelper.name}).Times(1)
	ctx.EXPECT().Send(fmt.Sprintf("账单类别：%s，请输入账单金额", parseBillPatchesHelper.category)).Times(1)

	err := suite.hub.OnEmpty(ctx)

	suite.Equal("饮食", parseBillPatchesHelper.text)
	suite.True(parseBillPatchesHelper.called)
	suite.NoError(err)
}

func (suite *HandlersHubTestSuite) TestOnEmptyIfEmptyText() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()

	ctx := telebotmock.NewMockContext(gomock.NewController(suite.T()))
	ctx.EXPECT().Text().Return("").Times(1)
	suite.userStateManager.EXPECT().SetUserState(0, nil).Times(0)
	err := suite.hub.OnEmpty(ctx)
	suite.Nil(err)
}

func (suite *HandlersHubTestSuite) TestOnCreatingBill() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()

	// mock teleService.GetBaseUserID
	type GetBaseUserID struct {
		called           bool
		paramTeleUserID  int64
		returnBaseUserID uint
		returnErr        error
	}
	getBaseUserID := &GetBaseUserID{returnBaseUserID: 255, returnErr: nil}
	patches := gomonkey.ApplyMethod(reflect.TypeOf(suite.hub.teleService), "GetBaseUserID",
		func(_ *telegram.Service, teleUserID int64) (uint, error) {
			getBaseUserID.called = true
			getBaseUserID.paramTeleUserID = teleUserID
			return getBaseUserID.returnBaseUserID, getBaseUserID.returnErr
		})
	defer patches.Reset()

	// mock billService.CreateNewBill
	type CreateNewBill struct {
		called        bool
		paramUserID   uint
		paramAmount   float64
		paramCategory string
		paramOpts     []billDAL.CreateBillOptions
		returnBill    *models.Bill
		returnErr     error
	}
	createNewBill := &CreateNewBill{returnBill: &models.Bill{}, returnErr: nil}
	patches = gomonkey.ApplyMethod(reflect.TypeOf(suite.billService), "CreateNewBill",
		func(_ *bill.Service, userID uint, amount float64, category string, opts ...billDAL.CreateBillOptions) (*models.Bill, error) {
			createNewBill.called = true
			createNewBill.paramUserID = userID
			createNewBill.paramAmount = amount
			createNewBill.paramCategory = category
			createNewBill.paramOpts = opts
			return createNewBill.returnBill, createNewBill.returnErr
		})
	defer patches.Reset()

	ctx := telebotmock.NewMockContext(gomock.NewController(suite.T()))
	userState := &UserState{BillCategory: strings.Pointer("娱乐"), BillName: strings.Pointer("龙珠漫画")}
	ctx.EXPECT().Text().Return("23.14")

	ctx.EXPECT().Sender().Return(&telebot.User{ID: 7})
	defer func() { suite.True(getBaseUserID.called) }()
	defer func() { suite.Equal(int64(7), getBaseUserID.paramTeleUserID) }()

	defer func() { suite.True(createNewBill.called) }()
	defer func() { suite.Equal(getBaseUserID.returnBaseUserID, createNewBill.paramUserID) }()
	defer func() { suite.Equal(-23.14, createNewBill.paramAmount) }()
	defer func() { suite.Equal("娱乐", createNewBill.paramCategory) }()
	defer func() { suite.Len(createNewBill.paramOpts, 1) }()
	defer func() { suite.Equal("龙珠漫画", *createNewBill.paramOpts[0].Name) }()

	suite.userStateManager.EXPECT().ClearUserState(int64(7)).Return(nil).Times(1)
	ctx.EXPECT().Send(&NewBillSender{createNewBill.returnBill})

	err := suite.hub.OnCreatingBill(ctx, userState)
	suite.NoError(err)
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

func TestParseAmount(t *testing.T) {
	testCases := map[string]float64{
		"+1.23": 1.23,
		"1.23":  -1.23,
		".2":    -0.2,
	}
	for text, expectedAmount := range testCases {
		amount, err := parseAmount(text)
		assert.NoError(t, err)
		assert.Equal(t, expectedAmount, amount)
	}

	badCases := []string{"", "+-1", "abc", "++2"}
	for _, text := range badCases {
		_, err := parseAmount(text)
		assert.Error(t, err)
	}
}

func TestGetDayRange(t *testing.T) {
	testTime := time.Date(2011, 1, 2, 9, 5, 7, 1234, time.UTC)
	begin, end := getDayRange(testTime)
	assert.Equal(t, time.Date(2011, 1, 2, 0, 0, 0, 0, time.UTC), begin)
	assert.Equal(t, time.Date(2011, 1, 3, 0, 0, 0, 0, time.UTC), end)

	testTime = time.Date(2011, 2, 28, 0, 0, 0, 0, time.UTC)
	begin, end = getDayRange(testTime)
	assert.Equal(t, time.Date(2011, 2, 28, 0, 0, 0, 0, time.UTC), begin)
	assert.Equal(t, time.Date(2011, 3, 1, 0, 0, 0, 0, time.UTC), end)
}
