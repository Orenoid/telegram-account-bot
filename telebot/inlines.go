package telebot

import (
	"encoding/json"
	"gopkg.in/telebot.v3"
	"time"
)

var (
	prevDayBillBtnTmpl = telebot.Btn{Text: "<️", Unique: "dayBillSelectorBtn"}
	nextDayBillBtnTmpl = telebot.Btn{Text: ">️", Unique: "dayBillSelectorBtn"}
	cancelBillBtnTmpl  = telebot.Btn{Text: "撤销", Unique: "cancelBillBtn"}
)

type DayBillBtnData struct {
	Year, Month, Day int
}

// PrevDayBillBtn 根据传入的年月日，生成用于将账单切换到前一天的按钮
func PrevDayBillBtn(year, month, day int) telebot.Btn {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	prevDate := date.Add(-24 * time.Hour)
	data := &DayBillBtnData{prevDate.Year(), int(prevDate.Month()), prevDate.Day()}
	dataRaw, _ := json.Marshal(data)
	return telebot.Btn{
		Text:   prevDayBillBtnTmpl.Text,
		Unique: prevDayBillBtnTmpl.Unique,
		Data:   string(dataRaw),
	}
}

// NextDayBillBtn 根据传入的年月日，生成用于将账单切换到后一天的按钮
func NextDayBillBtn(year, month, day int) telebot.Btn {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	nextDate := date.Add(24 * time.Hour)
	data := &DayBillBtnData{nextDate.Year(), int(nextDate.Month()), nextDate.Day()}
	dataRaw, _ := json.Marshal(data)
	return telebot.Btn{
		Text:   nextDayBillBtnTmpl.Text,
		Unique: nextDayBillBtnTmpl.Unique,
		Data:   string(dataRaw),
	}
}

type CancelBillData struct {
	BillID uint
}

// CancelBillBtn 构造用于取消某个订单的按钮
func CancelBillBtn(billID uint) telebot.Btn {
	data := &CancelBillData{billID}
	dataRaw, _ := json.Marshal(data)
	return telebot.Btn{
		Text:   cancelBillBtnTmpl.Text,
		Unique: cancelBillBtnTmpl.Unique,
		Data:   string(dataRaw),
	}
}
