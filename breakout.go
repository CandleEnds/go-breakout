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
	mgl "github.com/go-gl/mathgl/mgl32"
	"runtime"
	"time"
)

const (
	WindowWidth   = 800
	WindowHeight  = 600
	WindowTitle   = "App"
	TimePerUpdate = time.Duration(1.0 / 60.0 * float32(time.Second))
)

var gPause = false
var gPaddle *Paddle = nil
var gBall *Ball = nil
var gVP *mgl.Mat4 = nil
var gBlocks []*Block

func glfwErrorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func glfwKeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if gPaddle != nil &&
		gPaddle.GetController()(gPaddle, key, scancode, action, mods) {
		return
	}

	if action == glfw.Press && key == glfw.KeyP {
		gPause = !gPause
	}
}

func PopulateBlocks(sceneSize mgl.Vec2) {
	// Number of blocks
	horizBlocks := 8
	vertBlocks := 4

	// Padding around blockfield
	var vertStart float32 = .55 * sceneSize[1]
	var horizStart float32 = .1 * sceneSize[0]

	// Amount of space taken up by whole blockfield
	horizSpace := float32(.8)
	vertSpace := float32(.3)

	blockWidth := sceneSize[0] * horizSpace / float32(horizBlocks)
	blockHeight := sceneSize[1] * vertSpace / float32(vertBlocks)
	blockSize := mgl.Vec2{blockWidth, blockHeight}

	gBlocks = make([]*Block, horizBlocks*vertBlocks)

	color := mgl.Vec3{0, 1, 0}

	for r := 0; r < vertBlocks; r++ {
		posy := float32(r)*blockHeight + vertStart

		for c := 0; c < horizBlocks; c++ {
			posx := float32(c)*blockWidth + horizStart
			gBlocks[r*horizBlocks+c] = MakeBlock(blockSize, mgl.Vec2{posx, posy}, color)
		}
	}

}

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

	height := float32(2)
	width := height * float32(WindowWidth) / float32(WindowHeight)
	stageSize := mgl.Vec2{width, height}

	gPaddle = MakePaddle(0.6, stageSize)
	gBall = MakeBall(0.05, mgl.Vec2{width / 2, height / 2})
	PopulateBlocks(stageSize)

	VP := mgl.Ortho(0, width, 0, height, -4, 4)
	gVP = &VP

	previousTime := time.Now()
	var lag time.Duration
	for !window.ShouldClose() {
		glfw.PollEvents()

		// Escape key press closes window
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		currentTime := time.Now()
		elapsed := currentTime.Sub(previousTime)
		previousTime = currentTime

		if gPause {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		lag += elapsed

		// Constant time-step updates
		for lag >= TimePerUpdate {

			gPaddle.Update(stageSize)
			gBall.Update(stageSize)
			//update blocks?

			// Collision handling
			var colliders []Collider
			// ball is dynamic, others are static
			colliders = append(colliders, gPaddle)
			colliders = append(colliders, gBall)
			for _, b := range gBlocks {
				colliders = append(colliders, b)
			}

			CollideAll(colliders)

			var killBlocks []int
			for index, b := range gBlocks {
				if !b.alive {
					killBlocks = append(killBlocks, index)
				}
			}

			for i := len(killBlocks) - 1; i >= 0; i-- {
				idx := killBlocks[i]
				gBlocks = append(gBlocks[:idx], gBlocks[idx+1:]...)
			}

			lag -= TimePerUpdate
		}

		// Render once per loop
		ClearScreen()

		for _, b := range gBlocks {
			b.renderer.Draw(b.pos, *gVP)
		}

		gPaddle.Draw(*gVP)
		gBall.renderer.Draw(gBall.pos, *gVP)

		window.SwapBuffers()

	}
}
