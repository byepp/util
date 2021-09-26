package decimalutil

import (
	"bytes"
	"github.com/shopspring/decimal"
	"strings"
)

var (
	Hundred = decimal.NewFromInt(100)
)

// 计算增长率 返回内容带百分号
func GrowthRate(presentAmount, historyAmount decimal.Decimal) string {
	if historyAmount.IsZero() || historyAmount.IsNegative() || presentAmount.IsNegative() {// 参数不正确
		if presentAmount.IsZero() {
			return "+0.00%"
		}
		return "NaN%"
	}
	var buf bytes.Buffer
	if presentAmount.GreaterThan(historyAmount) { // 如果是增长需要增加+号
		buf.WriteString("+")
	}
	buf.WriteString(presentAmount.Sub(historyAmount).Div(historyAmount).Mul(Hundred).StringFixed(2))
	buf.WriteString("%")
	return buf.String()
}


// 计算增长率 返回内容带百分号
func GrowthRateFixed(presentAmount, historyAmount decimal.Decimal, places int) string {
	if historyAmount.IsZero() || historyAmount.IsNegative() || presentAmount.IsNegative() {// 参数不正确
		if presentAmount.IsZero() {
			if places > 0 {
				return "+0." + strings.Repeat("0", places) + "%"
			}
			return "+0%"
		}
		return "NaN%"
	}
	var buf bytes.Buffer
	if presentAmount.GreaterThan(historyAmount) { // 如果是增长需要增加+号
		buf.WriteString("+")
	}
	buf.WriteString(presentAmount.Sub(historyAmount).Div(historyAmount).Mul(Hundred).StringFixed(int32(places)))
	buf.WriteString("%")
	return buf.String()
}

// 计算使用率 返回内容带百分号
func UsageRate(val, total decimal.Decimal) string {
	if total.IsZero() || val.IsNegative() || total.IsNegative() {// 参数不正确
		if val.IsZero() {
			return "0.00%"
		}
		return "NaN%"
	}
	var buf bytes.Buffer
	buf.WriteString(val.Div(total).Mul(Hundred).StringFixed(2))
	buf.WriteString("%")
	return buf.String()
}
