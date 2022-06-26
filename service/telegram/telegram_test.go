package telegram

import (
	"github.com/golang/mock/gomock"
	"github.com/orenoid/telegram-account-bot/dal/bill"
	"github.com/orenoid/telegram-account-bot/dal/telegram"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"strconv"
	"testing"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	teleRepo := telegram.NewMockRepository(ctrl)

	service := NewService(teleRepo)
	assert.NotNil(t, service)
	assert.True(t, service.teleRepo == teleRepo)
}

type ServiceTestSuite struct {
	suite.Suite

	teleMockCtrl *gomock.Controller
	teleRepo     *telegram.MockRepository
	billMockCtrl *gomock.Controller
	billRepo     *bill.MockRepository
	teleService  *Service
}

func (suite *ServiceTestSuite) SetupTest(t *testing.T) func() {
	suite.teleMockCtrl = gomock.NewController(t)
	suite.teleRepo = telegram.NewMockRepository(suite.teleMockCtrl)
	suite.billMockCtrl = gomock.NewController(t)
	suite.billRepo = bill.NewMockRepository(suite.billMockCtrl)
	suite.teleService = NewService(suite.teleRepo)

	return func() {
		suite.teleMockCtrl = nil
		suite.teleRepo = nil
		suite.billMockCtrl = nil
		suite.billRepo = nil
		suite.teleService = nil
	}
}

func (suite *ServiceTestSuite) TestCreateOrUpdateUser() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()

	userID := int64(rand.Int())
	userName := strconv.Itoa(rand.Int())
	chatID := int64(rand.Int())

	// userFromRepo 重新使用新的随机值，是为了校验 service 返回的 user 是由 repo 提供的
	userFromRepo := &models.TelegramUser{BaseUserID: uint(rand.Int()), UserName: "whatever", ChatID: int64(rand.Int())}
	suite.teleRepo.EXPECT().CreateOrUpdateTelegramUser(userID, userName, chatID).Return(userFromRepo, nil)

	userFromService, err := suite.teleService.CreateOrUpdateTelegramUser(userID, userName, chatID)
	suite.NoError(err)
	suite.Equal(userFromRepo, userFromService)
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
