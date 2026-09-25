package app

import (
	"FastCopyPast/fonts"
	"context"
	"fmt"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"golang.design/x/clipboard"
)

type ViewUI interface {
	Run(gtx layout.Context, th *material.Theme)
	Events(gtx layout.Context)
}

type App struct {
	Window     *app.Window
	Th         *material.Theme
	CancelRead context.CancelFunc
	MousePos   f32.Point
	Title      string
	Version    string
}

func NewApp(title, version string) *App {
	bufferInterceptor()

	return &App{
		Window:  new(app.Window),
		Th:      initTheme(),
		Title:   title,
		Version: version,
	}
}

func initTheme() *material.Theme {
	collection := fonts.ParseFont(fonts.SteticaFonts)
	th := material.NewTheme()
	th.Shaper = text.NewShaper(
		text.WithCollection(collection),
	)
	th.Face = fonts.SteticaMediumItalic
	return th
}

func bufferInterceptor() {
	err := clipboard.Init()
	if err != nil {
		panic(fmt.Sprintf("Не удалось инициализировать буфер обмена: %v", err))
	}
}

func (a *App) Run(v ViewUI) error {
	a.Window.Option(
		app.Title(a.Title),
		app.Decorated(false),
		app.Size(unit.Dp(316), unit.Dp(113)),
		app.MinSize(unit.Dp(316), unit.Dp(113)),
	)
	var ops op.Ops

	for {
		switch e := a.Window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			v.Events(gtx)

			v.Run(gtx, a.Th)

			e.Frame(gtx.Ops)
		}
	}
}
