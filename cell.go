package main

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
)

type CellStyle struct {
	BGColor color.RGBA
}

func Cell(bgColor color.RGBA) CellStyle {
	return CellStyle{
		BGColor: bgColor,
	}
}

func (c CellStyle) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	min := gtx.Constraints.Min
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			paint.ColorOp{Color: c.BGColor}.Add(gtx.Ops)
			paint.PaintOp{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Min.X - 1),
					Y: float32(gtx.Constraints.Min.Y),
				}},
			}.Add(gtx.Ops)
			return layout.Dimensions{
				Size: gtx.Constraints.Min,
			}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = min
			return w(gtx)
		}),
	)
}
