package telebot

import (
	"encoding/json"
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

func (hub *HandlersHub) HandleMonthCommand(ctx telebot.Context) error {
	sender := ctx.Sender()
	now := time.Now()
	begin, end := getMonthRange(now)

	baseUserID, err := hub.teleService.GetBaseUserID(sender.ID)
	if err != nil {
		return err
	}
	bills, err := hub.billService.GetUserBillsByCreateTime(baseUserID,
		billDAL.GetUserBillsByCreateTimeOptions{GreaterThan: begin, GreaterOrEqual: true, LessThan: end})
	if err != nil {
		return err
	}
	var sendable telebot.Sendable = &MonthBillsSender{Bills: bills, Year: now.Year(), Month: int(now.Month())}
	return ctx.Send(sendable)
}

func (hub *HandlersHub) HandleCancelCommand(ctx telebot.Context) error {
	sender := ctx.Sender()
	if sender == nil {
		return nil
	}
	err := hub.userStateManager.ClearUserState(sender.ID)
	if err != nil {
		return err
	}
	err = ctx.Send("已取消")
	return err
}

func getMonthRange(t time.Time) (time.Time, time.Time) {
	currentYear, currentMonth, _ := t.Date()
	currentLocation := t.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return firstOfMonth, lastOfMonth
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
	err = ctx.Send(fmt.Sprintf("账单类别：%s，请输入账单金额\n默认记做支出，若想记为收入，可在金额前带上\"+\"号\n若想取消本次操作，请输入 /cancel", category))
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

// HandleDayBillSelectionCallback 处理切换日期账单的回调事件
func (hub *HandlersHub) HandleDayBillSelectionCallback(ctx telebot.Context) error {
	// 解析回调按钮的日期
	callback := ctx.Callback()
	if callback == nil {
		return nil
	}
	data := DayBillBtnData{}
	err := json.Unmarshal([]byte(callback.Data), &data)
	if err != nil {
		return errors.WithStack(err)
	}
	// 查询当日账单
	baseUserID, err := hub.teleService.GetBaseUserID(callback.Sender.ID)
	if err != nil {
		return err
	}
	date := time.Date(data.Year, time.Month(data.Month), data.Day, 0, 0, 0, 0, time.Local)
	begin, end := getDayRange(date)
	bills, err := hub.billService.GetUserBillsByCreateTime(baseUserID,
		billDAL.GetUserBillsByCreateTimeOptions{GreaterThan: begin, GreaterOrEqual: true, LessThan: end})
	if err != nil {
		return err
	}
	// 更新消息，切换账单
	sender := &DateBillsSender{bills, data.Year, data.Month, data.Day, false}
	err = ctx.Edit(sender.Text(), sender.ReplyMarkup())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// HandleMonthBillSelectionCallback 处理切换月度账单的回调事件
func (hub *HandlersHub) HandleMonthBillSelectionCallback(ctx telebot.Context) error {
	callback := ctx.Callback()
	if callback == nil {
		return nil
	}
	data := MonthBillBtnData{}
	err := json.Unmarshal([]byte(callback.Data), &data)
	if err != nil {
		return errors.WithStack(err)
	}
	// 查询月度账单
	baseUserID, err := hub.teleService.GetBaseUserID(callback.Sender.ID)
	if err != nil {
		return err
	}
	date := time.Date(data.Year, time.Month(data.Month), 1, 0, 0, 0, 0, time.Local)
	begin, end := getMonthRange(date)
	bills, err := hub.billService.GetUserBillsByCreateTime(baseUserID,
		billDAL.GetUserBillsByCreateTimeOptions{GreaterThan: begin, GreaterOrEqual: true, LessThan: end})
	if err != nil {
		return err
	}
	// 更新消息，切换账单
	sender := &MonthBillsSender{bills, data.Year, data.Month}
	err = ctx.Edit(sender.Text(), sender.ReplyMarkup())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (hub *HandlersHub) HandleCancelBillCallback(ctx telebot.Context) error {
	callback := ctx.Callback()
	if callback == nil {
		return nil
	}
	data := CancelBillData{}
	err := json.Unmarshal([]byte(callback.Data), &data)
	if err != nil {
		return errors.WithStack(err)
	}
	err = hub.billService.CancelBillAndUpdateUserBalance(data.BillID)
	if err != nil {
		return err
	}
	err = ctx.Edit("已撤销账单并回滚余额")
	return errors.WithStack(err)
}
