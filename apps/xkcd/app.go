package xkcd

import (
	"fmt"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lcox74/tulip/constants"
)

const (
	XKCD_URL       = "https://xkcd.com/"
	XKCD_REFRESH   = 60
	XKCD_LINES     = 5
	XKCD_TEXT_SIZE = 20
	XKCD_PADDING   = 20
)

type XkcdApp struct {
	Code  int
	Title string
	Img   string

	comicImage rl.Texture2D
	camera     rl.Camera2D
	lines      []string

	LastFetch time.Time
}

func (app *XkcdApp) Init() {
	app.LastFetch = time.Now()
	app.camera = rl.Camera2D{
		Zoom: 1.0,
	}

	fetchXkcd(app)
	downloadAndProcessImage(app)

}

func (app *XkcdApp) Update() {

	app.lines = wrapText(app.Title, 140, XKCD_TEXT_SIZE, constants.COLOR_FG)
	if len(app.lines) > XKCD_LINES {
		app.lines = app.lines[:XKCD_LINES]
		app.lines = append(app.lines, "...")
	}

	offset := (240 - (XKCD_TEXT_SIZE * (len(app.lines) + 2))) / 2.0
	app.camera.Target.Y = -float32(offset)
}

func (app *XkcdApp) Draw() {
	rl.ClearBackground(constants.COLOR_BG)
	// rl.DrawRectangle(XKCD_PADDING, XKCD_PADDING, 200, 200, constants.COLOR_FG)
	rl.DrawTexture(app.comicImage,
		XKCD_PADDING+int32((200-app.comicImage.Width)/2.0),
		XKCD_PADDING+int32((200-app.comicImage.Height)/2.0),
		rl.White,
	)

	rl.BeginMode2D(app.camera)
	rl.DrawText(fmt.Sprintf("#%d", app.Code), 240, 0, XKCD_TEXT_SIZE, constants.COLOR_FG)
	for i, line := range app.lines {
		rl.DrawText(line, 240, int32(XKCD_TEXT_SIZE*(i+1)), XKCD_TEXT_SIZE, constants.COLOR_FG)
	}
	rl.EndMode2D()

}

func wrapText(text string, width, size int32, color rl.Color) []string {
	words := strings.Split(text, " ")

	lines := []string{}

	var line string
	for _, word := range words {
		newLine := line + word + " "
		if rl.MeasureText(newLine, size) > width {
			lines = append(lines, line)
			line = word + " "
		} else {
			line = newLine
		}
	}
	return append(lines, line)
}
