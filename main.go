package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type ScreenshotTool struct {
	*widgets.QWidget
	captureButton *widgets.QPushButton
	overlay       *ScreenshotOverlay
	coordsLabel   *widgets.QLabel
	timesInput    *widgets.QLineEdit
	captureCount  int
	rightArrow    *widgets.QRadioButton
	leftArrow     *widgets.QRadioButton
}

type ScreenshotOverlay struct {
	*widgets.QWidget
	startX, startY int
	endX, endY     int
	isDragging     bool
	parent         *ScreenshotTool
}

func NewScreenshotTool() *ScreenshotTool {
	tool := &ScreenshotTool{
		QWidget: widgets.NewQWidget(nil, 0),
	}
	tool.SetWindowTitle("Screenshot Tool")
	tool.SetFixedSize2(250, 250)

	layout := widgets.NewQVBoxLayout()
	tool.SetLayout(layout)

	inputLayout := widgets.NewQHBoxLayout()
	inputLabel := widgets.NewQLabel2("Number of captures:", nil, 0)
	tool.timesInput = widgets.NewQLineEdit(nil)
	tool.timesInput.SetPlaceholderText("Enter number")
	inputLayout.AddWidget(inputLabel, 0, 0)
	inputLayout.AddWidget(tool.timesInput, 0, 0)
	layout.AddLayout(inputLayout, 0)

	// Add radio buttons for arrow keys
	tool.rightArrow = widgets.NewQRadioButton2("Right Arrow", nil)
	tool.leftArrow = widgets.NewQRadioButton2("Left Arrow", nil)
	layout.AddWidget(tool.rightArrow, 0, 0)
	layout.AddWidget(tool.leftArrow, 0, 0)
	tool.rightArrow.SetChecked(true)

	tool.captureButton = widgets.NewQPushButton2("Capture Screenshot", nil)
	tool.captureButton.ConnectClicked(tool.startCapture)
	layout.AddWidget(tool.captureButton, 0, 0)

	tool.coordsLabel = widgets.NewQLabel2("Coordinates will appear here", nil, 0)
	tool.coordsLabel.SetAlignment(core.Qt__AlignCenter)
	layout.AddWidget(tool.coordsLabel, 0, 0)

	return tool
}

func (t *ScreenshotTool) startCapture(checked bool) {
	times, err := strconv.Atoi(t.timesInput.Text())
	if err != nil || times <= 0 {
		widgets.QMessageBox_Warning(nil, "Invalid Input", "Please enter a valid positive number.", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}

	t.captureCount = times
	t.Hide()
	t.overlay = NewScreenshotOverlay(t)
	t.overlay.Show()
}

func NewScreenshotOverlay(parent *ScreenshotTool) *ScreenshotOverlay {
	overlay := &ScreenshotOverlay{
		QWidget: widgets.NewQWidget(nil, core.Qt__FramelessWindowHint|core.Qt__WindowStaysOnTopHint|core.Qt__Tool),
		parent:  parent,
	}
	overlay.SetWindowTitle("Screenshot Overlay")
	overlay.SetAttribute(core.Qt__WA_TranslucentBackground, true)
	overlay.SetStyleSheet("background-color: transparent;")

	desktop := widgets.QApplication_Desktop()
	screenRect := desktop.ScreenGeometry(desktop.PrimaryScreen())
	overlay.SetGeometry(screenRect)

	overlay.ConnectMousePressEvent(overlay.OnMousePress)
	overlay.ConnectMouseMoveEvent(overlay.OnMouseMove)
	overlay.ConnectMouseReleaseEvent(overlay.OnMouseRelease)
	overlay.ConnectPaintEvent(overlay.OnPaint)
	overlay.ConnectKeyPressEvent(overlay.OnKeyPress)

	return overlay
}

func (o *ScreenshotOverlay) OnMousePress(event *gui.QMouseEvent) {
	o.startX, o.startY = event.GlobalX(), event.GlobalY()
	o.endX, o.endY = o.startX, o.startY
	o.isDragging = true
	o.Update()
}

func (o *ScreenshotOverlay) OnMouseMove(event *gui.QMouseEvent) {
	if o.isDragging {
		o.endX, o.endY = event.GlobalX(), event.GlobalY()
		o.Update()
	}
}

func (o *ScreenshotOverlay) OnMouseRelease(event *gui.QMouseEvent) {
	o.isDragging = false
	go o.captureMultipleTimes()
}

func (o *ScreenshotOverlay) OnPaint(event *gui.QPaintEvent) {
	painter := gui.NewQPainter2(o)
	defer painter.DestroyQPainter()

	painter.SetCompositionMode(gui.QPainter__CompositionMode_SourceOver)
	painter.FillRect2(0, 0, o.Width(), o.Height(), gui.NewQBrush3(gui.NewQColor3(0, 0, 0, 1), core.Qt__SolidPattern))

	if o.isDragging || (o.startX != o.endX && o.startY != o.endY) {
		pen := gui.NewQPen3(gui.NewQColor3(255, 0, 0, 255))
		pen.SetWidth(2)
		painter.SetPen(pen)
		painter.SetBrush(gui.NewQBrush3(gui.NewQColor3(255, 255, 255, 0), core.Qt__SolidPattern))

		x := min(o.startX, o.endX)
		y := min(o.startY, o.endY)
		width := abs(o.endX - o.startX)
		height := abs(o.endY - o.startY)

		painter.DrawRect2(x-o.X(), y-o.Y(), width, height)

		painter.SetPen(gui.NewQPen3(gui.NewQColor3(255, 0, 0, 255)))
		text := fmt.Sprintf("(%d, %d) to (%d, %d)", x, y, x+width, y+height)
		painter.DrawText3(x-o.X()+5, y-o.Y()+20, text)
	}
}

func (o *ScreenshotOverlay) OnKeyPress(event *gui.QKeyEvent) {
	if event.Key() == int(core.Qt__Key_Escape) {
		o.Close()
		o.parent.Show()
	}
}

func (o *ScreenshotOverlay) HandleArrowKey(right bool) {
	var key core.Qt__Key
	if right {
		key = core.Qt__Key_Right
	} else {
		key = core.Qt__Key_Left
	}

	event := gui.NewQKeyEvent(core.QEvent__KeyPress, int(key), core.Qt__NoModifier, "", false, 1)
	core.QCoreApplication_PostEvent(o.QWidget, event, int(core.Qt__NormalEventPriority))
}

func (o *ScreenshotOverlay) captureMultipleTimes() {
	x := min(o.startX, o.endX)
	y := min(o.startY, o.endY)
	width := abs(o.endX - o.startX)
	height := abs(o.endY - o.startY)

	screen := gui.QGuiApplication_PrimaryScreen()

	coordsText := fmt.Sprintf("Captured: (%d, %d) to (%d, %d)\nSize: %dx%d", x, y, x+width, y+height, width, height)
	o.parent.coordsLabel.SetText(coordsText)

	// Create the "img" directory if it doesn't exist
	imgDir := filepath.Join(".", "img")
	if err := os.MkdirAll(imgDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	for i := 0; i < o.parent.captureCount; i++ {
		o.Hide()
		pixmap := screen.GrabWindow(0, x, y, width, height)
		o.Show()

		filename := filepath.Join(imgDir, fmt.Sprintf("screenshot_%d.webp", i+1))

		// Convert QPixmap to QImage
		image := pixmap.ToImage()

		// Save as WebP
		if image.Save(filename, "WEBP", 100) {
			fmt.Printf("Screenshot saved: %s\n", filename)
		} else {
			fmt.Printf("Failed to save screenshot: %s\n", filename)
		}

		if i < o.parent.captureCount-1 {
			// Handle arrow key press
			o.HandleArrowKey(o.parent.rightArrow.IsChecked())
			time.Sleep(500 * time.Millisecond) // Wait for 0.5 second after key press
		}
	}

	o.Close()
	o.parent.Show()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	tool := NewScreenshotTool()
	tool.Show()

	app.Exec()
}
