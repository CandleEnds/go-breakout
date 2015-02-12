package main

import (
	mgl "github.com/go-gl/mathgl/mgl64"
)

func Negate(v mgl.Vec2) mgl.Vec2 {
	return mgl.Vec2{-v[0], -v[1]}
}

func Sign(a float64) float64 {
	if a < 0 {
		return -1
	} else {
		return 1
	}
}

func MidPt(a, b mgl.Vec2) mgl.Vec2 {
	return a.Add(b).Mul(0.5)
}
