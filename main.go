package main

import (
	"github.com/tg567/bemfa-wol/gui"
	"github.com/tg567/bemfa-wol/utils"
)

func main() {
	if err := utils.LoadConfig(); err != nil {
		return
	}
	w := gui.LoadWindow()
	w.ShowAndRun()
}
