package main

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Ball struct {
	renderer *RenderComponent
	pos      mgl.Vec2
	speed    float32
	velocity mgl.Vec2
	size     mgl.Vec2
}

func MakeBall(radius float32, position mgl.Vec2) *Ball {

	rect := mgl.Vec2{radius * 2, radius * 2}
	renderComp := MakeRenderRect(rect, 1, "./ball.png")
	var speed float32 = 1 * float32(TimePerUpdate.Seconds())
	velocity := mgl.Vec2{0, -1}.Normalize()
	position[0] -= radius
	position[1] -= radius
	return &Ball{renderComp, position, speed, velocity, rect}
}

func (b *Ball) Update(stageSize mgl.Vec2) {
	b.pos = b.pos.Add(b.velocity.Mul(b.speed))

	if b.pos[0]+b.size[0] > stageSize[0] {
		b.pos[0] = stageSize[0] - b.size[0]
		b.velocity[0] = -b.velocity[0]
	}
	if b.pos[0] < 0 {
		b.pos[0] = 0
		b.velocity[0] = -b.velocity[0]
	}
	if b.pos[1]+b.size[1] > stageSize[1] {
		b.pos[1] = stageSize[1] - b.size[1]
		b.velocity[1] = -b.velocity[1]
	}
	if b.pos[1] < 0 {
		b.pos[1] = stageSize[1] / 2
	}
	b.velocity = b.velocity.Normalize()
}

func (b *Ball) GetPos() mgl.Vec2 {
	return b.pos
}

func (b *Ball) GetSize() mgl.Vec2 {
	return b.size
}
