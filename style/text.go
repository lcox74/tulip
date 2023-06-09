package style

import (
    rl "github.com/gen2brain/raylib-go/raylib"
    "strings"
)


type TextAlignmnet uint8
const (
    AlignLeft TextAlignmnet = iota
    AlignCenter
    AlignRight
)

// DrawText wraps text to fit within a given width but has no max lines.
func DrawText(text string, x, y, width, size int32, color rl.Color, align TextAlignmnet) {
    DrawTextEx(text, x, y, width, 0, size, color, align, false)
}

// DrawTextEx wraps text to fit within a given width and max lines. If debug is
// true then it will draw a rectangle around the text area and circles at the
// corners to indicate a bounding box.
func DrawTextEx(text string, x, y, width, maxLines, size int32, color rl.Color, align TextAlignmnet, debug bool) {
    lines := WrapText(text, width, size)

    if len(lines) > int(maxLines) && maxLines > 0 {
        lines = lines[:maxLines]
    }

    switch align {
    case AlignLeft:
        drawTextLeft(lines, x, y, size, color)
    case AlignCenter:
        drawTextCenter(lines, x, y, width, size, color)
    case AlignRight:
        drawTextRight(lines, x, y, width, size, color)
    }

    if debug {
        rl.DrawRectangleLines(x, y, width, size * int32(len(lines)), rl.Red)
        rl.DrawCircle(x, y, 2, rl.Blue)
        rl.DrawCircle(x, y + size * int32(len(lines)), 2, rl.Blue)
        rl.DrawCircle(x + width, y, 2, rl.Blue)
        rl.DrawCircle(x + width, y + size * int32(len(lines)), 2, rl.Blue)
    }
}

// WrapText wraps text to fit within a given width and returns a slice of lines.
func WrapText(text string, width, size int32) (lines []string) {
    var line string

    for _, section := range strings.Split(text, "\n") {
        line = ""
        for _, word := range strings.Split(section, " ") {
            buildLine := line + word

            for rl.MeasureText(word, size) > width {

                newLine, remainder := breakWord(buildLine, width, size - 1)
                lines = append(lines, newLine + "-")
                buildLine = remainder
                word = remainder
            }

            // Check if word is too long and split across multiple lines
            if rl.MeasureText(buildLine, size) > width {
                lines = append(lines, strings.TrimSpace(line))
                line = word
                continue
            }

            line = buildLine + " "
        }

        lines = append(lines, strings.TrimSpace(line))

    }

	return
}

func drawTextLeft(lines []string, x, y, size int32, color rl.Color) {
    for i, line := range lines {
        rl.DrawText(line, x, y + size * int32(i), size, color)
    }
}

func drawTextCenter(lines []string, x, y, width, size int32, color rl.Color) {
    for i, line := range lines {
        offset := int32(float32(width - rl.MeasureText(line, size)) / 2.0)
        rl.DrawText(line, x + offset, y + size * int32(i), size, color)
    }
}

func drawTextRight(lines []string, x, y, width, size int32, color rl.Color) {
    for i, line := range lines {
        offset := width - rl.MeasureText(line, size)
        rl.DrawText(line, x + offset, y + size * int32(i), size, color)
    }
}

func breakWord(word string, width, size int32) (left, right string) {
    left, right = word, ""

    for rl.MeasureText(left, size) > width {
        left, right = left[:len(left)-1], left[len(left)-1:] + right
    }

    return
}