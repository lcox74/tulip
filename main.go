package main

import (
	"github.com/lcox74/tulip/apps/xkcd"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(400, 240, "Tulip v0.1.0")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// Start App
	app := xkcd.XkcdApp{}
	app.Init()

	for !rl.WindowShouldClose() {
		app.Update()

		rl.BeginDrawing()
		app.Draw()
		rl.EndDrawing()
	}
}
