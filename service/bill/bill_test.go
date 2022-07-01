package bill

import (
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/orenoid/telegram-account-bot/dal/bill"
	"github.com/orenoid/telegram-account-bot/dal/user"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"math/rand"
	"testing"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	billRepo := bill.NewMockRepository(ctrl)
	userRepo := user.NewMockRepository(ctrl)
	service := NewService(billRepo, userRepo)
	assert.True(t, billRepo == service.billRepo)
	assert.True(t, userRepo == service.userRepo)
}

type BillServiceTestSuite struct {
	suite.Suite

	userMockCtrl *gomock.Controller
	userRepo     *user.MockRepository

	billMockCtrl *gomock.Controller
	billRepo     *bill.MockRepository
	billService  *Service
}

func (suite *BillServiceTestSuite) SetupTest(t *testing.T) func() {
	suite.billMockCtrl = gomock.NewController(t)
	suite.billRepo = bill.NewMockRepository(suite.billMockCtrl)
	suite.userMockCtrl = gomock.NewController(t)
	suite.userRepo = user.NewMockRepository(suite.userMockCtrl)

	suite.billService = NewService(suite.billRepo, suite.userRepo)

	return func() {
		suite.billMockCtrl.Finish()
		suite.billMockCtrl = nil
		suite.billRepo = nil
		suite.billMockCtrl = nil
		suite.userRepo = nil
		suite.billService = nil
	}
}

func (suite *BillServiceTestSuite) TestCreateBill() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()

	params := []struct {
		userID     uint
		billAmount float64
		category   string
		billName   string
	}{
		{userID: 42, billAmount: 22.33, category: "饮食", billName: ""},
		{userID: 43, billAmount: 300, category: "娱乐", billName: "十三机兵防卫圈"},
	}
	for i, param := range params {
		suite.Run(fmt.Sprintf("param%d", i), func() {
			var optsI []interface{}
			var opts []bill.CreateBillOptions
			newBill := &models.Bill{
				UserID: param.userID, Amount: decimal.NewFromFloat(param.billAmount), Category: param.category, Model: gorm.Model{ID: uint(rand.Int())},
			}
			if param.billName != "" {
				newBill.Name = sql.NullString{String: param.billName, Valid: true}
				optsI = append(optsI, bill.CreateBillOptions{Name: &param.billName})
				opts = append(opts, bill.CreateBillOptions{Name: &param.billName})
			}
			suite.userRepo.EXPECT().CheckUserExists(param.userID).Return(true, nil)
			suite.billRepo.EXPECT().CreateBillAndUpdateUserBalance(param.userID, param.billAmount, param.category, optsI...).Return(newBill, nil)
			returnedBill, err := suite.billService.CreateNewBill(param.userID, param.billAmount, param.category, opts...)
			suite.NoError(err)
			suite.Equal(newBill, returnedBill)
		})

	}
}

func (suite *BillServiceTestSuite) TestCreateBillIfUserNotFound() {
	tearDown := suite.SetupTest(suite.T())
	defer tearDown()

	userID := uint(42)
	suite.userRepo.EXPECT().CheckUserExists(userID).Return(false, nil)
	suite.billRepo.EXPECT().CreateBillAndUpdateUserBalance(userID, 0, "").Times(0)
	returnedBill, err := suite.billService.CreateNewBill(userID, 0, "")
	suite.ErrorContains(err, "user not exists")
	suite.Nil(returnedBill)
}

func TestBillServiceSuite(t *testing.T) {
	suite.Run(t, new(BillServiceTestSuite))
}
