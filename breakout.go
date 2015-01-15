// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Originally put together by github.com/segfault88, but
// I thought it might be useful to somebody else too.

// It took me quite a lot of frustration and messing around
// to get a basic example of glfw3 with modern OpenGL (3.3)
// with shaders etc. working. Hopefully this will save you
// some trouble. Enjoy!

package main

import (
	"fmt"
	glfw "github.com/go-gl/glfw3"
	"runtime"
	"time"
)

const (
	WindowWidth   = 800
	WindowHeight  = 600
	WindowTitle   = "App"
	TimePerUpdate = time.Duration(1.0 / 60.0 * float32(time.Second))
)

func glfwErrorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func glfwKeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if gPaddle != nil {
		gPaddle.GetController()(gPaddle, key, scancode, action, mods)
	}
}

var gPaddle *Paddle = nil

func main() {
	// lock glfw/gl calls to a single thread
	runtime.LockOSThread()

	// Initialize glfw
	glfw.SetErrorCallback(glfwErrorCallback)

	if !glfw.Init() {
		panic("Failed to initialize GLFW")
	}
	defer glfw.Terminate()

	// Open glfw window, with GL2.1 context
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, 0)

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, WindowTitle, nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.SetKeyCallback(glfwKeyCallback)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	InitGL()

	gPaddle = MakePaddle(0.6)

	previousTime := time.Now()
	var lag time.Duration
	//lag := time.Duration{0}
	for !window.ShouldClose() {
		currentTime := time.Now()
		elapsed := currentTime.Sub(previousTime)
		previousTime = currentTime
		lag += elapsed

		glfw.PollEvents()

		// Constant time-step updates
		for lag >= TimePerUpdate {
			gPaddle.Update()
			lag -= TimePerUpdate
		}

		// Render once per loop
		ClearScreen()
		gPaddle.Draw()
		window.SwapBuffers()

		// Escape key press closes window
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}
