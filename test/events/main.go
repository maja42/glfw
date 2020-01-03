// events hooks every available callback and outputs their arguments.
package main

import (
	"fmt"
	"time"

	"github.com/maja42/gl/render"
	"github.com/maja42/glfw"
)

var counter int = -1

// getCounter returns event index.
func getCounter() int {
	counter++
	return counter
}

// Window -> Id.
var windowIds = make(map[*glfw.Window]int)

func getWindowId(w *glfw.Window) int {
	return windowIds[w]
}

var startedProcess = time.Now()

// getTime returns time in seconds since process was started.
func getTime() float64 {
	return time.Since(startedProcess).Seconds()
}

func charString(char rune) string {
	return fmt.Sprintf("%#q", char)
}

func PosCallback(w *glfw.Window, x int, y int) {
	fmt.Printf("%08x to %v at %0.3f: Window position: %v %v\n",
		getCounter(), getWindowId(w), getTime(),
		x, y)
}

func SizeCallback(w *glfw.Window, width int, height int) {
	fmt.Printf("%08x to %v at %0.3f: Window size: %v %v\n",
		getCounter(), getWindowId(w), getTime(),
		width, height)
}

func FramebufferSizeCallback(w *glfw.Window, width int, height int) {
	fmt.Printf("%08x to %v at %0.3f: Framebuffer size: %v %v\n",
		getCounter(), getWindowId(w), getTime(),
		width, height)
}

func CloseCallback(w *glfw.Window) {
	fmt.Printf("%08x to %v at %0.3f: Window close\n",
		getCounter(), getWindowId(w), getTime())
}

func RefreshCallback(w *glfw.Window) {
	fmt.Printf("%08x to %v at %0.3f: Window refresh\n",
		getCounter(), getWindowId(w), getTime())
}

func FocusCallback(w *glfw.Window, focused bool) {
	focusedString := map[bool]string{
		true:  "focused",
		false: "defocused",
	}

	fmt.Printf("%08x to %v at %0.3f: Window %s\n",
		getCounter(), getWindowId(w), getTime(),
		focusedString[focused])
}

func IconifyCallback(w *glfw.Window, iconified bool) {
	iconifiedString := map[bool]string{
		true:  "iconified",
		false: "restored",
	}

	fmt.Printf("%08x to %v at %0.3f: Window was %s\n",
		getCounter(), getWindowId(w), getTime(),
		iconifiedString[iconified])
}

func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("%08x to %v at %0.3f: Mouse button %v (%s) (with%s) was %s\n",
		getCounter(), getWindowId(w), getTime(),
		button, button.String(), mods.String(), action.String())
}

func CursorPosCallback(w *glfw.Window, x float64, y float64) {
	fmt.Printf("%08x to %v at %0.3f: Cursor position: %f %f\n",
		getCounter(), getWindowId(w), getTime(),
		x, y)
}

func CursorEnterCallback(w *glfw.Window, entered bool) {
	enteredString := map[bool]string{
		true:  "entered",
		false: "left",
	}

	fmt.Printf("%08x to %v at %0.3f: Cursor %s window\n",
		getCounter(), getWindowId(w), getTime(),
		enteredString[entered])
}

func ScrollCallback(w *glfw.Window, x float64, y float64) {
	fmt.Printf("%08x to %v at %0.3f: Scroll: %0.3f %0.3f\n",
		getCounter(), getWindowId(w), getTime(),
		x, y)
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("%08x to %v at %0.3f: Key 0x%04x Scancode 0x%04x (%s) (with%s) was %s\n",
		getCounter(), getWindowId(w), getTime(),
		key, scancode, key.String(), mods.String(), action.String())
}

func CharCallback(w *glfw.Window, char rune) {
	fmt.Printf("%08x to %v at %0.3f: Character 0x%08x (%s) input\n",
		getCounter(), getWindowId(w), getTime(),
		char, charString(char))
}

func CharModsCallback(w *glfw.Window, char rune, mods glfw.ModifierKey) {
	fmt.Printf("%08x to %v at %0.3f: Character 0x%08x (%s) with modifiers (with%s) input\n",
		getCounter(), getWindowId(w), getTime(),
		char, charString(char), mods.String())
}

func DropCallback(w *glfw.Window, names []string) {
	fmt.Printf("%08x to %v at %0.3f: Drop input\n",
		getCounter(), getWindowId(w), getTime())
	for i, name := range names {
		fmt.Printf("  %v: %q\n", i, name)
	}
}

func main() {
	renderThread := render.New()
	defer renderThread.Terminate()

	err := glfw.Init(renderThread, nil)
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	fmt.Println("Library initialized.")

	window, err := glfw.CreateWindow(640, 480, "Event Linter", nil, nil)
	if err != nil {
		panic(err)
	}
	windowIds[window] = 1 // First (and only) window has id 1.

	window.SetPosCallback(PosCallback)
	window.SetSizeCallback(SizeCallback)
	window.SetFramebufferSizeCallback(FramebufferSizeCallback)
	window.SetCloseCallback(CloseCallback)
	window.SetRefreshCallback(RefreshCallback)
	window.SetFocusCallback(FocusCallback)
	window.SetIconifyCallback(IconifyCallback)
	window.SetMouseButtonCallback(MouseButtonCallback)
	window.SetCursorPosCallback(CursorPosCallback)
	window.SetCursorEnterCallback(CursorEnterCallback)
	window.SetScrollCallback(ScrollCallback)
	window.SetKeyCallback(KeyCallback)
	window.SetCharCallback(CharCallback)
	window.SetDropCallback(DropCallback)

	fmt.Println("Main loop starting.")

	for !window.ShouldClose() {
		glfw.WaitEvents()
	}
}
