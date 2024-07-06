package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

var (
	uidEntry, topicEntry, macEntry, broadcastAddressEntry, sshUserServerEntry *widget.Entry
	typeRadio                                                                 *widget.RadioGroup
	saveButton, saveAndRunButton                                              *widget.Button
)

func initWindowWidget() {
	uidEntry = widget.NewEntry()
	topicEntry = widget.NewEntry()
	topicEntry.SetPlaceHolder("xxx006")
	macEntry = widget.NewEntry()
	macEntry.SetPlaceHolder("00:00:00:00:00:00")
	broadcastAddressEntry = widget.NewEntry()
	broadcastAddressEntry.SetPlaceHolder("192.168.1.255")
	sshUserServerEntry = widget.NewEntry()
	sshUserServerEntry.SetPlaceHolder("root@192.168.1.1")
	typeRadio = widget.NewRadioGroup([]string{"TCP", "MQTT"}, func(s string) {})
	typeRadio.SetSelected("TCP")
	typeRadio.Horizontal = true
	saveButton = widget.NewButton("保存配置", func() {})
	saveButton.Resize(fyne.NewSize(50, 30))
	saveAndRunButton = widget.NewButton("保存并运行", func() {})
	saveAndRunButton.Resize(fyne.NewSize(50, 30))

	saveButton.OnTapped = func() {

	}
}

func LoadWindow() fyne.Window {
	gui := app.New()
	gui.Settings().SetTheme(&CNTheme{})

	w := gui.NewWindow("小爱同学控制电脑开关机")
	if desk, ok := gui.(desktop.App); ok {
		m := fyne.NewMenu("",
			fyne.NewMenuItem("显示窗口", func() {
				w.Show()
			}))
		desk.SetSystemTrayMenu(m)

		w.SetCloseIntercept(func() {
			w.Hide()
		})
	}

	initWindowWidget()
	form := widget.NewForm(
		widget.NewFormItem("通信类型", typeRadio),
		widget.NewFormItem("巴法uid", uidEntry),
		widget.NewFormItem("巴法topic", topicEntry),
		widget.NewFormItem("mac地址", macEntry),
		widget.NewFormItem("广播地址", broadcastAddressEntry),
		widget.NewFormItem("ssh用户@ip地址", sshUserServerEntry))
	w.SetContent(container.NewPadded(container.NewVBox(form, container.NewCenter(container.NewHBox(saveButton, saveAndRunButton)))))
	w.Resize(fyne.NewSize(800, 600))

	return w
}
