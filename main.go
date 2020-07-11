package main

import (
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/danbrakeley/sfv"
	"github.com/sqweek/dialog"
)

func main() {
	titleAndVersion := "sfvui v0.1.0"
	exeName := filepath.Base(os.Args[0])

	if len(os.Args) != 2 {
		dialog.Message(
			"%s was passed %d command line arguments, but was expecting 1.\n\nUsage: %s <filename.sfv>",
			exeName, len(os.Args)-1, exeName,
		).Title(titleAndVersion).Error()
		os.Exit(-1)
	}

	sfvFileName := os.Args[1]
	sf, err := sfv.CreateFromFile(sfvFileName)
	if err != nil {
		dialog.Message(
			"%s has encountered an error while reading \"%s\":\n\n%v", exeName, sfvFileName, err,
		).Title(titleAndVersion).Error()
		os.Exit(-1)
	}

	results := sf.Verify()

	go func() {
		defer os.Exit(0)
		w := app.NewWindow(app.Title(titleAndVersion+" - ©2020 Dan Brakeley"), app.Size(unit.Dp(800), unit.Dp(700)))
		if err := mainGio(w, results); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

type gioState struct {
	List *layout.List
}

func NewState() *gioState {
	return &gioState{
		List: &layout.List{
			Axis: layout.Vertical,
		},
	}
}

func mainGio(w *app.Window, sfvResults sfv.VerifyResults) error {
	// init
	th := material.NewTheme(gofont.Collection())
	state := NewState()

	// event loop
	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)

				widgets := make([]layout.Widget, len(sfvResults.Files))

				cellInset := layout.UniformInset(unit.Dp(8))
				for i, v := range sfvResults.Files {
					fe := v
					widgets[i] = func(gtx layout.Context) layout.Dimensions {
						cellColor := color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}
						if fe.ExpectedCRC32 != fe.ActualCRC32 {
							cellColor = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
						}
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(2, func(gtx layout.Context) layout.Dimensions {
								return Cell(cellColor).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return cellInset.Layout(gtx, material.H6(th, fe.Filename).Layout)
								})
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return Cell(cellColor).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return cellInset.Layout(gtx, material.H6(th, fe.ExpectedCRC32).Layout)
								})
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return Cell(cellColor).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return cellInset.Layout(gtx, material.H6(th, fe.ActualCRC32).Layout)
								})
							}),
						)
					}
				}

				state.List.Layout(gtx, len(widgets), func(gtx layout.Context, i int) layout.Dimensions {
					return layout.UniformInset(unit.Dp(1)).Layout(gtx, widgets[i])
				})

				e.Frame(gtx.Ops)
			}
		}
	}
}

func layoutWidget(ctx layout.Context, width, height int) layout.Dimensions {
	return layout.Dimensions{
		Size: image.Point{
			X: width,
			Y: height,
		},
	}
}
