package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var FONT_SPACING = 1

type Button struct {
	Text         string
	X, Y         int32
	W, H         int32
	FontSize     float32
	Rect         rl.Rectangle
	Callback     func()
	Color        rl.Color
	HoverColor   rl.Color
	IsHovered    bool
	TextColor    rl.Color
	TextSize     rl.Vector2
	TextPosition rl.Vector2
}

func NewButton(text string, x, y int32, callback func(), col, hoverCol rl.Color) *Button {
	w := int32(150)
	h := int32(50)
	fontSize := float32(20)

	textSize := rl.MeasureTextEx(rl.GetFontDefault(), text, fontSize, float32(FONT_SPACING))
	textPosition := rl.NewVector2(
		float32(x)+(float32(w)-textSize.X)/2,
		float32(y)+(float32(h)-textSize.Y)/2,
	)

	return &Button{
		Text:         text,
		X:            x,
		Y:            y,
		W:            w,
		H:            h,
		FontSize:     fontSize,
		Rect:         rl.NewRectangle(float32(x), float32(y), float32(w), float32(h)),
		Callback:     callback,
		Color:        col,
		HoverColor:   hoverCol,
		IsHovered:    false,
		TextColor:    GB_0,
		TextSize:     textSize,
		TextPosition: textPosition,
	}
}

func (b *Button) Update(mousePos rl.Vector2) {
	b.IsHovered = rl.CheckCollisionPointRec(mousePos, b.Rect)
}

func (b *Button) WasClicked() {
	if b.IsHovered && rl.IsMouseButtonPressed(rl.MouseLeftButton) && b.Callback != nil {
		b.Callback()
	}
}

func (b *Button) Draw() {
	drawColor := b.Color
	if b.IsHovered {
		drawColor = b.HoverColor
	}

	rl.DrawRectangle(b.X, b.Y, b.W, b.H, drawColor)
	rl.DrawTextEx(rl.GetFontDefault(), b.Text, b.TextPosition, b.FontSize, float32(FONT_SPACING), b.TextColor)
}

func (b *Button) PrepForDrop() {

}
