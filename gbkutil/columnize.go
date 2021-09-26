package gbkutil

import (
	"bytes"
	"strings"
)

type ColumnAlign int

const (
	ColumnAlignLeft   ColumnAlign = 0
	ColumnAlignRight  ColumnAlign = 1
	ColumnAlignCenter ColumnAlign = 2
)

type Columnize struct {
	title       string
	columns     []string
	columnMasks []string
	result      string
	rows        [][]string
	// 表宽度
	tableWidth int
	// 列宽度
	columnWidths []int
	// 列对齐
	columnAligns []ColumnAlign
	// 列个数
	columnCount int
	isShowBorder bool
}

func NewColumnize() *Columnize {
	return &Columnize{
		title:        "",
		columns:      []string{},
		columnMasks:  []string{},
		result:       "",
		rows:         [][]string{},
		tableWidth:   0,
		columnWidths: []int{},
		columnCount:  0,
		isShowBorder: true,
	}
}

func NewColumnizeHideBorder() *Columnize {
	return &Columnize{
		title:        "",
		columns:      []string{},
		columnMasks:  []string{},
		result:       "",
		rows:         [][]string{},
		tableWidth:   0,
		columnWidths: []int{},
		columnCount:  0,
		isShowBorder: false,
	}
}

func (col *Columnize) String() string {
	return col.result
}

func (col *Columnize) SetTitle(title string) {
	col.title = strings.TrimSpace(title)
	col.calc()
}

func (col *Columnize) SetColumns(columns []string) {
	for i, column := range columns {
		columns[i] = strings.TrimSpace(column)
	}
	col.columns = columns
	col.calc()
}

func (col *Columnize) SetColumnMasks(columnMasks []string) {
	for i, columnMask := range columnMasks {
		columnMasks[i] = strings.TrimSpace(columnMask)
	}
	col.columnMasks = columnMasks
	col.calc()
}

func (col *Columnize) SetRows(rows [][]string) {
	col.rows = rows
	col.calc()
}

func (col *Columnize) AddRow(row []string) {
	col.rows = append(col.rows, row)
	col.calc()
}

// 设置是否显示边框
func (col *Columnize) SetShowBorder(isShowBorder bool) {
	col.isShowBorder = isShowBorder
}

func (col *Columnize) calc() {
	var buf bytes.Buffer
	col.columnWidths = []int{}
	col.columnAligns = []ColumnAlign{}
	col.columnCount = 0
	col.tableWidth = 0
	// 准备
	for _, row := range col.rows {
		for i, td := range row {
			tdLen := RuneLen(td)
			if len(col.columnWidths) <= i { // 不存在则新的就是最长
				col.columnWidths = append(col.columnWidths, tdLen)
				col.columnAligns = append(col.columnAligns, ColumnAlignLeft)
				col.columnCount++
			} else if col.columnWidths[i] < tdLen { // 遇到更长的了
				col.columnWidths[i] = tdLen
			}
		}
	}
	if len(col.columnMasks) > 0 { // 有列长度描述优先使用描述长度
		for i, columnMask := range col.columnMasks {
			maskLen := len(columnMask)
			if len(col.columnWidths) <= i { // 不存在则新的就是最长
				col.columnWidths = append(col.columnWidths, maskLen)
				col.columnAligns = append(col.columnAligns, ColumnAlignLeft)
				col.columnCount++
			} else if col.columnWidths[i] < maskLen {
				col.columnWidths[i] = maskLen
			}
			if columnMask[len(columnMask)-1] == ':' {// 右对齐
				if columnMask[0] == ':' {// 居中对齐
					col.columnAligns[i] = ColumnAlignCenter
				} else {// 右对齐
					col.columnAligns[i] = ColumnAlignRight
				}
			} else {// 默认左对齐
				col.columnAligns[i] = ColumnAlignLeft
			}
		}
	}
	if len(col.columns) > 0 {
		for i, column := range col.columns {
			columnWidth := RuneLen(column)
			if len(col.columnWidths) <= i { // 不存在则新的就是最长
				col.columnWidths = append(col.columnWidths, columnWidth)
				col.columnAligns = append(col.columnAligns, ColumnAlignLeft)
				col.columnCount++
			} else if col.columnWidths[i] < columnWidth {
				col.columnWidths[i] = columnWidth
			}

		}
	}
	// 生成
	for i := 0; i < col.columnCount; i++ {
		if col.isShowBorder {
			col.tableWidth += col.columnWidths[i] + 3
		} else {
			col.tableWidth += col.columnWidths[i] + 1
		}
	}
	col.tableWidth++

	if len(col.title) > 0 {
		titleLen := RuneLen(col.title)
		if col.tableWidth == 1 {// 还没有列的时候
			if col.isShowBorder {
				col.tableWidth = len(col.title) + 4
			} else {
				col.tableWidth = len(col.title) + 2
			}
		} else if col.tableWidth < titleLen {
			col.tableWidth = titleLen + 2
		}
		if col.isShowBorder {
			buf.WriteString("|")
		}
		buf.WriteString(PadCenterSpace(col.title, col.tableWidth-2))
		if col.isShowBorder {
			buf.WriteString("|")
		}
		buf.WriteString("\n")
		if col.isShowBorder {
			buf.WriteString("|")
		}
		buf.WriteString(strings.Repeat("-", col.tableWidth-2))
		if col.isShowBorder {
			buf.WriteString("|")
		}
		buf.WriteString("\n")
	}
	// Header
	for i := 0; i < col.columnCount; i++ {
		if i > 0 {
			buf.WriteString(" ")
		}
		if col.isShowBorder {
			buf.WriteString("| ")
		}
		if len(col.columns) > i {// 有列信息
			switch col.columnAligns[i] {
			case ColumnAlignLeft:
				buf.WriteString(PadRightSpace(col.columns[i], col.columnWidths[i]))
			case ColumnAlignRight:
				buf.WriteString(PadLeftSpace(col.columns[i], col.columnWidths[i]))
			case ColumnAlignCenter:
				buf.WriteString(PadCenterSpace(col.columns[i], col.columnWidths[i]))
			default:
				buf.WriteString(PadRightSpace(col.columns[i], col.columnWidths[i]))
			}
		} else {// 没有列信息
			buf.WriteString(PadLeftSpace("", col.columnWidths[i]))
		}
	}
	if col.isShowBorder {
		buf.WriteString("|")
	}
	buf.WriteString("\n")
	// Header Line
	for i := 0; i < col.columnCount; i++ {
		if i > 0 {
			buf.WriteString(" ")
		}
		if col.isShowBorder {
			buf.WriteString("| ")
		}
		if len(col.columnMasks) > i {
			buf.WriteString(PadCenterSpace(col.columnMasks[i], col.columnWidths[i]))
		} else {
			buf.WriteString(strings.Repeat("-", col.columnWidths[i]))
		}
	}
	if col.isShowBorder {
		buf.WriteString("|")
	}
	buf.WriteString("\n")
	// Body
	for _, row := range col.rows {
		for i, s := range row {
			if i > 0 {
				buf.WriteString(" ")
			}
			if col.isShowBorder {
				buf.WriteString("| ")
			}
			if len(col.columnAligns) > i {// 有对齐信息
				switch col.columnAligns[i] {
				case ColumnAlignLeft:
					buf.WriteString(PadRightSpace(s, col.columnWidths[i]))
				case ColumnAlignRight:
					buf.WriteString(PadLeftSpace(s, col.columnWidths[i]))
				case ColumnAlignCenter:
					buf.WriteString(PadCenterSpace(s, col.columnWidths[i]))
				default:
					buf.WriteString(PadRightSpace(s, col.columnWidths[i]))
				}
			} else {// 没有对齐信息，默认左对齐
				buf.WriteString(PadRightSpace(s, col.columnWidths[i]))
			}
		}
		if col.isShowBorder {
			buf.WriteString("|")
		}
		buf.WriteString("\n")
	}

	col.result = buf.String()
}
