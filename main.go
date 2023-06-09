package main

import (
	"github.com/lcox74/tulip/apps/xkcd"
	"github.com/lcox74/tulip/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const VERSON = "Tulip v0.1.0"

func main() {
	rl.InitWindow(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT, VERSON)
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// rl.SetTraceLog(rl.LogDebug)

	// Start App
	app := xkcd.XkcdApp{}
	app.Init()

	for !rl.WindowShouldClose() {
		
		// Update the App's state
		app.Update()

		// Draw the App
		rl.BeginDrawing()
		app.Draw()
		rl.EndDrawing()
	}
}
