package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BALL_BASE_VALUE = 5
	LEFT            = -1
	RIGHT           = 1
	BOTTOM_ROW      = 20
)

var (
	BALL_COLOR         = GB_1
	NEXT_MOVE  float32 = 0.25
)

type Ball struct {
	IsActive  bool
	Row       int
	Column    int
	Value     int
	MoveFloat float32
	Position  Vec2
}

func (b *Ball) Draw(texture *rl.Texture2D) {
	if b.IsActive {
		DrawSpriteInCell(*texture, int(b.Position.X), int(b.Position.Y), BALL_COLOR)
	}
}

func (b *Ball) MoveY() {
	if IsPegBelow(b.Position, pegs) {
		rl.PlaySound(PEG_HIT_0)
		if b.Position.X == 0 {
			b.MoveX(RIGHT)
		} else if b.Position.X == 26 {
			b.MoveX(LEFT)
		} else {
			dir := LEFT
			if rand.Intn(2) == RIGHT {
				dir = RIGHT
			}
			b.MoveX(dir)
		}
	}
	b.Position.Y++ // += 1
	if b.Position.Y >= BOTTOM_ROW {
		rl.PlaySound(CUP_HIT_0)
		cup_index := CheckWhichCupBallIn(int(b.Position.X))
		cups[cup_index].IncreaseExp(b.Value)
		player_score += cups[cup_index].GetValue()
		b.IsActive = false
	}
}

func (b *Ball) MoveX(direction int) {
	b.Position.X += direction
}

func (b *Ball) Update() {
	if b.IsActive {
		if b.MoveFloat > 0 {
			b.MoveFloat -= rl.GetFrameTime()
		} else {
			b.MoveY()
			b.MoveFloat = NEXT_MOVE
		}
	}
}
