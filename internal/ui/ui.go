package ui

import (
	"FastCopyPast/internal/app"
	"image"
	"image/color"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Opacity interface {
	UpdateOpacity(alpha uint8)
}

type (
	C = layout.Context
	D = layout.Dimensions
)

var (
	black    = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	red      = color.NRGBA{R: 255, A: 255}
	green    = color.NRGBA{G: 200, A: 255}
	blue     = color.NRGBA{G: 100, B: 255, A: 255}
	purple   = color.NRGBA{R: 255, G: 0, B: 255, A: 255}
	elements = color.NRGBA{R: 10, G: 110, B: 80, A: 255}
)

type UI struct {
	a               *app.App
	activeSwc       widget.Bool
	aboveAll        widget.Bool
	newLinePast     widget.Bool
	BufferField     *widget.Editor
	SeparatorField  *widget.Editor
	KegelField      *widget.Editor
	OpacityValue    *widget.Float
	closeIcon       widget.Icon
	settingIcon     widget.Icon
	hoverIcon       widget.Icon
	brightnessIcon  widget.Icon
	powerIcon       widget.Icon
	closeBtn        *widget.Clickable
	settingBtn      *widget.Clickable
	hoverBtn        *widget.Clickable
	black           *widget.Clickable
	red             *widget.Clickable
	green           *widget.Clickable
	blue            *widget.Clickable
	purple          *widget.Clickable
	SettingsOn      bool
	prevWindowSize  image.Point
	divider         string
	bufferTextSize  int
	bufferTextColor *color.NRGBA
}

func NewUI(a *app.App) *UI {
	cl, _ := widget.NewIcon(icons.NavigationClose)
	set, _ := widget.NewIcon(icons.ActionSettings)
	hov, _ := widget.NewIcon(icons.NavigationExpandMore)
	br, _ := widget.NewIcon(icons.ActionSettingsBrightness)
	pwr, _ := widget.NewIcon(icons.ActionSettingsPower)

	return &UI{
		a:           a,
		BufferField: &widget.Editor{},
		SeparatorField: &widget.Editor{
			Alignment: text.Middle,
		},
		KegelField: &widget.Editor{
			Alignment: text.Middle,
			Filter:    "0,1,2,3,4,5,6,7,8,9",
			MaxLen:    2,
		},
		OpacityValue: &widget.Float{
			Value: 1,
		},
		closeIcon:       *cl,
		settingIcon:     *set,
		hoverIcon:       *hov,
		brightnessIcon:  *br,
		powerIcon:       *pwr,
		closeBtn:        new(widget.Clickable),
		settingBtn:      new(widget.Clickable),
		hoverBtn:        new(widget.Clickable),
		black:           new(widget.Clickable),
		red:             new(widget.Clickable),
		green:           new(widget.Clickable),
		blue:            new(widget.Clickable),
		purple:          new(widget.Clickable),
		bufferTextSize:  16,
		bufferTextColor: &black,
	}
}

func (ui *UI) Run(gtx C, th *material.Theme) {
	r := image.Rectangle{Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y}}
	area := clip.Rect(r).Push(gtx.Ops)
	system.ActionInputOp(system.ActionMove).Add(gtx.Ops)
	area.Pop()

	if ui.OpacityValue.Update(gtx) {
		alpha := ui.OpacityValue.Value * 255

		if alpha < 55 {
			alpha = 55
		}

		app.UpdateOpacity(uint8(alpha))
	}

	layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					//buffer field
					layout.Flexed(1, func(gtx C) D {
						return layout.Inset{Right: unit.Dp(5), Bottom: unit.Dp(5)}.Layout(gtx, func(gtx C) D {
							border := widget.Border{
								Color:        color.NRGBA{G: 0, A: 60},
								Width:        unit.Dp(3),
								CornerRadius: unit.Dp(4),
							}

							return border.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
									minW := gtx.Dp(unit.Dp(260))
									minH := gtx.Dp(unit.Dp(65))

									if gtx.Constraints.Min.X < minW {
										gtx.Constraints.Min.X = minW
									}
									if gtx.Constraints.Min.Y < minH {
										gtx.Constraints.Min.Y = minH
									}

									f := material.Editor(th, ui.BufferField, "")
									f.TextSize = unit.Sp(ui.bufferTextSize)
									f.Color = *ui.bufferTextColor
									f.Font.Typeface = "SteticaMedium"

									return f.Layout(gtx)
								})
							})
						})
					}),
					//close, setting buttons
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(btnParams(ui.closeBtn, ui.closeIcon, th, "close").Layout),
							layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
							layout.Rigid(btnParams(ui.settingBtn, ui.settingIcon, th, "setting").Layout),
							layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
							layout.Rigid(btnParams(ui.hoverBtn, ui.hoverIcon, th, "hover").Layout),
						)
					}),
				)
			}),
			//settings
			layout.Rigid(func(gtx C) D {
				if !ui.SettingsOn {
					return D{}
				}
				//return layout.Flex{Axis: layout.Vertical}.Layout(gtx)
				return layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceEnd,
				}.Layout(gtx,
					//on/off; separator
					layout.Rigid(settingRow(th, &ui.activeSwc, "active on/off", "вкл/выкл", func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(material.Label(th, unit.Sp(16), "Разделитель: ").Layout),
							layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
							layout.Rigid(func(gtx C) D { return fieldsParam(gtx, th, ui.SeparatorField) }),
						)
					})),
					// Отступ между строками настроек
					layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
					//above all
					layout.Rigid(settingRow(th, &ui.aboveAll, "above all", "поверх всех окон", nil)),
					// Отступ между строками настроек
					layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
					//column sort
					layout.Rigid(settingRow(th, &ui.newLinePast, "new line", "вставка с новой строки", nil)),
					//kegel
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(material.Label(th, unit.Sp(16), "Размер шрифта: ").Layout),
							layout.Rigid(func(gtx C) D {
								kegel := ui.KegelField.Text()
								if kegel == "" || kegel == "0" {
									ui.KegelField.SetText("16")
								}
								return fieldsParam(gtx, th, ui.KegelField)
							}),
							//text color
							layout.Rigid(colorButtons(ui.black, black, th)),
							layout.Rigid(colorButtons(ui.red, red, th)),
							layout.Rigid(colorButtons(ui.green, green, th)),
							layout.Rigid(colorButtons(ui.blue, blue, th)),
							layout.Rigid(colorButtons(ui.purple, purple, th)),
						)
					}),
					//opacity
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(material.Label(th, unit.Sp(16), "Прозрачность: ").Layout),
							layout.Rigid(func(gtx C) D {
								gtx.Constraints.Min.X = gtx.Dp(100)
								s := material.Slider(th, ui.OpacityValue)
								s.Color = elements
								return s.Layout(gtx)
							}),
						)
					}),
					//version
					layout.Rigid(about(th, ui.a.Version)),
				)
			}),
		)
	})
}

func settingRow(th *material.Theme, swc *widget.Bool, hint string, labelText string, extra layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Переключатель
			layout.Rigid(swcParams(swc, th, hint).Layout),
			// Отступ между переключателем и текстом вместо пробелов
			layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
			// Текст опции
			layout.Rigid(material.Label(th, unit.Sp(16), labelText).Layout),
			// Опциональный кастомный элемент (например, инпут "Разделитель")
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if extra == nil {
					return layout.Dimensions{}
				}
				// Добавляем небольшой отступ перед кастомным полем
				return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, extra)
			}),
		)
	}
}

func swcParams(sw *widget.Bool, th *material.Theme, desc string) material.SwitchStyle {
	s := material.Switch(th, sw, desc)
	s.Color.Enabled = elements

	return s
}

func btnParams(btn *widget.Clickable, btnIcon widget.Icon, th *material.Theme, desc string) material.IconButtonStyle {
	b := material.IconButton(th, btn, &btnIcon, desc)
	b.Size = unit.Dp(25)
	b.Inset = layout.UniformInset(unit.Dp(2))
	b.Background = elements

	return b
}

func fieldsParam(gtx layout.Context, th *material.Theme, w *widget.Editor) D {
	border := widget.Border{
		Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 60},
		CornerRadius: unit.Dp(4),
		Width:        unit.Dp(2),
	}

	return border.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Dp(25)
		return layout.UniformInset(4).Layout(gtx, material.Editor(th, w, "").Layout)
	})
}

func colorButtons(btn *widget.Clickable, col color.NRGBA, th *material.Theme) layout.Widget {
	b := material.Button(th, btn, "")
	b.CornerRadius = unit.Dp(9)
	b.Background = col
	b.Inset = layout.Inset{Left: unit.Dp(10), Right: unit.Dp(10), Bottom: unit.Dp(1), Top: unit.Dp(1)}

	return func(gtx C) D {
		return layout.Inset{Left: unit.Dp(5)}.Layout(gtx, b.Layout)
	}
}

func about(th *material.Theme, ver string) layout.Widget {
	l := material.Label(th, unit.Sp(14), "FastCopyPast version "+ver)
	l.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 120}
	l.Font.Typeface = "SteticaRegular"

	return func(gtx C) D {
		return l.Layout(gtx)
	}
}
