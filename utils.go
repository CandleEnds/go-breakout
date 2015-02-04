package main

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

func Negate(v mgl.Vec2) mgl.Vec2 {
	return mgl.Vec2{-v[0], -v[1]}
}

func Sign(a float32) float32 {
	if a < 0 {
		return -1
	} else {
		return 1
	}
}

func Sign64(a float64) float64 {
	if a < 0 {
		return -1
	} else {
		return 1
	}
}

func Abs(a float32) float32 {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func MidPt(a, b mgl.Vec2) mgl.Vec2 {
	return a.Add(b).Mul(0.5)
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}
