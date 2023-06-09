package xkcd

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lcox74/tulip/constants"
	"github.com/lcox74/tulip/style"
)

const (
	XKCD_URL          = "https://xkcd.com/"
	XKCD_URL_RANDOM   = "https://c.xkcd.com/random/comic/"
	XKCD_REFRESH      = 60
	XKCD_MAX_LINES    = 6
	XKCD_LINE_WIDTH   = 140
	XKCD_IMG_MAX_SIZE = 200
	XKCD_TEXT_SIZE    = 20
	XKCD_PADDING      = 20
)

type XkcdApp struct {
	comic [2]xkcdComic

	offset int32
	lines  []string

	LastFetch time.Time
}

func (app *XkcdApp) Init() {
	// Initialize the App's buffered comic
	app.comic = [2]xkcdComic{EmptyXkcdComic(), EmptyXkcdComic()}
	app.processText()

	// Fetch the first comic
	go app.fetchComic()
	app.LastFetch = time.Now()
}

func (app *XkcdApp) fetchComic() {
	newComic := xkcdComic{}
	newComic.FetchRandom()
	newComic.Process(XKCD_IMG_MAX_SIZE)

	// Once processed then set it to the App's buffered commic to eventually
	// swap into the current comic
	app.comic[1] = newComic
}

// Update the App's state
func (app *XkcdApp) Update() {

	// Check if it's time to fetch a new comic
	if time.Since(app.LastFetch).Seconds() > XKCD_REFRESH {
		go app.fetchComic()
		app.LastFetch = time.Now()
	}

	// Check if comic is ready to display
	if app.comic[1].rawImage != nil {
		app.swapComics()
		app.comic[1].Clean()

		// Prepare new comic for display
		app.currentComic().LoadTexture()
		app.processText()
	}
}

// Draw the current comic and text to display
func (app *XkcdApp) Draw() {
	var comic = app.currentComic()

	rl.ClearBackground(constants.COLOR_BG)

	// Display Comic Image
	rl.DrawTexture(comic.ComicTexture,
		XKCD_PADDING+int32((XKCD_IMG_MAX_SIZE-comic.ComicTexture.Width)/2.0),
		XKCD_PADDING+int32((XKCD_IMG_MAX_SIZE-comic.ComicTexture.Height)/2.0),
		rl.White,
	)

	// Vertically Center the comic code and title text using a camera wa
	style.DrawTextEx(fmt.Sprintf("#%d\n%s", comic.Code, comic.Title),
		// Position (X, Y)
		XKCD_IMG_MAX_SIZE + 2 * XKCD_PADDING, app.offset,

		// Bounding Box (Width, MaxLines)
		XKCD_LINE_WIDTH,
		XKCD_MAX_LINES,
		
		// Styling (TextSize, Color, Alignment)
		XKCD_TEXT_SIZE,
		constants.COLOR_FG,
		style.AlignLeft,

		// Debug Mode
		constants.TULIP_DEBUG,
	)

}

func (app *XkcdApp) currentComic() *xkcdComic {
	return &app.comic[0]
}

func (app *XkcdApp) swapComics() {
	app.comic[0], app.comic[1] = app.comic[1], app.comic[0]
}

// Processes the comic code and title to be displayed in a bounding box. If the
// title is too long it will cut it off after a few lines and add an ellipsis.
// The text is also vertically centered calculating an offset. The resulting 
// text is stored in the App's lines variable and looks like this:
//
// #1234
// This is a really
// long title that
// will be cut off
// after a few lines
// ...
func (app *XkcdApp) processText() {
	// Bound the Text to a text box
	app.lines = style.WrapText(
		fmt.Sprintf("#%d\n%s", app.currentComic().Code, app.currentComic().Title),
		XKCD_LINE_WIDTH, XKCD_TEXT_SIZE,
	)

	// Cut off excess lines and add an ellipsis
	if len(app.lines) > XKCD_MAX_LINES {
		app.lines = app.lines[:XKCD_MAX_LINES]
		app.lines = append(app.lines, "...")
	}

	// Vertically Center the comic code and title text using a camera as its
	// easier than trying to calculate the offset
	app.offset = int32((240 - (XKCD_TEXT_SIZE * len(app.lines))) / 2.0)
}
