package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Peg struct {
	Row      int
	Column   int
	Level    int
	Position Vec2
	Texture  rl.Texture2D
	Color    rl.Color
}

func (p *Peg) Draw(texture *rl.Texture2D) {
	DrawSpriteInCell(*texture, int(p.Position.X), int(p.Position.Y), p.Color)
}

func IsPegBelow(position Vec2, pegs []*Peg) bool {
	for _, peg := range pegs {
		if position.X == peg.Position.X && position.Y+1 == peg.Position.Y {
			return true
		}
	}
	return false
}
