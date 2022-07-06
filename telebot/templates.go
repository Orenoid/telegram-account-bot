package telebot

import (
	"fmt"
	"github.com/orenoid/telegram-account-bot/models"
	"github.com/shopspring/decimal"
	"sort"
	"strings"
)

type Template interface {
	Render() string
}

type BillCreatedTemplate struct {
	bill *models.Bill
}

func (template *BillCreatedTemplate) Render() string {
	format :=
		`%s：%s 元，类别：%s
点击 "撤销" 可撤回该账单并回滚余额`
	var inOrOut, amountStr string
	if template.bill.Amount.LessThan(decimal.NewFromFloat(0)) {
		inOrOut, amountStr = "支出", template.bill.Amount.Abs().String()
	} else {
		inOrOut, amountStr = "收入", template.bill.Amount.String()
	}
	return fmt.Sprintf(format, inOrOut, amountStr, template.bill.Category)
}

type MonthTitleTemplate struct {
	Year  int
	Month int
}

func (template *MonthTitleTemplate) Render() string {
	return fmt.Sprintf("%d年%d月收支统计", template.Year, template.Month)
}

type DateTitleTemplate struct {
	Year     int
	Month    int
	Day      int
	ShowYear bool
}

func (template *DateTitleTemplate) Render() string {
	if template.ShowYear {
		return fmt.Sprintf("%d年%d月%d日收支统计", template.Year, template.Month, template.Day)
	} else {
		return fmt.Sprintf("%d月%d日收支统计", template.Month, template.Day)

	}
}

type BillListTemplate struct {
	Bills         []*models.Bill
	MergeCategory bool // 对相同类别的账单求总和后展示
}

func (template *BillListTemplate) Render() string {
	if len(template.Bills) == 0 {
		return "暂无收支记录"
	}

	result := &strings.Builder{}
	expendSection := template.expendSection()
	if len(expendSection) > 0 {
		result.Write([]byte(expendSection))
	}
	incomeSection := template.incomeSection()
	if len(expendSection) > 0 && len(incomeSection) > 0 {
		// 支出与收入记录之间多间隔一行
		result.Write([]byte("\n\n"))
	}
	if len(incomeSection) > 0 {
		result.Write([]byte(incomeSection))
	}

	return result.String()
}

// expendSection 支出记录
func (template *BillListTemplate) expendSection() string {
	expendSum := decimal.NewFromFloat(0)
	categoryMapping := map[string]decimal.Decimal{}
	var expendBills []*models.Bill

	// 筛选支出账单
	for _, bill := range template.Bills {
		if bill.Amount.LessThan(decimal.NewFromFloat(0)) {
			expendBills = append(expendBills, bill)
			expendSum = expendSum.Add(bill.Amount)
			categoryMapping[bill.Category] = categoryMapping[bill.Category].Add(bill.Amount)
		}
	}
	if len(expendBills) == 0 {
		return ""
	}

	result := &strings.Builder{}
	result.Write([]byte(fmt.Sprintf("合计支出：%s 元\n\n", expendSum.Abs().String())))
	if !template.MergeCategory {
		// 展示所有账单时，按创建时间排序
		sort.Sort(billsSortByCreateTime(expendBills))
		count := 0
		for _, bill := range expendBills {
			count++
			result.Write([]byte(fmt.Sprintf("       - %s：%s 元", bill.Category, bill.Amount.Abs().String())))
			if count < len(expendBills) {
				result.Write([]byte("\n"))
			}
		}
	} else {
		count := 0
		// 按类别支出总和由大到小排序
		for _, category := range template.getSortedCategories(categoryMapping) {
			amount := categoryMapping[category]
			count++
			result.Write([]byte(fmt.Sprintf("       - %s：%s 元", category, amount.Abs().String())))
			if count < len(categoryMapping) {
				result.Write([]byte("\n"))
			}
		}
	}
	return result.String()
}

// incomeSection 收入记录
func (template *BillListTemplate) incomeSection() string {
	var incomeBills []*models.Bill
	incomeSum := decimal.NewFromFloat(0)
	categoryMapping := map[string]decimal.Decimal{}

	// 筛选收入账单
	for _, bill := range template.Bills {
		if !bill.Amount.LessThanOrEqual(decimal.NewFromFloat(0)) {
			incomeBills = append(incomeBills, bill)
			incomeSum = incomeSum.Add(bill.Amount)
			categoryMapping[bill.Category] = categoryMapping[bill.Category].Add(bill.Amount)
		}
	}
	if len(incomeBills) == 0 {
		return ""
	}

	result := &strings.Builder{}
	result.Write([]byte(fmt.Sprintf("合计收入：%s 元\n\n", incomeSum.String())))
	if !template.MergeCategory {
		// 展示所有账单时，按创建时间排序
		sort.Sort(billsSortByCreateTime(incomeBills))
		count := 0
		for _, bill := range incomeBills {
			count++
			result.Write([]byte(fmt.Sprintf("       - %s：%s 元", bill.Category, bill.Amount.String())))
			if count < len(incomeBills) {
				result.Write([]byte("\n"))
			}
		}
	} else {
		count := 0
		categories := template.getSortedCategories(categoryMapping)
		// 按类别收入总和由大到小排序
		for i := len(categories) - 1; i >= 0; i-- {
			category := categories[i]
			amount := categoryMapping[category]
			count++
			result.Write([]byte(fmt.Sprintf("       - %s：%s 元", category, amount.Abs().String())))
			if count < len(categoryMapping) {
				result.Write([]byte("\n"))
			}
		}
	}
	return result.String()
}

func (template *BillListTemplate) getSortedCategories(categoriesMapping map[string]decimal.Decimal) []string {
	var sortHelper [][2]interface{}
	for cate, amount := range categoriesMapping {
		sortHelper = append(sortHelper, [2]interface{}{cate, amount})
	}
	sort.Sort(sortByAmount(sortHelper))
	var categories []string
	for _, item := range sortHelper {
		categories = append(categories, item[0].(string))
	}
	return categories
}

type sortByAmount [][2]interface{}

func (items sortByAmount) Len() int {
	return len(items)
}

func (items sortByAmount) Less(i, j int) bool {
	return items[i][1].(decimal.Decimal).LessThan(items[j][1].(decimal.Decimal))
}

func (items sortByAmount) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

type billsSortByCreateTime []*models.Bill

func (bills billsSortByCreateTime) Len() int {
	return len(bills)
}

func (bills billsSortByCreateTime) Less(i, j int) bool {
	return bills[i].CreatedAt.Before(bills[j].CreatedAt)
}

func (bills billsSortByCreateTime) Swap(i, j int) {
	bills[i], bills[j] = bills[j], bills[i]
}
