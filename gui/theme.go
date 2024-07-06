package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CNTheme struct{}

var _ fyne.Theme = (*CNTheme)(nil)

// resourceSIMHEITTF 对应的是 bundle.go 中的变量名
func (m CNTheme) Font(fyne.TextStyle) fyne.Resource {
	return resourceSIMHEITTF
}

func (*CNTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*CNTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*CNTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
