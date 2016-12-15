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
	glfw "github.com/go-gl/glfw3/v3.2/glfw"
	mgl32 "github.com/go-gl/mathgl/mgl32"
	mgl "github.com/go-gl/mathgl/mgl64"
	"math"
	"runtime"
	"time"
)

const (
	WindowWidth   = 600
	WindowHeight  = 800
	WindowTitle   = "Cylinoid"
	TimePerUpdate = time.Duration(time.Second / 60.0)
)

var gPause = false
var gPaddle *Paddle = nil
var gBall *Ball = nil
var gBlocks []*Block
var gCamPos mgl.Vec3
var gLevelWidth float64

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

	if action == glfw.Press {
		inc := float64(0.05)
		switch key {
		case glfw.KeyW:
			gCamPos[1] += inc
		case glfw.KeyA:
			gCamPos[0] -= inc
		case glfw.KeyS:
			gCamPos[1] -= inc
		case glfw.KeyD:
			gCamPos[0] += inc
		case glfw.KeyQ:
			gCamPos[2] += inc
		case glfw.KeyE:
			gCamPos[2] -= inc
		}
	}

}

func PopulateBlocks(sceneSize mgl.Vec2) {
	// Number of blocks
	horizBlocks := 10
	vertBlocks := 4

	// Padding around blockfield
	var vertStartNorm float64 = .7 // .55
	var horizStartNorm float64 = 0 // .1
	var vertStart float64 = vertStartNorm * sceneSize[1]
	var horizStart float64 = horizStartNorm * sceneSize[0]

	// Amount of space taken up by whole blockfield
	horizSpace := float64(1) //float64(.8)
	vertSpace := float64(.3)

	blockWidth := sceneSize[0] * horizSpace / float64(horizBlocks)
	blockHeight := sceneSize[1] * vertSpace / float64(vertBlocks)
	blockSize := mgl.Vec2{blockWidth, blockHeight}

	gBlocks = make([]*Block, horizBlocks*vertBlocks)

	color := mgl.Vec3{0, 1, 0}

	for r := 0; r < vertBlocks; r++ {
		posy := float64(r)*blockHeight + vertStart

		for c := 0; c < horizBlocks; c++ {
			posx := float64(c)*blockWidth + horizStart
			gBlocks[r*horizBlocks+c] = MakeBlock(blockSize, mgl.Vec2{posx, posy}, color)
		}
	}

}

func main() {
	// lock glfw/gl calls to a single thread
	runtime.LockOSThread()

	// Initialize glfw
	// glfw.SetErrorCallback(glfwErrorCallback)

	if glfw.Init() != nil {
		panic("Failed to initialize GLFW")
	}
	defer glfw.Terminate()

	// Open glfw window, with GL4.1 context
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, 0)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, WindowTitle, nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.SetKeyCallback(glfwKeyCallback)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	InitGL()

	height := float64(2)
	width := height * float64(WindowWidth) / float64(WindowHeight)
	gLevelWidth = width
	stageSize := mgl.Vec2{width, height}

	gPaddle = MakePaddle(0.4, stageSize)
	gBall = MakeBall(0.05, mgl.Vec2{width / 2, height / 2})
	PopulateBlocks(stageSize)

	gCamPos = mgl.Vec3{0, 5, 11}
	persp := mgl32.Perspective(45, float32(width/height), 0.1, 100)

	//VP := mgl.Ortho(-width/2, width/2, 0, height*2, -4, 4)

	previousTime := time.Now()

	cameraPos := float64(stageSize[0] / 2)

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
			if len(gBlocks) == 0 {
				PopulateBlocks(stageSize)
			}

			lag -= TimePerUpdate
		}

		// Render once per loop
		ClearScreen()

		// Camera logic
		c := cameraPos
		p := float64(gPaddle.pos[0] + gPaddle.size[0]/2)
		dirLeft := true
		dist := c - p
		stageWidth := float64(stageSize[0])
		maxDist := .2 * stageWidth
		if math.Abs(dist) > stageWidth/2 {
			dist = (stageWidth - math.Abs(dist)) * Sign(dist)
			dirLeft = false
		}
		if math.Abs(dist) > maxDist {
			moveDist := (math.Abs(dist) - maxDist) * Sign(dist)
			if dirLeft {
				cameraPos -= moveDist
			} else {
				cameraPos += moveDist
			}
			if cameraPos > stageWidth {
				cameraPos -= stageWidth
			} else if cameraPos < 0 {
				cameraPos += stageWidth
			}

		}
		//fmt.Println((cameraPos - p) / stageWidth)

		//model := mgl.HomogRotate3DY(-gPaddle.pos[0] / gLevelWidth * 2 * math.Pi)
		//"Model" transformation is the view angle, emulates camera
		//model := mgl.HomogRotate3DY(-float64(cameraPos / stageWidth * 2 * math.Pi))
		model := mgl32.Ident4()
		view := mgl32.LookAt(
			float32(gCamPos[0]), float32(gCamPos[1]), float32(gCamPos[2]),
			0, 3, 0, //gCamPos[0], gCamPos[1], gCamPos[2]+1,
			0, 1, 0)
		MVP := persp.Mul4(view.Mul4(model))

		for _, b := range gBlocks {
			b.renderer.Draw(b.pos, MVP)
		}

		gPaddle.Draw(MVP)
		gBall.renderer.Draw(gBall.pos, MVP)

		window.SwapBuffers()

	}
}
