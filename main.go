// package main

// import (
// 	"fmt"
// 	"runtime"

// 	"github.com/k4sud0n/ebook-to-pdf/capture"
// 	"github.com/k4sud0n/ebook-to-pdf/convert"
// )

// func main() {
// 	filename := "screenshot.png"

// 	if runtime.GOOS == "darwin" {
// 		err := capture.MacOS(filename, "100,100,400,400")
// 		if err != nil {
// 			fmt.Println("Failed to capture screenshot:", err)
// 			return
// 		}
// 		fmt.Println("Screenshot saved")

//			convert.PDF(filename)
//		} else {
//			fmt.Println("This example only works on macOS.")
//		}
//	}
package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/go-vgo/robotgo"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Initialize the Fyne application
	a := app.New()
	w := a.NewWindow("Mouse Drag Example")

	// Create a label to display the mouse coordinates
	label := canvas.NewText("Drag to move", color.Black)
	label.TextSize = 20
	label.Move(fyne.NewPos(10, 10))

	// Create a rectangle to represent the dragging area
	rect := canvas.NewRectangle(color.Gray{Y: 128}) // Gray rectangle
	rect.SetMinSize(fyne.NewSize(100, 100))        // Initial size of rectangle

	// Create a container with a white background
	background := canvas.NewRectangle(color.White)
	background.Resize(fyne.NewSize(800, 600)) // Set the initial size of the background

	content := container.NewWithoutLayout(background, rect, label)
	w.SetContent(content)

	// Variables to store the mouse drag state
	var dragging bool
	var startX, startY int

	go func() {
		for {
			// Get current mouse coordinates and whether the left button is pressed
			x, y := robotgo.GetMousePos()
			isMousePressed := robotgo.IsMouseDown(robotgo.MouseLeft)

			if isMousePressed {
				if !dragging {
					// Start dragging
					dragging = true
					startX, startY = x, y
					label.Text = fmt.Sprintf("Dragging from: X: %d, Y: %d", startX, startY)
				} else {
					// Update rectangle position and size
					width := x - startX
					height := y - startY
					rect.Move(fyne.NewPos(startX, startY))
					rect.Resize(fyne.NewSize(width, height))
					label.Text = fmt.Sprintf("Dragging: X: %d, Y: %d", x, y)
				}
			} else if dragging {
				// End dragging
				dragging = false
				label.Text = fmt.Sprintf("Dragged to: X: %d, Y: %d", x, y)
			}

			// Refresh the Fyne UI
			w.Canvas().Refresh(rect)
			w.Canvas().Refresh(label)

			// Sleep for a short duration to reduce CPU usage
			time.Sleep(30 * time.Millisecond)
		}
	}()

	// Set window size and show it
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
