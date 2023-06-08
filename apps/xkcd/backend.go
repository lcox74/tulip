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
)

var (
	ErrXkcdCantRequest   = errors.New("cant request xkcd")
	ErrXkcdReadResponse  = errors.New("cant read xkcd response")
	ErrXkcdDownloadImage = errors.New("cant download xkcd comic image")
	ErrXkcdParseImage    = errors.New("cant parse xkcd comic image")
)

// Fetch a random comic from XKCD and store it in the app
func fetchXkcd(app *XkcdApp) error {

	// Send GET request
	resp, err := http.Get(XKCD_URL)
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
	parseXkcdHtmlNode(app, doc)
	return nil
}

// Parse HTML to find the comic title, image, and alt text
func parseXkcdHtmlNode(app *XkcdApp, n *html.Node) {
	var err error

	// Parse HTML node
	if n.Type == html.ElementNode && n.Data == "meta" {
		if htmlAttrContainsKey(n.Attr, "content") && htmlAttrContainsKey(n.Attr, "property") {
			switch {
			case htmlAttrContains(n.Attr, "property", "og:title"):
				app.Title = htmlAttrValFromKey(n.Attr, "content")
			case htmlAttrContains(n.Attr, "property", "og:image"):
				app.Img = htmlAttrValFromKey(n.Attr, "content")
			case htmlAttrContains(n.Attr, "property", "og:url"):
				url := htmlAttrValFromKey(n.Attr, "content")
				parts := strings.Split(strings.Trim(url, "/"), "/")

				// Parse the comic code from the URL
				if app.Code, err = strconv.Atoi(parts[len(parts)-1]); err != nil {
					fmt.Printf("[XKCD] Can't parse XKCD code: %s\n", err)
				}
			}
		}
	}

	// Recurse through child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseXkcdHtmlNode(app, c)
	}
}

func downloadAndProcessImage(app *XkcdApp) error {
	if app.Img == "" {
		return nil
	}

	// Download image
	r, err := http.Get(app.Img)
	if err != nil {
		fmt.Printf("[XKCD] Failed to fetch Image: %s\n", err)
		return ErrXkcdDownloadImage
	}

	// Decode image
	img, _, err := image.Decode(r.Body)
	if err != nil {
		fmt.Printf("[XKCD] Failed to decode Image: %s\n", err)
		return ErrXkcdParseImage
	}

	// Process image
	rlImg := rl.NewImageFromImage(img)
	rl.ImageColorGrayscale(rlImg)

	// Scale image
	scale := 1.0
	maxImgLen := math.Max(float64(rlImg.Width), float64(rlImg.Height))
	if maxImgLen > 200 {
		scale = 200 / maxImgLen
	}
	rl.ImageResize(rlImg,
		int32(float64(rlImg.Width)*scale),
		int32(float64(rlImg.Height)*scale),
	)

	app.comicImage = rl.LoadTextureFromImage(rlImg)

	return nil
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
