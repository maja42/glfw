// +build !js

package glfw

import "C"
import "github.com/go-gl/glfw/v3.3/glfw"

type Hint int

// Init related hints. (Use with glfw.InitHint)
const (
	JoystickHatButtons  = Hint(glfw.JoystickHatButtons)  // Specifies whether to also expose joystick hats as buttons, for compatibility with earlier versions of GLFW that did not have glfwGetJoystickHats.
	CocoaChdirResources = Hint(glfw.CocoaChdirResources) // Specifies whether to set the current directory to the application to the Contents/Resources subdirectory of the application's bundle, if present.
	CocoaMenubar        = Hint(glfw.CocoaMenubar)        // Specifies whether to create a basic menu bar, either from a nib or manually, when the first window is created, which is when AppKit is initialized.
)

// Window related hints/attributes.
const (
	Focused                = Hint(glfw.Focused)                // Specifies whether the window will be given input focus when created. This hint is ignored for full screen and initially hidden windows.
	Iconified              = Hint(glfw.Iconified)              // Specifies whether the window will be minimized.
	Maximized              = Hint(glfw.Maximized)              // Specifies whether the window is maximized.
	Visible                = Hint(glfw.Visible)                // Specifies whether the window will be initially visible.
	Hovered                = Hint(glfw.Hovered)                // Specifies whether the cursor is currently directly over the content area of the window, with no other windows between. See Cursor enter/leave events for details.
	Resizable              = Hint(glfw.Resizable)              // Specifies whether the window will be resizable by the user.
	Decorated              = Hint(glfw.Decorated)              // Specifies whether the window will have window decorations such as a border, a close widget, etc.
	Floating               = Hint(glfw.Floating)               // Specifies whether the window will be always-on-top.
	AutoIconify            = Hint(glfw.AutoIconify)            // Specifies whether fullscreen windows automatically iconify (and restore the previous video mode) on focus loss.
	CenterCursor           = Hint(glfw.CenterCursor)           // Specifies whether the cursor should be centered over newly created full screen windows. This hint is ignored for windowed mode windows.
	TransparentFramebuffer = Hint(glfw.TransparentFramebuffer) // Specifies whether the framebuffer should be transparent.
	FocusOnShow            = Hint(glfw.FocusOnShow)            // Specifies whether the window will be given input focus when glfwShowWindow is called.
	ScaleToMonitor         = Hint(glfw.ScaleToMonitor)         // Specified whether the window content area should be resized based on the monitor content scale of any monitor it is placed on. This includes the initial placement when the window is created.
)

// Context related hints.
const (
	ClientAPI               = Hint(glfw.ClientAPI)               // Specifies which client API to create the context for. Hard constraint.
	ContextVersionMajor     = Hint(glfw.ContextVersionMajor)     // Specifies the client API version that the created context must be compatible with.
	ContextVersionMinor     = Hint(glfw.ContextVersionMinor)     // Specifies the client API version that the created context must be compatible with.
	ContextRobustness       = Hint(glfw.ContextRobustness)       // Specifies the robustness strategy to be used by the context.
	ContextReleaseBehavior  = Hint(glfw.ContextReleaseBehavior)  // Specifies the release behavior to be used by the context.
	OpenGLForwardCompatible = Hint(glfw.OpenGLForwardCompatible) // Specifies whether the OpenGL context should be forward-compatible. Hard constraint.
	OpenGLDebugContext      = Hint(glfw.OpenGLDebugContext)      // Specifies whether to create a debug OpenGL context, which may have additional error and performance issue reporting functionality. If OpenGL ES is requested, this hint is ignored.
	OpenGLProfile           = Hint(glfw.OpenGLProfile)           // Specifies which OpenGL profile to create the context for. Hard constraint.
	ContextCreationAPI      = Hint(glfw.ContextCreationAPI)      // Specifies which context creation API to use to create the context.
)

// Framebuffer related hints.
const (
	ContextRevision        = Hint(glfw.ContextRevision)
	RedBits                = Hint(glfw.RedBits)                // Specifies the desired bit depth of the default framebuffer.
	GreenBits              = Hint(glfw.GreenBits)              // Specifies the desired bit depth of the default framebuffer.
	BlueBits               = Hint(glfw.BlueBits)               // Specifies the desired bit depth of the default framebuffer.
	AlphaBits              = Hint(glfw.AlphaBits)              // Specifies the desired bit depth of the default framebuffer.
	DepthBits              = Hint(glfw.DepthBits)              // Specifies the desired bit depth of the default framebuffer.
	StencilBits            = Hint(glfw.StencilBits)            // Specifies the desired bit depth of the default framebuffer.
	AccumRedBits           = Hint(glfw.AccumRedBits)           // Specifies the desired bit depth of the accumulation buffer.
	AccumGreenBits         = Hint(glfw.AccumGreenBits)         // Specifies the desired bit depth of the accumulation buffer.
	AccumBlueBits          = Hint(glfw.AccumBlueBits)          // Specifies the desired bit depth of the accumulation buffer.
	AccumAlphaBits         = Hint(glfw.AccumAlphaBits)         // Specifies the desired bit depth of the accumulation buffer.
	AuxBuffers             = Hint(glfw.AuxBuffers)             // Specifies the desired number of auxiliary buffers.
	Stereo                 = Hint(glfw.Stereo)                 // Specifies whether to use stereoscopic rendering. Hard constraint.
	Samples                = Hint(glfw.Samples)                // Specifies the desired number of samples to use for multisampling. Zero disables multisampling.
	SRGBCapable            = Hint(glfw.SRGBCapable)            // Specifies whether the framebuffer should be sRGB capable.
	RefreshRate            = Hint(glfw.RefreshRate)            // Specifies the desired refresh rate for full screen windows. If set to zero, the highest available refresh rate will be used. This hint is ignored for windowed mode windows.
	DoubleBuffer           = Hint(glfw.DoubleBuffer)           // Specifies whether the framebuffer should be double buffered. You nearly always want to use double buffering. This is a hard constraint.
	CocoaGraphicsSwitching = Hint(glfw.CocoaGraphicsSwitching) // Specifies whether to in Automatic Graphics Switching, i.e. to allow the system to choose the integrated GPU for the OpenGL context and move it between GPUs if necessary or whether to force it to always run on the discrete GPU.
	CocoaRetinaFramebuffer = Hint(glfw.CocoaRetinaFramebuffer) // Specifies whether to use full resolution framebuffers on Retina displays.
)

// Naming related hints. (Use with glfw.WindowHintString)
const (
	CocoaFrameNAME  = Hint(glfw.CocoaFrameNAME)  // Specifies the UTF-8 encoded name to use for autosaving the window frame, or if empty disables frame autosaving for the window.
	X11ClassName    = Hint(glfw.X11ClassName)    // Specifies the desired ASCII encoded class parts of the ICCCM WM_CLASS window property.nd instance parts of the ICCCM WM_CLASS window property.
	X11InstanceName = Hint(glfw.X11InstanceName) // Specifies the desired ASCII encoded instance parts of the ICCCM WM_CLASS window property.nd instance parts of the ICCCM WM_CLASS window property.
)

const (
	// These hints are used for WebGL contexts, ignored on desktop.
	PremultipliedAlpha = noopHint
	PreserveDrawingBuffer
	PreferLowPowerToHighPerformance
	FailIfMajorPerformanceCaveat
)

// noopHint is ignored.
const noopHint Hint = -1

func WindowHint(target Hint, hint int) {
	if target == noopHint {
		return
	}

	glfw.WindowHint(glfw.Hint(target), hint)
}
