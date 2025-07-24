package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	CUP_LANES_WIDTH = 3
	CUP_BASE        = 529
	CUP_TOP         = (CUP_BASE - 40)
	CUP_FONT_SIZE   = 40
)

type Cup struct {
	Level     int
	Exp       int
	NeededExp int
	FillWidth int
	FillY     int
	Lanes     []int
	Rect      rl.Rectangle
	Color     rl.Color
	TextPos   Vec2
}

func (c *Cup) Draw() {
	rl.DrawRectangleLinesEx(c.Rect, 1, GB_1)
	rl.DrawRectangle(int32(c.Rect.X+2), int32(c.Rect.Y), int32(c.FillWidth), int32(44), GB_FADE)
	rl.DrawRectangle(int32(c.Rect.X+2), int32(c.FillY), int32(c.FillWidth), int32(c.Exp), GB_1)

	rl.DrawText(
		fmt.Sprintf("%d", c.Level),
		int32(c.TextPos.X),
		int32(c.Rect.Y),
		CUP_FONT_SIZE,
		GB_0)
}

func (c *Cup) IncreaseExp(amount int) {
	c.Exp += amount
	c.UpdateFillRect()
}

func (c *Cup) UpdateFillRect() {
	c.FillY = CUP_BASE - c.Exp
	if c.FillY == CUP_TOP {
		c.LevelUp()
	}
}

func (c *Cup) GetValue() int {
	// TODO: Add function
	return c.Level
}

func (c *Cup) LevelUp() {
	rl.PlaySound(CUP_UP_2)
	c.Level++
	c.UpdateTextPos()
	c.Exp = 0
	// TODO: Add exp function
	c.NeededExp = 5
	c.FillY = CUP_BASE
}

func (c *Cup) UpdateTextPos() {
	c.TextPos.X = (int(c.Rect.X + (c.Rect.Width / 2) - float32(rl.MeasureText(fmt.Sprintf("%d", c.Level), CUP_FONT_SIZE))))
}
