package pkg

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

type Table struct {
	header      []string
	footer      []string
	isSetBorder bool
	w           io.Writer
}

func NewTable(w io.Writer) *Table {
	return &Table{
		isSetBorder: true,
		w:           w,
	}
}

func (d *Table) SetHeader(header []string) {
	d.header = header
}

func (d *Table) SetBorder(isSetBorder bool) {
	d.isSetBorder = isSetBorder
}

func (d *Table) Write(data [][]string) {
	table := tablewriter.NewWriter(d.w)

	// header
	table.SetHeader(d.header)

	// footer
	if d.footer != nil {
		table.SetFooter(d.footer)
	}

	// border
	table.SetBorder(d.isSetBorder)

	// data
	table.AppendBulk(data)

	table.Render()
}
