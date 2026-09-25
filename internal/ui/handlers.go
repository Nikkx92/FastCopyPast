package ui

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"golang.design/x/clipboard"
)

func (ui *UI) Events(gtx layout.Context) {
	for {
		keyEvent, ok := gtx.Event(

			pointer.Filter{
				Target: &ui.a.MousePos,
				Kinds:  pointer.Press,
			},
		)
		if !ok {
			break
		}

		switch ke := keyEvent.(type) {

		case pointer.Event:
			if ke.Kind == pointer.Press {
				ui.a.Window.Perform(system.ActionMove)
			}

		}
	}

	if ui.closeBtn.Clicked(gtx) {
		os.Exit(0)
	}

	if ui.settingBtn.Clicked(gtx) {
		if !ui.SettingsOn {
			ui.SettingsOn = true
			ui.prevWindowSize = gtx.Constraints.Min
			minHeight := gtx.Dp(270)

			ui.a.Window.Option(app.MinSize(unit.Dp(316), unit.Dp(minHeight)))

			if ui.prevWindowSize.Y < minHeight {
				delta := minHeight - ui.prevWindowSize.Y
				ui.a.Window.Option(app.Size(unit.Dp(ui.prevWindowSize.X), unit.Dp(ui.prevWindowSize.Y+delta)))
			}

		} else {
			ui.a.Window.Option(app.MinSize(unit.Dp(316), unit.Dp(113)))
			ui.SettingsOn = false
			ui.a.Window.Option(app.Size(unit.Dp(ui.prevWindowSize.X), unit.Dp(ui.prevWindowSize.Y)))
		}
	}

	if ui.hoverBtn.Clicked(gtx) {
		ui.a.Window.Perform(system.ActionMinimize)
	}

	if ui.aboveAll.Update(gtx) {
		if ui.aboveAll.Value {
			ui.a.Window.Option(app.TopMost(true))
		} else {
			ui.a.Window.Option(app.TopMost(false))
		}
	}

	if ui.activeSwc.Update(gtx) {
		if ui.activeSwc.Value {
			ctx, cancel := context.WithCancel(context.Background())
			changes := clipboard.Watch(ctx, clipboard.FmtText)
			ui.a.CancelRead = cancel
			go ui.reader(ctx, changes)
		} else {
			ui.a.CancelRead()
		}
	}

	//quotes for divider
	for {
		ev, ok := ui.SeparatorField.Update(gtx)
		if !ok {
			break
		}

		if _, ok := ev.(widget.ChangeEvent); ok {
			currentText := ui.SeparatorField.Text()

			if currentText == "''" || currentText == "" {
				ui.SeparatorField.SetText("")
				break
			}

			hasStartQuote := currentText[0] == '\''
			hasEndQuote := currentText[len(currentText)-1] == '\''

			mutated := false
			newText := currentText

			if !hasStartQuote {
				newText = "'" + newText
				mutated = true
			}

			if !hasEndQuote {
				newText = newText + "'"
				mutated = true
			}

			if mutated {
				start, end := ui.SeparatorField.Selection()

				if !hasStartQuote {
					start++
					end++
				}

				ui.SeparatorField.SetText(newText)
				ui.SeparatorField.SetCaret(start, end)
			}

		}
	}

	//kegel set
	for {
		ev, ok := ui.KegelField.Update(gtx)
		if !ok {
			break
		}

		if _, ok := ev.(widget.ChangeEvent); ok {
			i, _ := strconv.Atoi(ui.KegelField.Text())
			ui.bufferTextSize = i
			ui.a.Window.Invalidate()
		}
	}

	if ui.black.Clicked(gtx) {
		ui.bufferTextColor = &black
	}

	if ui.red.Clicked(gtx) {
		ui.bufferTextColor = &red
	}

	if ui.green.Clicked(gtx) {
		ui.bufferTextColor = &green
	}

	if ui.blue.Clicked(gtx) {
		ui.bufferTextColor = &blue
	}

	if ui.purple.Clicked(gtx) {
		ui.bufferTextColor = &purple
	}
}

func (ui *UI) reader(ctx context.Context, changes <-chan clipboard.Data) {
	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-changes:
			if !ok {
				fmt.Println("can't read from buf")
			}

			if ui.BufferField.SelectionLen() > 0 {
				continue
			}

			divider := ""
			if len(ui.SeparatorField.Text()) >= 2 && len(ui.BufferField.Text()) != 0 {
				divider = ui.SeparatorField.Text()[1 : len(ui.SeparatorField.Text())-1]
			}

			newLine := ""
			if ui.newLinePast.Value && len(ui.BufferField.Text()) != 0 {
				newLine = "\n"
			}

			currentContent := ui.BufferField.Text() + newLine + string(data.Bytes) + divider
			ui.BufferField.SetText(currentContent)
			ui.BufferField.SetCaret(len(currentContent), len(currentContent))
			ui.a.Window.Invalidate()

		}
	}
}
