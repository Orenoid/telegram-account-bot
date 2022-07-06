package telebot

import (
	"fmt"
	billDAL "github.com/orenoid/telegram-account-bot/dal/bill"
	"github.com/orenoid/telegram-account-bot/service/bill"
	"github.com/orenoid/telegram-account-bot/service/telegram"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type HandlersHub struct {
	teleService      *telegram.Service
	billService      *bill.Service
	userStateManager UserStateManager
}

func NewHandlerHub(billService *bill.Service, teleService *telegram.Service, userStateManager UserStateManager) *HandlersHub {
	return &HandlersHub{billService: billService, teleService: teleService, userStateManager: userStateManager}
}

func (hub *HandlersHub) HandleStartCommand(ctx telebot.Context) error {
	chat := ctx.Chat()
	if chat == nil {
		return errors.New("nil chat of context")
	}
	sender := ctx.Sender()
	_, err := hub.teleService.CreateOrUpdateTelegramUser(sender.ID, sender.Username, chat.ID)
	if err != nil {
		return err
	}
	// TODO set default keyboard
	err = ctx.Send("hello") // TODO send help message
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (hub *HandlersHub) HandleDayCommand(ctx telebot.Context) error {
	sender := ctx.Sender()
	now := time.Now()
	begin, end := getDayRange(now)

	baseUserID, err := hub.teleService.GetBaseUserID(sender.ID)
	if err != nil {
		return err
	}
	bills, err := hub.billService.GetUserBillsByCreateTime(baseUserID,
		billDAL.GetUserBillsByCreateTimeOptions{GreaterThan: begin, GreaterOrEqual: true, LessThan: end})
	if err != nil {
		return err
	}
	return ctx.Send(&DateBillsSender{bills, now.Year(), int(now.Month()), now.Day(), false})
}

// 获取某个时刻当天的0点-24点范围
func getDayRange(t time.Time) (time.Time, time.Time) {
	begin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tomorrow := t.Add(24 * time.Hour)
	end := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	return begin, end
}

func (hub *HandlersHub) HandleText(ctx telebot.Context) error {
	userState, err := hub.userStateManager.GetUserState(ctx.Sender().ID)
	if err != nil {
		return err
	}
	switch userState.Type {
	case Empty:
		return hub.OnEmpty(ctx)
	case CreatingBill:
		return hub.OnCreatingBill(ctx, userState)
	}
	return nil
}

func (hub *HandlersHub) OnEmpty(ctx telebot.Context) error {
	text := ctx.Text()
	if len(text) == 0 {
		return nil
	}
	category, name := ParseBill(text)
	err := hub.userStateManager.SetUserState(ctx.Sender().ID,
		&UserState{Type: CreatingBill, BillCategory: &category, BillName: name},
	)
	if err != nil {
		return err
	}
	err = ctx.Send(fmt.Sprintf("账单类别：%s，请输入账单金额", category))
	return errors.WithStack(err)
}

func (hub *HandlersHub) OnCreatingBill(ctx telebot.Context, userState *UserState) error {
	amount, err := parseAmount(ctx.Text())
	if err != nil {
		return err
	}
	sender := ctx.Sender()
	baseUserID, err := hub.teleService.GetBaseUserID(sender.ID)
	if err != nil {
		return err
	}
	newBill, err := hub.billService.CreateNewBill(
		baseUserID, amount, *userState.BillCategory, billDAL.CreateBillOptions{Name: userState.BillName},
	)
	if err != nil {
		return err
	}
	err = hub.userStateManager.ClearUserState(sender.ID)
	if err != nil {
		return err
	}
	err = ctx.Send(&NewBillSender{newBill})
	return errors.WithStack(err)
}

var validAmount = regexp.MustCompile("^([+-]?)([0-9]*\\.?[0-9]+)$")

//parseAmount 解析数额，若前面不带 "+"，则默认会解析为负数（平时大多数时候为支出）
func parseAmount(text string) (float64, error) {
	matchResult := validAmount.FindStringSubmatch(text)
	if len(matchResult) == 0 {
		return 0, errors.Errorf("invalid amount text: %s", text)
	} else if len(matchResult) == 3 {
		amount, err := strconv.ParseFloat(matchResult[2], 64)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		if matchResult[1] != "+" {
			amount = -amount
		}
		return amount, nil
	}
	return 0, errors.Errorf("invalid amount text: %s", text)
}

func ParseBill(text string) (string, *string) {
	ss := strings.SplitN(text, " ", 2)
	var category, name string
	if len(ss) == 1 {
		category = ss[0]
		return category, nil
	} else {
		category, name = ss[0], ss[1]
		return category, &name
	}
}
