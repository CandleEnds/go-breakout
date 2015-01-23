package main

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Collider interface {
	GetSize() mgl.Vec2
	GetPos() mgl.Vec2
	//IsStatic() bool
	Collided(other Collider)
	ResolveCollision(projVecs []mgl.Vec2)
}

type Rect struct {
	lower mgl.Vec2
	upper mgl.Vec2
}

func (r *Rect) Height() float32 {
	return r.upper[1] - r.lower[1]
}

func (r *Rect) Width() float32 {
	return r.upper[0] - r.lower[0]
}

func (r *Rect) Center() mgl.Vec2 {
	return r.lower.Add(r.upper).Mul(0.5)
}

func CollideAll(colliders []Collider) {
	colls := make(map[Collider][]mgl.Vec2)
	// try collision with all colliders against all other colliders + all colliders
	for i := 0; i < len(colliders); i++ {
		for j := i + 1; j < len(colliders); j++ {
			a := colliders[i]
			b := colliders[j]
			if collides, pv := Collide(a, b); collides {
				a.Collided(b)
				b.Collided(a)
				colls[a] = append(colls[a], pv)
				colls[b] = append(colls[b], Negate(pv))
			}
		}
	}

	for collider, pvs := range colls {
		collider.ResolveCollision(pvs)
	}

}

//Returns projection vector for c1, negate to use for c2
func Collide(c1 Collider, c2 Collider) (bool, mgl.Vec2) {
	lower1 := c1.GetPos()              //x1, y1
	upper1 := lower1.Add(c1.GetSize()) //x2, y2
	lower2 := c2.GetPos()              //x3, y3
	upper2 := lower2.Add(c2.GetSize()) //x4, y4

	lowerOverlap := mgl.Vec2{} //x5, y5
	lowerOverlap[0] = Max(lower1[0], lower2[0])
	lowerOverlap[1] = Max(lower1[1], lower2[1])

	upperOverlap := mgl.Vec2{}
	upperOverlap[0] = Min(upper1[0], upper2[0])
	upperOverlap[1] = Min(upper1[1], upper2[1])

	isGood := (lowerOverlap[0] < upperOverlap[0]) && (lowerOverlap[1] < upperOverlap[1])

	overRect := Rect{lowerOverlap, upperOverlap}

	center1 := MidPt(lower1, upper1)
	center2 := MidPt(lower2, upper2)

	var ymul float32
	var xmul float32

	if center1[1] > center2[1] {
		ymul = 1
	} else {
		ymul = -1
	}

	if center1[0] > center2[0] {
		xmul = 1
	} else {
		xmul = -1
	}

	var projVec mgl.Vec2
	if overRect.Height() < overRect.Width() {
		projVec = mgl.Vec2{0, ymul * overRect.Height()}
	} else {
		projVec = mgl.Vec2{xmul * overRect.Width(), 0}
	}

	return isGood, projVec
}
