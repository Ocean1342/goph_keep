package printer

import "github.com/spf13/cobra"

// Printer - отвечает за печать информации в консоль пользователюддддд
type Printer struct {
	cmd *cobra.Command
}

func New(cmd *cobra.Command) *Printer {
	return &Printer{
		cmd: cmd,
	}
}

func (p Printer) Println(data string) {
	p.cmd.Println(data)
}

func (p Printer) PrintErr(data string) {
	p.cmd.PrintErrln(data)
}
