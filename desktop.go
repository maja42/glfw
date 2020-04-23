// +build !js

package glfw

import "C"
import (
	"io"
	"os"
	"strings"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var enqueue func(blocking bool, fn func())
var contextWatcher ContextWatcher

type RenderThread interface {
	Enqueue(blocking bool, fn func())
}

// Init initializes the library.
//
// Expects a render thread to execute commands.
// A valid ContextWatcher must be provided. It gets notified when context becomes current or detached.
// It should be provided by the GL bindings you are using, so you can do glfw.Init(renderThread, gl.ContextWatcher).
func Init(renderThread RenderThread, cw ContextWatcher) error {
	contextWatcher = cw
	enqueue = renderThread.Enqueue

	var err error
	enqueue(true, func() {
		err = glfw.Init()
	})
	return err
}

// Terminate destroys all remaining windows, frees any allocated resources and de-initializes the library.
func Terminate() {
	enqueue(false, func() {
		glfw.Terminate()
	})
}

// CreateWindow creates a window and its associated context. Most of the options
// controlling how the window and its context should be created are specified
// through Hint.
func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) (*Window, error) {
	var m *glfw.Monitor
	if monitor != nil {
		m = monitor.Monitor
	}
	var s *glfw.Window
	if share != nil {
		s = share.Window
	}

	var w *glfw.Window
	var err error
	enqueue(true, func() {
		w, err = glfw.CreateWindow(width, height, title, m, s)
	})
	if err != nil {
		return nil, err
	}

	window := &Window{Window: w}

	return window, err
}

// SwapInterval sets the swap interval for the current context, i.e. the number
// of screen updates to wait before swapping the buffers of a window and
// returning from SwapBuffers. This is sometimes called
// 'vertical synchronization', 'vertical retrace synchronization' or 'vsync'.
func SwapInterval(interval int) {
	enqueue(false, func() {
		glfw.SwapInterval(interval)
	})
}

// MakeContextCurrent makes the context of the window current.
func (w *Window) MakeContextCurrent() {
	enqueue(false, func() {
		w.Window.MakeContextCurrent()
		// In reality, context is available on each platform via GetGLXContext, GetWGLContext, GetNSGLContext, etc.
		// Pretend it is not available and pass nil, since it's not actually needed at this time.
		contextWatcher.OnMakeCurrent(nil)
	})
}

func DetachCurrentContext() {
	enqueue(false, func() {
		glfw.DetachCurrentContext()
		contextWatcher.OnDetach()
	})
}

func (w *Window) SwapBuffers() {
	enqueue(false, func() {
		w.Window.SwapBuffers()
	})
}

func (w *Window) Destroy() {
	enqueue(false, w.Window.Destroy)
}

func (w *Window) SetTitle(title string) {
	enqueue(false, func() {
		w.Window.SetTitle(title)
	})
}

func (w *Window) SetPos(xpos, ypos int) {
	enqueue(false, func() {
		w.Window.SetPos(xpos, ypos)
	})
}

func (w *Window) SetSize(width, height int) {
	enqueue(false, func() {
		w.Window.SetSize(width, height)
	})
}

func (w *Window) GetContentScale() (float32, float32) {
	var x, y float32
	enqueue(true, func() {
		x, y = w.Window.GetContentScale()
	})
	return x, y
}

func (w *Window) GetOpacity() float32 {
	var o float32
	enqueue(true, func() {
		o = w.Window.GetOpacity()
	})
	return o
}

func (w *Window) SetOpacity(opacity float32) {
	enqueue(false, func() {
		w.Window.SetOpacity(opacity)
	})
}

func (w *Window) Iconify() {
	enqueue(false, w.Window.Iconify)
}

func (w *Window) Restore() {
	enqueue(false, w.Window.Restore)
}

func (w *Window) Show() {
	enqueue(false, w.Window.Show)
}

func (w *Window) Hide() {
	enqueue(false, w.Window.Hide)
}

// SetAttrib function sets the value of an attribute of the specified window.
//
// The supported attributes are Decorated, Resizeable, Floating and AutoIconify.
//
// Some of these attributes are ignored for full screen windows. The new value
// will take effect if the window is later made windowed.
//
// Some of these attributes are ignored for windowed mode windows. The new value
// will take effect if the window is later made full screen.
//
// This function may only be called from the main thread.
func (w *Window) SetAttrib(attrib Hint, value int) {
	enqueue(false, func() {
		w.Window.SetAttrib(glfw.Hint(attrib), value)
	})
}

// GetAttrib returns an attribute of the window. There are many attributes,
// some related to the window and others to its context.
func (w *Window) GetAttrib(attrib Hint) int {
	var val int
	enqueue(true, func() {
		val = w.Window.GetAttrib(glfw.Hint(attrib))
	})
	return val
}

func (w *Window) SetClipboardString(str string) {
	enqueue(false, func() {
		w.Window.SetClipboardString(str)
	})
}

func (w *Window) GetClipboardString() string {
	var s string
	enqueue(false, func() {
		s = w.Window.GetClipboardString()
	})
	return s
}

type Window struct {
	*glfw.Window
}

type Monitor struct {
	*glfw.Monitor
}

func GetPrimaryMonitor() *Monitor {
	var m *glfw.Monitor
	enqueue(true, func() {
		m = glfw.GetPrimaryMonitor()
	})
	return &Monitor{Monitor: m}
}

func PollEvents() {
	enqueue(false, func() {
		glfw.PollEvents()
	})
}

type CursorPosCallback func(w *Window, xpos float64, ypos float64)

func (w *Window) SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback) {
	wrappedCbfun := func(_ *glfw.Window, xpos float64, ypos float64) {
		cbfun(w, xpos, ypos)
	}

	p := w.Window.SetCursorPosCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type KeyCallback func(w *Window, key Key, scancode int, action Action, mods ModifierKey)

func (w *Window) SetKeyCallback(cbfun KeyCallback) (previous KeyCallback) {
	wrappedCbfun := func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		cbfun(w, Key(key), scancode, Action(action), ModifierKey(mods))
	}

	p := w.Window.SetKeyCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type CharCallback func(w *Window, char rune)

func (w *Window) SetCharCallback(cbfun CharCallback) (previous CharCallback) {
	wrappedCbfun := func(_ *glfw.Window, char rune) {
		cbfun(w, char)
	}

	p := w.Window.SetCharCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type ScrollCallback func(w *Window, xoff float64, yoff float64)

func (w *Window) SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback) {
	wrappedCbfun := func(_ *glfw.Window, xoff float64, yoff float64) {
		cbfun(w, xoff, yoff)
	}

	p := w.Window.SetScrollCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type MouseButtonCallback func(w *Window, button MouseButton, action Action, mods ModifierKey)

func (w *Window) SetMouseButtonCallback(cbfun MouseButtonCallback) (previous MouseButtonCallback) {
	wrappedCbfun := func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		cbfun(w, MouseButton(button), Action(action), ModifierKey(mods))
	}

	p := w.Window.SetMouseButtonCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type FramebufferSizeCallback func(w *Window, width int, height int)

func (w *Window) SetFramebufferSizeCallback(cbfun FramebufferSizeCallback) (previous FramebufferSizeCallback) {
	wrappedCbfun := func(_ *glfw.Window, width int, height int) {
		cbfun(w, width, height)
	}

	p := w.Window.SetFramebufferSizeCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

func (w *Window) GetKey(key Key) Action {
	a := w.Window.GetKey(glfw.Key(key))
	return Action(a)
}

func (w *Window) GetMouseButton(button MouseButton) Action {
	a := w.Window.GetMouseButton(glfw.MouseButton(button))
	return Action(a)
}

func (w *Window) GetInputMode(mode InputMode) int {
	return w.Window.GetInputMode(glfw.InputMode(mode))
}

func (w *Window) SetInputMode(mode InputMode, value int) {
	w.Window.SetInputMode(glfw.InputMode(mode), value)
}

type Key glfw.Key

const (
	KeySpace        = Key(glfw.KeySpace)
	KeyApostrophe   = Key(glfw.KeyApostrophe)
	KeyComma        = Key(glfw.KeyComma)
	KeyMinus        = Key(glfw.KeyMinus)
	KeyPeriod       = Key(glfw.KeyPeriod)
	KeySlash        = Key(glfw.KeySlash)
	Key0            = Key(glfw.Key0)
	Key1            = Key(glfw.Key1)
	Key2            = Key(glfw.Key2)
	Key3            = Key(glfw.Key3)
	Key4            = Key(glfw.Key4)
	Key5            = Key(glfw.Key5)
	Key6            = Key(glfw.Key6)
	Key7            = Key(glfw.Key7)
	Key8            = Key(glfw.Key8)
	Key9            = Key(glfw.Key9)
	KeySemicolon    = Key(glfw.KeySemicolon)
	KeyEqual        = Key(glfw.KeyEqual)
	KeyA            = Key(glfw.KeyA)
	KeyB            = Key(glfw.KeyB)
	KeyC            = Key(glfw.KeyC)
	KeyD            = Key(glfw.KeyD)
	KeyE            = Key(glfw.KeyE)
	KeyF            = Key(glfw.KeyF)
	KeyG            = Key(glfw.KeyG)
	KeyH            = Key(glfw.KeyH)
	KeyI            = Key(glfw.KeyI)
	KeyJ            = Key(glfw.KeyJ)
	KeyK            = Key(glfw.KeyK)
	KeyL            = Key(glfw.KeyL)
	KeyM            = Key(glfw.KeyM)
	KeyN            = Key(glfw.KeyN)
	KeyO            = Key(glfw.KeyO)
	KeyP            = Key(glfw.KeyP)
	KeyQ            = Key(glfw.KeyQ)
	KeyR            = Key(glfw.KeyR)
	KeyS            = Key(glfw.KeyS)
	KeyT            = Key(glfw.KeyT)
	KeyU            = Key(glfw.KeyU)
	KeyV            = Key(glfw.KeyV)
	KeyW            = Key(glfw.KeyW)
	KeyX            = Key(glfw.KeyX)
	KeyY            = Key(glfw.KeyY)
	KeyZ            = Key(glfw.KeyZ)
	KeyLeftBracket  = Key(glfw.KeyLeftBracket)
	KeyBackslash    = Key(glfw.KeyBackslash)
	KeyRightBracket = Key(glfw.KeyRightBracket)
	KeyGraveAccent  = Key(glfw.KeyGraveAccent)
	KeyWorld1       = Key(glfw.KeyWorld1)
	KeyWorld2       = Key(glfw.KeyWorld2)
	KeyEscape       = Key(glfw.KeyEscape)
	KeyEnter        = Key(glfw.KeyEnter)
	KeyTab          = Key(glfw.KeyTab)
	KeyBackspace    = Key(glfw.KeyBackspace)
	KeyInsert       = Key(glfw.KeyInsert)
	KeyDelete       = Key(glfw.KeyDelete)
	KeyRight        = Key(glfw.KeyRight)
	KeyLeft         = Key(glfw.KeyLeft)
	KeyDown         = Key(glfw.KeyDown)
	KeyUp           = Key(glfw.KeyUp)
	KeyPageUp       = Key(glfw.KeyPageUp)
	KeyPageDown     = Key(glfw.KeyPageDown)
	KeyHome         = Key(glfw.KeyHome)
	KeyEnd          = Key(glfw.KeyEnd)
	KeyCapsLock     = Key(glfw.KeyCapsLock)
	KeyScrollLock   = Key(glfw.KeyScrollLock)
	KeyNumLock      = Key(glfw.KeyNumLock)
	KeyPrintScreen  = Key(glfw.KeyPrintScreen)
	KeyPause        = Key(glfw.KeyPause)
	KeyF1           = Key(glfw.KeyF1)
	KeyF2           = Key(glfw.KeyF2)
	KeyF3           = Key(glfw.KeyF3)
	KeyF4           = Key(glfw.KeyF4)
	KeyF5           = Key(glfw.KeyF5)
	KeyF6           = Key(glfw.KeyF6)
	KeyF7           = Key(glfw.KeyF7)
	KeyF8           = Key(glfw.KeyF8)
	KeyF9           = Key(glfw.KeyF9)
	KeyF10          = Key(glfw.KeyF10)
	KeyF11          = Key(glfw.KeyF11)
	KeyF12          = Key(glfw.KeyF12)
	KeyF13          = Key(glfw.KeyF13)
	KeyF14          = Key(glfw.KeyF14)
	KeyF15          = Key(glfw.KeyF15)
	KeyF16          = Key(glfw.KeyF16)
	KeyF17          = Key(glfw.KeyF17)
	KeyF18          = Key(glfw.KeyF18)
	KeyF19          = Key(glfw.KeyF19)
	KeyF20          = Key(glfw.KeyF20)
	KeyF21          = Key(glfw.KeyF21)
	KeyF22          = Key(glfw.KeyF22)
	KeyF23          = Key(glfw.KeyF23)
	KeyF24          = Key(glfw.KeyF24)
	KeyF25          = Key(glfw.KeyF25)
	KeyKP0          = Key(glfw.KeyKP0)
	KeyKP1          = Key(glfw.KeyKP1)
	KeyKP2          = Key(glfw.KeyKP2)
	KeyKP3          = Key(glfw.KeyKP3)
	KeyKP4          = Key(glfw.KeyKP4)
	KeyKP5          = Key(glfw.KeyKP5)
	KeyKP6          = Key(glfw.KeyKP6)
	KeyKP7          = Key(glfw.KeyKP7)
	KeyKP8          = Key(glfw.KeyKP8)
	KeyKP9          = Key(glfw.KeyKP9)
	KeyKPDecimal    = Key(glfw.KeyKPDecimal)
	KeyKPDivide     = Key(glfw.KeyKPDivide)
	KeyKPMultiply   = Key(glfw.KeyKPMultiply)
	KeyKPSubtract   = Key(glfw.KeyKPSubtract)
	KeyKPAdd        = Key(glfw.KeyKPAdd)
	KeyKPEnter      = Key(glfw.KeyKPEnter)
	KeyKPEqual      = Key(glfw.KeyKPEqual)
	KeyLeftShift    = Key(glfw.KeyLeftShift)
	KeyLeftControl  = Key(glfw.KeyLeftControl)
	KeyLeftAlt      = Key(glfw.KeyLeftAlt)
	KeyLeftSuper    = Key(glfw.KeyLeftSuper)
	KeyRightShift   = Key(glfw.KeyRightShift)
	KeyRightControl = Key(glfw.KeyRightControl)
	KeyRightAlt     = Key(glfw.KeyRightAlt)
	KeyRightSuper   = Key(glfw.KeyRightSuper)
	KeyMenu         = Key(glfw.KeyMenu)
)

var keyNames = map[Key]string{
	// Printable characters
	KeyA:            "A",
	KeyB:            "B",
	KeyC:            "C",
	KeyD:            "D",
	KeyE:            "E",
	KeyF:            "F",
	KeyG:            "G",
	KeyH:            "H",
	KeyI:            "I",
	KeyJ:            "J",
	KeyK:            "K",
	KeyL:            "L",
	KeyM:            "M",
	KeyN:            "N",
	KeyO:            "O",
	KeyP:            "P",
	KeyQ:            "Q",
	KeyR:            "R",
	KeyS:            "S",
	KeyT:            "T",
	KeyU:            "U",
	KeyV:            "V",
	KeyW:            "W",
	KeyX:            "X",
	KeyY:            "Y",
	KeyZ:            "Z",
	Key1:            "1",
	Key2:            "2",
	Key3:            "3",
	Key4:            "4",
	Key5:            "5",
	Key6:            "6",
	Key7:            "7",
	Key8:            "8",
	Key9:            "9",
	Key0:            "0",
	KeySpace:        "SPACE",
	KeyMinus:        "MINUS",
	KeyEqual:        "EQUAL",
	KeyLeftBracket:  "LEFT BRACKET",
	KeyRightBracket: "RIGHT BRACKET",
	KeyBackslash:    "BACKSLASH",
	KeySemicolon:    "SEMICOLON",
	KeyApostrophe:   "APOSTROPHE",
	KeyGraveAccent:  "GRAVE ACCENT",
	KeyComma:        "COMMA",
	KeyPeriod:       "PERIOD",
	KeySlash:        "SLASH",
	KeyWorld1:       "WORLD 1",
	KeyWorld2:       "WORLD 2",
	// Function keys
	KeyEscape:       "ESCAPE",
	KeyF1:           "F1",
	KeyF2:           "F2",
	KeyF3:           "F3",
	KeyF4:           "F4",
	KeyF5:           "F5",
	KeyF6:           "F6",
	KeyF7:           "F7",
	KeyF8:           "F8",
	KeyF9:           "F9",
	KeyF10:          "F10",
	KeyF11:          "F11",
	KeyF12:          "F12",
	KeyF13:          "F13",
	KeyF14:          "F14",
	KeyF15:          "F15",
	KeyF16:          "F16",
	KeyF17:          "F17",
	KeyF18:          "F18",
	KeyF19:          "F19",
	KeyF20:          "F20",
	KeyF21:          "F21",
	KeyF22:          "F22",
	KeyF23:          "F23",
	KeyF24:          "F24",
	KeyF25:          "F25",
	KeyUp:           "UP",
	KeyDown:         "DOWN",
	KeyLeft:         "LEFT",
	KeyRight:        "RIGHT",
	KeyLeftShift:    "LEFT SHIFT",
	KeyRightShift:   "RIGHT SHIFT",
	KeyLeftControl:  "LEFT CONTROL",
	KeyRightControl: "RIGHT CONTROL",
	KeyLeftAlt:      "LEFT ALT",
	KeyRightAlt:     "RIGHT ALT",
	KeyTab:          "TAB",
	KeyEnter:        "ENTER",
	KeyBackspace:    "BACKSPACE",
	KeyInsert:       "INSERT",
	KeyDelete:       "DELETE",
	KeyPageUp:       "PAGE UP",
	KeyPageDown:     "PAGE DOWN",
	KeyHome:         "HOME",
	KeyEnd:          "END",
	KeyKP0:          "KEYPAD 0",
	KeyKP1:          "KEYPAD 1",
	KeyKP2:          "KEYPAD 2",
	KeyKP3:          "KEYPAD 3",
	KeyKP4:          "KEYPAD 4",
	KeyKP5:          "KEYPAD 5",
	KeyKP6:          "KEYPAD 6",
	KeyKP7:          "KEYPAD 7",
	KeyKP8:          "KEYPAD 8",
	KeyKP9:          "KEYPAD 9",
	KeyKPDivide:     "KEYPAD DIVIDE",
	KeyKPMultiply:   "KEYPAD MULTPLY",
	KeyKPSubtract:   "KEYPAD SUBTRACT",
	KeyKPAdd:        "KEYPAD ADD",
	KeyKPDecimal:    "KEYPAD DECIMAL",
	KeyKPEqual:      "KEYPAD EQUAL",
	KeyKPEnter:      "KEYPAD ENTER",
	KeyPrintScreen:  "PRINT SCREEN",
	KeyNumLock:      "NUM LOCK",
	KeyCapsLock:     "CAPS LOCK",
	KeyScrollLock:   "SCROLL LOCK",
	KeyPause:        "PAUSE",
	KeyLeftSuper:    "LEFT SUPER",
	KeyRightSuper:   "RIGHT SUPER",
	KeyMenu:         "MENU",
}

func (k Key) String() string {
	name, ok := keyNames[k]
	if !ok {
		return "UNKNOWN"
	}
	return name
}

type MouseButton glfw.MouseButton

const (
	MouseButton1 = MouseButton(glfw.MouseButton1)
	MouseButton2 = MouseButton(glfw.MouseButton2)
	MouseButton3 = MouseButton(glfw.MouseButton3)

	MouseButtonLeft   = MouseButton(glfw.MouseButtonLeft)
	MouseButtonRight  = MouseButton(glfw.MouseButtonRight)
	MouseButtonMiddle = MouseButton(glfw.MouseButtonMiddle)
)

func (b MouseButton) String() string {
	switch b {
	case MouseButtonLeft:
		return "LEFT"
	case MouseButtonRight:
		return "RIGHT"
	case MouseButtonMiddle:
		return "MIDDLE"
	default:
		return "UNKNOWN"
	}
}

type Action glfw.Action

const (
	Release = Action(glfw.Release)
	Press   = Action(glfw.Press)
	Repeat  = Action(glfw.Repeat)
)

func (a Action) String() string {
	switch a {
	case Press:
		return "PRESSED"
	case Release:
		return "RELEASED"
	case Repeat:
		return "REPEATED"
	default:
		return "UNKNOWN"
	}
}

type InputMode int

const (
	CursorMode             = InputMode(glfw.CursorMode)
	StickyKeysMode         = InputMode(glfw.StickyKeysMode)
	StickyMouseButtonsMode = InputMode(glfw.StickyMouseButtonsMode)
)

const (
	CursorNormal   = int(glfw.CursorNormal)
	CursorHidden   = int(glfw.CursorHidden)
	CursorDisabled = int(glfw.CursorDisabled)
)

type ModifierKey int

const (
	ModShift   = ModifierKey(glfw.ModShift)
	ModControl = ModifierKey(glfw.ModControl)
	ModAlt     = ModifierKey(glfw.ModAlt)
	ModSuper   = ModifierKey(glfw.ModSuper)
)

func (m ModifierKey) String() string {
	if m == 0 {
		return "[]"
	}
	str := []string{}
	if m&ModShift != 0 {
		str = append(str, "SHIFT")
	}
	if m&ModControl != 0 {
		str = append(str, "CONTROL")
	}
	if m&ModAlt != 0 {
		str = append(str, "ALT")
	}
	if m&ModSuper != 0 {
		str = append(str, "SUPER")
	}
	return "[" + strings.Join(str, ",") + "]"
}

// Open opens a named asset. It's the caller's responsibility to close it when done.
//
// For now, assets are read directly from the current working directory.
func Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

// ---

func WaitEvents() {
	enqueue(true, func() {
		glfw.WaitEvents()
	})
}

func PostEmptyEvent() {
	glfw.PostEmptyEvent()
}

func DefaultWindowHints() {
	enqueue(false, func() {
		glfw.DefaultWindowHints()
	})
}

type CloseCallback func(w *Window)

func (w *Window) SetCloseCallback(cbfun CloseCallback) (previous CloseCallback) {
	wrappedCbfun := func(_ *glfw.Window) {
		cbfun(w)
	}

	p := w.Window.SetCloseCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

// MaximizeCallback is the function signature for window maximize callback functions.
type MaximizeCallback func(w *Window, iconified bool)

// SetMaximizeCallback sets the maximization callback of the specified window,
// which is called when the window is maximized or restored.
//
// This function must only be called from the main thread.
func (w *Window) SetMaximizeCallback(cbfun MaximizeCallback) MaximizeCallback {
	wrappedCbfun := func(_ *glfw.Window, iconified bool) {
		cbfun(w, iconified)
	}

	p := w.Window.SetMaximizeCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

// ContentScaleCallback is the function signature for window content scale
// callback functions.
type ContentScaleCallback func(w *Window, x, y float32)

// SetContentScaleCallback function sets the window content scale callback of
// the specified window, which is called when the content scale of the specified
// window changes.
//
// This function must only be called from the main thread.
func (w *Window) SetContentScaleCallback(cbfun ContentScaleCallback) ContentScaleCallback {
	wrappedCbfun := func(_ *glfw.Window, x, y float32) {
		cbfun(w, x, y)
	}

	p := w.Window.SetContentScaleCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type RefreshCallback func(w *Window)

func (w *Window) SetRefreshCallback(cbfun RefreshCallback) (previous RefreshCallback) {
	wrappedCbfun := func(_ *glfw.Window) {
		cbfun(w)
	}

	p := w.Window.SetRefreshCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type SizeCallback func(w *Window, width int, height int)

func (w *Window) SetSizeCallback(cbfun SizeCallback) (previous SizeCallback) {
	wrappedCbfun := func(_ *glfw.Window, width int, height int) {
		cbfun(w, width, height)
	}

	p := w.Window.SetSizeCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type CursorEnterCallback func(w *Window, entered bool)

func (w *Window) SetCursorEnterCallback(cbfun CursorEnterCallback) (previous CursorEnterCallback) {
	wrappedCbfun := func(_ *glfw.Window, entered bool) {
		cbfun(w, entered)
	}

	p := w.Window.SetCursorEnterCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type PosCallback func(w *Window, xpos int, ypos int)

func (w *Window) SetPosCallback(cbfun PosCallback) (previous PosCallback) {
	wrappedCbfun := func(_ *glfw.Window, xpos int, ypos int) {
		cbfun(w, xpos, ypos)
	}

	p := w.Window.SetPosCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type FocusCallback func(w *Window, focused bool)

func (w *Window) SetFocusCallback(cbfun FocusCallback) (previous FocusCallback) {
	wrappedCbfun := func(_ *glfw.Window, focused bool) {
		cbfun(w, focused)
	}

	p := w.Window.SetFocusCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type IconifyCallback func(w *Window, iconified bool)

func (w *Window) SetIconifyCallback(cbfun IconifyCallback) (previous IconifyCallback) {
	wrappedCbfun := func(_ *glfw.Window, iconified bool) {
		cbfun(w, iconified)
	}

	p := w.Window.SetIconifyCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type DropCallback func(w *Window, names []string)

func (w *Window) SetDropCallback(cbfun DropCallback) (previous DropCallback) {
	wrappedCbfun := func(_ *glfw.Window, names []string) {
		cbfun(w, names)
	}

	p := w.Window.SetDropCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}
