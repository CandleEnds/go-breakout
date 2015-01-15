package main

import (
	//"fmt"
	glfw "github.com/go-gl/glfw3"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Player interface {
	Draw()
	GetController() KeyHandleFunc
	GetPos() mgl.Vec2
	Update()
}

type KeyHandleFunc func(*Paddle, glfw.Key, int, glfw.Action, glfw.ModifierKey)

func PaddleHandleKey(paddle *Paddle,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey) {
	//direction of key movement, 1 is press, -1 is release
	var d int
	if action == glfw.Press {
		d = 1
	} else if action == glfw.Release {
		d = -1
	} else {
		return
	}

	if key == glfw.KeyLeft {
		paddle.Move(-d)
	} else if key == glfw.KeyRight {
		paddle.Move(d)
	}
}

type Paddle struct {
	renderer   *RenderComponent
	controller KeyHandleFunc
	pos        mgl.Vec2
	speed      float32
	velocity   int
	width      float32
}

func MakePaddle(width float32) *Paddle {
	if width > 1 {
		width = 1
	}
	rect := NewRect(width, 0.07)
	renderComp := MakeRenderRect(rect, "./paddle.png")
	pos := mgl.Vec2{-width / 2.0, -.95}
	speed := 1 * TimePerUpdate.Seconds()
	return &Paddle{renderComp, PaddleHandleKey, pos, float32(speed), 0, width}
}

func (p *Paddle) Draw() {
	p.renderer.Draw(p.pos)
}

func (p *Paddle) GetController() KeyHandleFunc {
	return p.controller
}

func (p *Paddle) GetPos() mgl.Vec2 {
	return p.pos
}

func (p *Paddle) Update() {
	p.pos[0] += p.speed * float32(p.velocity)
	if p.pos[0]+p.width > 1 {
		p.pos[0] = 1 - p.width
	} else if p.pos[0] < -1 {
		p.pos[0] = -1
	}
}

//negative is left, positive is right
func (p *Paddle) Move(dir int) {
	p.velocity += dir
}
