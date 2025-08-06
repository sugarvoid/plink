package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Vec2 struct {
	X int
	Y int
}

type FakeBall struct {
	Position Vec2
}

type CupData struct {
	Level int `json:"level"`
	FillY int `json:"fill_y"`
	Exp   int `json:"exp"`
}

type SaveData struct {
	Score  int64      `json:"score"`
	C_Data [9]CupData `json:"cup_data"`
}

const (
	DEBUG               = false
	SAVE_FILE           = "save.json"
	GAME_TITLE          = "Plink"
	SCREEN_WIDTH        = 960
	SCREEN_HEIGHT       = 540
	NUM_ROW             = 4
	NUM_COL             = 6
	LINE_THICKNESS      = 4.0
	FPS                 = 30
	BANNER_FONT_SIZE    = 38
	BANNER_SCROLL_SPEED = 2
	MAX_BANNER_LEN      = 64
	GRID_COLUMNS        = 27
	GRID_ROWS           = 18
	GRID_OFFSET_X       = 6
	GRID_OFFSET_Y       = 46
	ROW_COUNT           = 8
	COLS_PER_ROW        = 13
	PEG_COUNT           = (ROW_COUNT * COLS_PER_ROW) + 4
	NUM_CUPS            = 9
	MAX_BALLS           = 50
	CELL_SIZE           = 24
)

var (
	GB_0    = rl.NewColor(6, 12, 15, 255)
	GB_1    = rl.NewColor(83, 91, 78, 255)
	GB_2    = rl.NewColor(176, 179, 166, 255)
	GB_3    = rl.NewColor(239, 238, 232, 255)
	GB_FADE = rl.NewColor(6, 12, 15, 10)
)

var (
	PEG_TEXTURE  rl.Texture2D
	BALL_TEXTURE rl.Texture2D
	PEG_HIT_0    rl.Sound
	CUP_HIT_0    rl.Sound
	CUP_UP_2     rl.Sound
)

var (
	player_score int
	save_data    SaveData
)

var (
	BORDER_RECT_1 = rl.Rectangle{X: 0, Y: 0, Width: SCREEN_WIDTH, Height: 44}
	BORDER_RECT_2 = rl.Rectangle{X: 0, Y: 0, Width: SCREEN_WIDTH, Height: SCREEN_HEIGHT}
	BORDER_RECT_3 = rl.Rectangle{X: 0, Y: 40, Width: SCREEN_WIDTH - 300, Height: SCREEN_HEIGHT}
)

var (
	pegs    []*Peg
	balls   []*Ball
	cups    []Cup
	f_balls []*FakeBall
)

var (
	btn_drop *Button
)

func main() {

	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, GAME_TITLE)
	rl.InitAudioDevice()
	rl.SetTargetFPS(FPS)

	rand.NewSource(time.Now().UnixNano())

	defer rl.CloseWindow()

	btn_drop = NewButton("Drop", 54, 480, 670, 50, func() { DropBall() }, GB_1, GB_2)

	rl.SetExitKey(rl.KeyQ)

	PEG_TEXTURE = rl.LoadTexture("res/peg_lines.png")
	BALL_TEXTURE = rl.LoadTexture("res/ball_fill.png")
	//MONOGRAM = rl.LoadFontEx("res/monogram.ttf", 60, nil, 250)
	PEG_HIT_0 = rl.LoadSound("res/peg_hit_0.wav")
	CUP_HIT_0 = rl.LoadSound("res/cup_hit_0.wav")
	CUP_UP_2 = rl.LoadSound("res/cup_up_2.wav")

	InitPegs()
	InitFakeBalls()
	InitCups()
	LoadGame()

	for !rl.WindowShouldClose() {
		Update()
		rl.BeginDrawing()
		rl.ClearBackground(GB_3)
		Draw()
		rl.EndDrawing()
	}

	SaveGame()
	CleanUp()
}

func Update() {
	btn_drop.Update(rl.GetMousePosition())
	for _, b := range balls {
		b.Update()
	}
}

func Draw() {
	DrawLines()
	btn_drop.Draw()
	if DEBUG {
		DrawDebugGrid()
	}

	btn_drop.WasClicked()
	for _, peg := range pegs {
		peg.Draw(&PEG_TEXTURE)
	}

	for _, cup := range cups {
		cup.Draw()
	}

	// TODO: Find better way to hide top of cups
	rl.DrawLine(6, 486, 654, 486, GB_3)

	for _, b := range balls {
		b.Draw(&BALL_TEXTURE)
	}

	rl.DrawText(fmt.Sprintf("Score: %d", player_score), 42, 4, CUP_FONT_SIZE, GB_0)

	for _, fb := range f_balls {
		DrawSpriteInCell(BALL_TEXTURE, fb.Position.X, fb.Position.Y, GB_FADE)
	}
}

func DrawSpriteInCell(texture rl.Texture2D, col int, row int, color rl.Color) {
	cell_x := GRID_OFFSET_X + col*CELL_SIZE
	cell_y := GRID_OFFSET_Y + row*CELL_SIZE

	center_x := cell_x + CELL_SIZE/2
	center_y := cell_y + CELL_SIZE/2

	rl.DrawTexture(texture, int32(center_x-(int(texture.Width/2))), int32(center_y-int(texture.Height/2)), color)
}

func DrawLines() {
	rl.DrawRectangleLinesEx(BORDER_RECT_1, LINE_THICKNESS, GB_0)
	rl.DrawRectangleLinesEx(BORDER_RECT_2, LINE_THICKNESS, GB_0)
	rl.DrawRectangleLinesEx(BORDER_RECT_3, LINE_THICKNESS, GB_0)
}

func DrawDebugGrid() {
	for row := range GRID_ROWS {
		for col := range GRID_COLUMNS {
			// Calculate the position of the rectangle within the grid cell
			x := GRID_OFFSET_X + col*CELL_SIZE
			y := GRID_OFFSET_Y + row*CELL_SIZE

			rl.DrawRectangleLines(int32(x), int32(y), int32(CELL_SIZE), int32(CELL_SIZE), GB_0)
		}
	}
}

func DropBall() {
	lane := rand.Intn(GRID_COLUMNS)

	ball := &Ball{
		IsActive:  true,
		Value:     BALL_BASE_VALUE,
		MoveFloat: 0.5,
		Position:  Vec2{X: lane, Y: 0},
	}

	balls = append(balls, ball)

	if DEBUG {
		fmt.Println(len(balls))
		fmt.Printf("Ball address after adding: %p\n", ball)
		fmt.Printf("Balls x pos: %d\n", int32(ball.Position.X))
	}

	CleanUpBallSlice()
}

func InitCups() {
	for i := range NUM_CUPS {
		new_cup := Cup{
			Rect:      rl.Rectangle{X: float32(6 + (CELL_SIZE * 3 * i)), Y: 485, Width: CELL_SIZE * 3, Height: CELL_SIZE * 1.9},
			FillWidth: (CELL_SIZE * 3) - 4,
			Exp:       0,
			Level:     1,
		}
		new_cup.UpdateTextPos()

		for j := range CUP_LANES_WIDTH {
			new_cup.Lanes = append(new_cup.Lanes, i*3+j)
		}

		cups = append(cups, new_cup)
	}
}

func InitPegs() {
	i := 0
	for row := range ROW_COUNT {
		y := (2 + row*2) // Rows at y = 2, 4, 6, ..., 20
		startX := (0)
		if row%2 == 0 {
			startX = 1
		}

		neededPegs := COLS_PER_ROW
		if row%2 != 0 {
			neededPegs++
		}

		// Loop through columns to place pegs
		for col := range neededPegs {
			peg := &Peg{
				Position: Vec2{X: (startX) + (col * 2), Y: y},
				Color:    GB_1,
			}
			pegs = append(pegs, peg)
			i++
		}
	}
}

func InitFakeBalls() {
	for row := range 20 {
		switch row {
		case 0, 1, 3, 5, 7, 9, 11, 13, 15, 17:
			for col := range 27 {
				f_ball := &FakeBall{
					Position: Vec2{X: col, Y: row},
				}
				f_balls = append(f_balls, f_ball)
			}
		case 2, 4, 6, 8, 10, 12, 14, 16:
			start_x := 0
			if row == 2 || row == 6 || row == 10 || row == 14 {
				start_x = 0
			} else {
				start_x = 1
			}
			for i := start_x; i < 27; i += 2 {
				f_ball := &FakeBall{
					Position: Vec2{X: i, Y: row},
				}
				f_balls = append(f_balls, f_ball)
			}
		}
	}
}

func CleanUpBallSlice() {
	for i := 0; i < len(balls); i++ {
		//for i := range len(balls) {
		if !balls[i].IsActive {
			balls = slices.Delete(balls, i, i+1)
			i--
		}
	}
}

func CheckWhichCupBallIn(lane int) int {
	for cup_index := range NUM_CUPS {
		for lane_index := range CUP_LANES_WIDTH {
			if cups[cup_index].Lanes[lane_index] == lane {
				return cup_index
			}
		}
	}
	return -1
}

func SaveGame() {
	save_data.Score = int64(player_score)
	for i := range NUM_CUPS {
		save_data.C_Data[i].Exp = cups[i].Exp
		save_data.C_Data[i].FillY = cups[i].FillY
		save_data.C_Data[i].Level = cups[i].Level
	}

	json_data, err := json.MarshalIndent(save_data, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(SAVE_FILE, json_data, 0644)
	if err != nil {
		panic(err)
	}
	println("JSON file written successfully!")
}

func LoadGame() {
	json_data, err := os.ReadFile(SAVE_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			SaveGame()
			json_data, err = os.ReadFile(SAVE_FILE)
			if err != nil {
				panic(err)
			}
		}
	}

	var _save_data SaveData
	err = json.Unmarshal(json_data, &_save_data)
	if err != nil {
		panic(err)
	}

	player_score = int(_save_data.Score)

	for i := range NUM_CUPS {
		cups[i].Exp = _save_data.C_Data[i].Exp
		cups[i].FillY = _save_data.C_Data[i].FillY
		cups[i].Level = _save_data.C_Data[i].Level
		cups[i].UpdateTextPos()
	}
}

func CleanUp() {
	rl.UnloadSound(PEG_HIT_0)
	rl.UnloadSound(CUP_HIT_0)
	rl.UnloadSound(CUP_UP_2)
}
