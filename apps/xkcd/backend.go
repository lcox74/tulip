package xkcd

import (
	"errors"
	"fmt"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/net/html"

	"github.com/lcox74/tulip/style"
)

var (
	ErrXkcdCantRequest   = errors.New("cant request xkcd")
	ErrXkcdReadResponse  = errors.New("cant read xkcd response")
	ErrXkcdDownloadImage = errors.New("cant download xkcd comic image")
	ErrXkcdParseImage    = errors.New("cant parse xkcd comic image")
)

type xkcdComic struct {
	Code         int
	Title        string
	Img          string

	rawImage	 *rl.Image
	ComicTexture rl.Texture2D
	loaded	   	 bool
}

func EmptyXkcdComic() xkcdComic {
	return xkcdComic{
		Code:  0,
		Title: "Loading...",
		Img:   "",
	}
}

func (comic *xkcdComic) FetchRandom() error {
	// Send GET request
	resp, err := http.Get(XKCD_URL_RANDOM)
	if err != nil {
		fmt.Printf("[XKCD] Can't request XKCD: %s\n", err)
		return ErrXkcdCantRequest
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[XKCD] Can't read XKCD response: %s\n", err)
		return ErrXkcdReadResponse
	}
	defer resp.Body.Close()

	// Parse HTML
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}

	// Parse HTML nodes
	parseHtmlNode(comic, doc)

	return nil
}

// Clears the comic's state so the next comic can be loaded
func (comic *xkcdComic) Clean() {
	if !comic.loaded {
		return
	}

	// Clean resources from memory
	rl.UnloadTexture(comic.ComicTexture)
	
	comic.loaded = false
	comic.rawImage = nil
}

// Load the texture to the Comic, this will set the comic as ready to be 
// displayed.
//
// NOTE: This isn't in the Process() function because we need this to run on
//       the main thread.
func (comic *xkcdComic) LoadTexture() {
	if comic.rawImage == nil || comic.loaded {
		return
	}
	comic.ComicTexture = rl.LoadTextureFromImage(comic.rawImage)

	// Clean resources from memory
	rl.UnloadImage(comic.rawImage)

	// Set comic as loaded and ready to draw
	comic.loaded = true
}

func (comic *xkcdComic) Process(maxSize int) error {
	if comic.Img == "" {
		return nil
	}

	// Download image
	img, err := downloadComic(comic.Img)

	// Process image
	rlImg := rl.NewImageFromImage(img)
	if rlImg == nil {
		fmt.Printf("[XKCD] Failed to process Image: %s\n", err)
		return ErrXkcdParseImage
	}

	// Convert image to black and white as the 
	style.BlackWhiteImage(rlImg)

	// Scale image to fit the bounds of the maxSize square. We want to keep
	// the aspect ratio so we need to scale both the width and height by the
	// same amount.
	scale := 1.0
	maxImgLen := math.Max(float64(rlImg.Width), float64(rlImg.Height))
	if maxImgLen > float64(maxSize) {
		scale = float64(maxSize) / maxImgLen
	}
	rl.ImageResize(rlImg,
		int32(float64(rlImg.Width)*scale),
		int32(float64(rlImg.Height)*scale),
	)

	comic.rawImage = rlImg
	return nil
}

// Parse HTML to find the comic title, image, and alt text
func parseHtmlNode(comic *xkcdComic, n *html.Node) {
	var err error

	// Parse HTML node
	if n.Type == html.ElementNode && n.Data == "meta" {
		if htmlAttrContainsKey(n.Attr, "content") && htmlAttrContainsKey(n.Attr, "property") {
			switch {
			case htmlAttrContains(n.Attr, "property", "og:title"):
				comic.Title = htmlAttrValFromKey(n.Attr, "content")
			case htmlAttrContains(n.Attr, "property", "og:image"):
				comic.Img = htmlAttrValFromKey(n.Attr, "content")
			case htmlAttrContains(n.Attr, "property", "og:url"):
				url := htmlAttrValFromKey(n.Attr, "content")
				parts := strings.Split(strings.Trim(url, "/"), "/")

				// Parse the comic code from the URL
				if comic.Code, err = strconv.Atoi(parts[len(parts)-1]); err != nil {
					fmt.Printf("[XKCD] Can't parse XKCD code: %s\n", err)
				}
			}
		}
	}

	// Recurse through child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseHtmlNode(comic, c)
	}
}

// Download and decode an image from a URL
func downloadComic(url string) (image.Image, error) {
	if url == "" {
		return nil, nil
	}

	// Download image
	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("[XKCD] Failed to fetch Image: %s\n", err)
		return nil, ErrXkcdDownloadImage
	}

	// Decode image
	img, _, err := image.Decode(r.Body)
	if err != nil {
		fmt.Printf("[XKCD] Failed to decode Image: %s\n", err)
		return nil, ErrXkcdParseImage
	}

	return img, nil
}

// Check if an HTML attribute contains a specific key and value
func htmlAttrContains(attrs []html.Attribute, key string, val string) bool {
	for _, a := range attrs {
		if a.Key == key && a.Val == val {
			return true
		}
	}
	return false
}

func htmlAttrContainsKey(attrs []html.Attribute, key string) bool {
	for _, a := range attrs {
		if a.Key == key {
			return true
		}
	}
	return false
}
func htmlAttrValFromKey(attrs []html.Attribute, key string) string {
	for _, a := range attrs {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func limitInt(value *uint8, thresh, minVal, maxVal uint8) {
	if *value < thresh {
		*value = minVal
		return
	}
	*value = maxVal
}