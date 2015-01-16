package main

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Collider interface {
	GetSize() mgl.Vec2
	GetPos() mgl.Vec2
}

func max(a, b float32) float32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func Collide(c1 Collider, c2 Collider) bool {
	lower1 := c1.GetPos()              //x1, y1
	upper1 := lower1.Add(c1.GetSize()) //x2, y2
	lower2 := c2.GetPos()              //x3, y3
	upper2 := lower2.Add(c2.GetSize()) //x4, y4

	lowerOverlap := mgl.Vec2{} //x5, y5
	lowerOverlap[0] = max(lower1[0], lower2[0])
	lowerOverlap[1] = max(lower1[1], lower2[1])

	upperOverlap := mgl.Vec2{}
	upperOverlap[0] = min(upper1[0], upper2[0])
	upperOverlap[1] = min(upper1[1], upper2[1])

	return (lowerOverlap[0] < upperOverlap[0]) && (lowerOverlap[1] < upperOverlap[1])
}

type Mover interface {
	GetPos() mgl.Vec2
	GetVelocity() mgl.Vec2
	SetPos()
	SetVelocity()
}
