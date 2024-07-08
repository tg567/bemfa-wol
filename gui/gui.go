package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/tg567/bemfa-wol/connection"
	"github.com/tg567/bemfa-wol/utils"
)

var (
	uidEntry, topicEntry, macEntry, broadcastAddressEntry, sshUserServerEntry *widget.Entry
	typeRadio                                                                 *widget.RadioGroup
	saveButton, saveAndRunButton                                              *widget.Button
	running                                                                   bool
	closeChan                                                                 chan struct{}
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
	typeRadio = widget.NewRadioGroup([]string{utils.WOL_TYPE_MQTT, utils.WOL_TYPE_TCP}, func(s string) {})
	typeRadio.SetSelected(utils.WOL_TYPE_TCP)
	typeRadio.Horizontal = true
	saveButton = widget.NewButton("保存配置", func() {})
	saveAndRunButton = widget.NewButton("保存并运行", func() {})

	uidEntry.SetText(utils.WolConfig.UID)
	topicEntry.SetText(utils.WolConfig.Topic)
	macEntry.SetText(utils.WolConfig.Mac)
	broadcastAddressEntry.SetText(utils.WolConfig.Broadcast)
	sshUserServerEntry.SetText(utils.WolConfig.SSH)
	if utils.WolConfig.Type != utils.WOL_TYPE_MQTT && utils.WolConfig.Type != utils.WOL_TYPE_TCP {
		typeRadio.SetSelected(utils.WOL_TYPE_TCP)
	} else {
		typeRadio.SetSelected(utils.WolConfig.Type)
	}

	saveButton.OnTapped = func() { saveConfig() }
	saveAndRunButton.OnTapped = saveAndRunFunc
}

func saveConfig() {
	utils.WolConfig.UID = uidEntry.Text
	utils.WolConfig.Topic = topicEntry.Text
	utils.WolConfig.Mac = macEntry.Text
	utils.WolConfig.Broadcast = broadcastAddressEntry.Text
	utils.WolConfig.SSH = sshUserServerEntry.Text
	utils.WolConfig.Type = typeRadio.Selected
	if err := utils.SaveConfig(); err != nil {
		utils.Println("保存配置文件错误", err)
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
	w.SetContent(container.NewPadded(container.NewVBox(form, container.NewGridWithColumns(2, saveButton, saveAndRunButton))))
	w.Resize(fyne.NewSize(400, 0))

	return w
}

var saveAndRunFunc = func() {
	if running {
		running = false
		close(closeChan)
	} else {
		closeChan = make(chan struct{})
		saveConfig()
		if utils.WolConfig.UID == "" || utils.WolConfig.Topic == "" || utils.WolConfig.Mac == "" || utils.WolConfig.Broadcast == "" {
			utils.Println("唤醒参数不能为空")
			return
		}
		switch utils.WolConfig.Type {
		case utils.WOL_TYPE_MQTT:
			connection.MqttWOL(closeChan)
		case utils.WOL_TYPE_TCP:
			go connection.TcpWOL(closeChan)
		default:
			utils.Println("通信类型错误")
			return
		}
		running = true
		saveAndRunButton.SetText("停止")
	}
}
