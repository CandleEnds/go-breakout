package main

import (
	mgl "github.com/go-gl/mathgl/mgl64"
)

type Block struct {
	//For drawing
	renderer *RenderComponent
	//For colliding
	pos   mgl.Vec2
	size  mgl.Vec2
	alive bool
}

func MakeBlock(size, pos mgl.Vec2, color mgl.Vec3) *Block {
	renderComp := MakeRenderRect(size, 0, "./greenblock.png")
	return &Block{renderComp, pos, size, true}
}

func (b *Block) GetPos() mgl.Vec2 {
	return b.pos
}

func (b *Block) GetSize() mgl.Vec2 {
	return b.size
}

func (b *Block) Collided(c Collider, overlap Rect) {
	b.alive = false
}

func (b *Block) ResolveCollision(pv []mgl.Vec2) {

}

func (b *Block) Impulse(v mgl.Vec2) {

}
