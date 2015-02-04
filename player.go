package main

import (
	//"fmt"
	glfw "github.com/go-gl/glfw3"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type KeyHandleFunc func(*Paddle, glfw.Key, int, glfw.Action, glfw.ModifierKey) bool

func PaddleHandleKey(paddle *Paddle,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey) bool {
	//direction of key movement, 1 is press, -1 is release
	var d int
	if action == glfw.Press {
		d = 1
	} else if action == glfw.Release {
		d = -1
	} else {
		return false
	}

	if key == glfw.KeyLeft {
		paddle.Move(-d)
	} else if key == glfw.KeyRight {
		paddle.Move(d)
	} else {
		return false
	}

	return true
}

type Paddle struct {
	renderer   *RenderComponent
	controller KeyHandleFunc
	pos        mgl.Vec2
	speed      float32
	velocity   int
	size       mgl.Vec2
}

func MakePaddle(width float32, sceneSize mgl.Vec2) *Paddle {
	size := mgl.Vec2{width, 0.15}
	renderComp := MakeRenderRect(size, 0, "./greenblock.png")
	pos := mgl.Vec2{(sceneSize[0] - width) / 2, 0.05 * sceneSize[1]}
	speed := 1 * TimePerUpdate.Seconds()
	return &Paddle{renderComp, PaddleHandleKey, pos, float32(speed), 0, size}
}

func (p *Paddle) Draw(VP mgl.Mat4) {
	p.renderer.Draw(p.pos, VP)
}

func (p *Paddle) GetController() KeyHandleFunc {
	return p.controller
}

func (p *Paddle) GetPos() mgl.Vec2 {
	return p.pos
}

func (p *Paddle) GetSize() mgl.Vec2 {
	return p.size
}

func (p *Paddle) Update(stageSize mgl.Vec2) {
	p.pos[0] += p.speed * float32(p.velocity)
	if p.pos[0] > stageSize[0] {
		p.pos[0] -= stageSize[0]
		//p.pos[0] = stageSize[0] - p.size[0]
	} else if p.pos[0] < 0 {
		p.pos[0] += stageSize[0]
		//p.pos[0] = 0
	}

	//fmt.Println(p.pos[0])
}

//negative is left, positive is right
func (p *Paddle) Move(dir int) {
	p.velocity += dir
}

func (p *Paddle) Collided(c Collider, overlap Rect) {
	/*
		impulse := mgl.Vec2{0, 0}
		if overlap.Height() > overlap.Width() {
			impulse[0] = 1
		} else {
			center := overlap.Center()[0]
			padcenter := p.pos[0] + p.size[0]/2
			norm := (center - padcenter) / p.size[0] * 2
			impulse[0] = norm
		}
		fmt.Println(impulse[0])
		c.Impulse(impulse)
	*/
}

func (p *Paddle) ResolveCollision(pv []mgl.Vec2) {

}

func (p *Paddle) Impulse(v mgl.Vec2) {

}
