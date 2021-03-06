package main

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"math"
)

type Ball struct {
	renderer *RenderComponent
	pos      mgl.Vec2
	speed    float64
	velocity mgl.Vec2
	size     mgl.Vec2
}

func MakeBall(radius float64, position mgl.Vec2) *Ball {

	rect := mgl.Vec2{radius * 2, radius * 2}
	renderComp := MakeRenderRect(rect, 0, "./ball.png")
	//renderComp := MakeRenderCube(radius*2, "./ball.png")
	var speed float64 = 1.3 * TimePerUpdate.Seconds()
	velocity := mgl.Vec2{.6, -.8}.Normalize()
	position[0] -= radius
	position[1] -= radius
	return &Ball{renderComp, position, speed, velocity, rect}
}

func (b *Ball) Update(stageSize mgl.Vec2) {
	b.pos = b.pos.Add(b.velocity.Mul(b.speed))

	if b.pos[0] > stageSize[0] {
		b.pos[0] -= stageSize[0]
		//b.pos[0] = stageSize[0] - b.size[0]
		//b.velocity[0] = -b.velocity[0]
	}
	if b.pos[0] < 0 {
		b.pos[0] += stageSize[0]
		//b.pos[0] = 0
		//b.velocity[0] = -b.velocity[0]
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
	return b.pos.Add(b.size.Mul(.25))
}

func (b *Ball) GetSize() mgl.Vec2 {
	return b.size.Mul(.5)
}

func (b *Ball) Collided(other Collider, overlap Rect) {

}

//compare two projection vectors for a better fit
//use larger if same direction, else average (arbitrary decision)
func BetterProjVal(a, b float64) float64 {
	if Sign(a) == Sign(b) {
		return Sign(a) * math.Max(math.Abs(a), math.Abs(b))
	} else {
		return (a + b) / 2
	}
}

//pvs = projection vectors
func (b *Ball) ResolveCollision(pvs []mgl.Vec2) {
	var finalProjVec mgl.Vec2
	for _, pv := range pvs {
		finalProjVec[0] = BetterProjVal(finalProjVec[0], pv[0])
		finalProjVec[1] = BetterProjVal(finalProjVec[1], pv[1])
	}

	b.pos = b.pos.Add(finalProjVec)

	finalProjVec = finalProjVec.Normalize()

	if finalProjVec[0] != 0 {
		b.velocity[0] *= -1
	}
	if finalProjVec[1] != 0 {
		b.velocity[1] *= -1
	}
}

func (b *Ball) Impulse(v mgl.Vec2) {
	b.velocity = b.velocity.Add(v)
}
