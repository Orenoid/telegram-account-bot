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

const helpMessage = `欢迎使用记账机器人！

以下是可用的命令：
	/start - 开始使用
	/day - 查看当日账单
	/month - 查看当月账单
	/set_keyboard - 设置快捷键盘
	/cancel - 取消当前操作

如何记账：
	直接向机器人发送要记录的账单类别
	待机器人回复后，再发送金额，即可完成本次记账

	金额前面带上\"+\"号表示收入，不带表示支出
	账单类别可以使用快捷键盘，也可以手动输入

如何自定义快捷键盘：
	请按照以下格式输入你想要设置的快捷键盘，例如：

	饮食,出行,杂项|娱乐,购物,房租|工资,基金

	其中\"｜\"表示换行
	在上面的例子中，则表示设置一个三行的快捷键盘，第一行设置了「饮食」、「出行」、「杂项」三个账单类别，以此类推
`

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
	defaultKeyboard := textToKeyboard("饮食,出行,杂项|娱乐,购物,房租|工资")
	err = ctx.Send(helpMessage, &telebot.ReplyMarkup{ReplyKeyboard: defaultKeyboard})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (hub *HandlersHub) HandleHelpCommand(ctx telebot.Context) error {
	err := ctx.Send(helpMessage)
	return errors.WithStack(err)
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

func (hub *HandlersHub) HandleSetKeyboardCommand(ctx telebot.Context) error {
	sender := ctx.Sender()
	if sender == nil {
		return nil
	}
	// 记录用户状态
	err := hub.userStateManager.SetUserState(sender.ID, &UserState{Type: SettingKeyboard})
	if err != nil {
		return err
	}
	// 发送提示信息
	err = ctx.Send("快捷键盘用于设置一些日常生活中的支出/收入类别，用于快速记录\n\n请按照以下格式输入你想要设置的快捷键盘，例如：\n\n饮食,出行,杂项|娱乐,购物|工资,基金\n\n其中\"｜\"表示换行，在上面的例子中，则表示设置一个三行的快捷键盘，第一行设置了「饮食」、「出行」、「杂项」三个账单类别，以此类推")
	return errors.WithStack(err)
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
	case SettingKeyboard:
		return hub.OnSettingKeyboard(ctx)
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

func (hub *HandlersHub) OnSettingKeyboard(ctx telebot.Context) error {
	keyboardStr := ctx.Text()
	keyboard := textToKeyboard(keyboardStr)
	err := hub.userStateManager.ClearUserState(ctx.Sender().ID)
	if err != nil {
		return err
	}
	err = ctx.Send("已设置", &telebot.ReplyMarkup{ReplyKeyboard: keyboard})
	return errors.WithStack(err)
}

func textToKeyboard(text string) [][]telebot.ReplyButton {
	var result [][]telebot.ReplyButton
	rows := strings.Split(text, "|")
	for _, rowStr := range rows {
		categories := strings.Split(rowStr, ",")
		btnsInRow := make([]telebot.ReplyButton, 0, len(categories))
		for _, category := range categories {
			btn := telebot.ReplyButton{Text: category}
			btnsInRow = append(btnsInRow, btn)
		}
		result = append(result, btnsInRow)
	}
	return result
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
	err = ctx.Edit("已撤销账单")
	return errors.WithStack(err)
}
