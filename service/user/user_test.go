package user

import (
	"github.com/golang/mock/gomock"
	"github.com/orenoid/account-bot/dal/user"
	"github.com/orenoid/account-bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"math/rand"
	"testing"
	"time"
)

func TestNewUserService(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ur := user.NewMockRepository(ctrl)
	us, err := NewUserService(ur)

	assert.NoError(t, err)
	assert.True(t, ur == us.userRepo)
}

type UserServiceTestSuite struct {
	suite.Suite
	userMockCtrl *gomock.Controller
	userRepo     *user.MockRepository
	userService  *Service
}

func (suite *UserServiceTestSuite) SetupTest(t *testing.T) func() {
	suite.userMockCtrl = gomock.NewController(t)
	suite.userRepo = user.NewMockRepository(suite.userMockCtrl)
	var err error
	suite.userService, err = NewUserService(suite.userRepo)
	if err != nil {
		panic(err)
	}
	return func() {
		suite.userMockCtrl.Finish()
		suite.userMockCtrl = nil
		suite.userRepo = nil
		suite.userService = nil
	}
}

func (suite *UserServiceTestSuite) TestCreateUser() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()
	expectUser := &models.User{
		Model: gorm.Model{ID: uint(rand.Int()), CreatedAt: getRandTime(), UpdatedAt: getRandTime()},
	}
	suite.userRepo.EXPECT().CreateUser().Return(expectUser, nil)

	newUser, err := suite.userService.CreateUser()
	suite.NoError(err)
	suite.Equal(expectUser, newUser)
}

func (suite *UserServiceTestSuite) TestSetUserBalanceSuccessFully() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()

	userID := uint(rand.Uint64())
	balance := rand.Float64()
	suite.userRepo.EXPECT().CheckUserExists(userID).Return(true, nil)
	suite.userRepo.EXPECT().SetUserBalance(userID, balance).Return(balance, nil)

	returnedBalance, err := suite.userService.SetUserBalance(userID, balance)
	suite.NoError(err)
	suite.Equal(balance, returnedBalance)
}

func (suite *UserServiceTestSuite) TestSetUserBalanceIfUserNotFound() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()

	userID := uint(rand.Uint64())
	suite.userRepo.EXPECT().CheckUserExists(userID).Return(false, nil)
	suite.userRepo.EXPECT().SetUserBalance(userID, 0).Times(0)

	_, err := suite.userService.SetUserBalance(userID, 0)
	suite.ErrorContains(err, "user not found")
}

func getRandTime() time.Time {
	return time.Unix(rand.Int63(), 0)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
