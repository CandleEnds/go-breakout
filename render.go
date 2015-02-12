package main

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glu"
	mgl32 "github.com/go-gl/mathgl/mgl32"
	mgl "github.com/go-gl/mathgl/mgl64"
	"image"
	"image/png"
	"io/ioutil"
	"os"
)

//Break into 4 vertices width-wise, so 3 quads, so 6 tris
func VertexifyRect(r mgl.Vec2, d float32) (vertices []float32, indices []uint16) {
	lX := float32(0)
	mX := float32(r[0] / 3)
	nX := float32(r[0] * 2 / 3)
	hX := float32(r[0])
	lY := float32(0)

	hY := float32(r[1])
	h := float32(.95)
	l := float32(.05)
	m := float32(1.0 / 3.0)
	n := float32(2.0 / 3.0)

	vertices = []float32{
		lX, lY, d, l, l,
		mX, lY, d, m, l,
		nX, lY, d, n, l,
		hX, lY, d, h, l,
		lX, hY, d, l, h,
		mX, hY, d, m, h,
		nX, hY, d, n, h,
		hX, hY, d, h, h,
	}
	indices = []uint16{
		0, 5, 4,
		0, 1, 5,
		1, 6, 5,
		1, 2, 6,
		2, 7, 6,
		2, 3, 7,
	}
	return
}

//TODO: Not working yet, maybe something to do with tex coords
func VertexifyCube(size float32) (vertices []float32, indices []uint16) {
	vertices = []float32{
		0, 0, 0, 0, 0,
		0, 0, 1, 0, 1,
		0, 1, 0, 1, 0,
		0, 1, 1, 1, 1,
		1, 0, 0, 0, 0,
		1, 0, 1, 0, 1,
		1, 1, 0, 1, 0,
		1, 1, 1, 1, 1,
	}
	indices = []uint16{
		0, 4, 2,
		2, 4, 6,
		4, 5, 6,
		6, 5, 7,
		3, 2, 6,
		3, 6, 7,
		1, 0, 2,
		1, 2, 3,
		5, 1, 7,
		1, 3, 7,
		1, 5, 4,
		0, 1, 4,
	}
	return
}

func MakeRenderCircle(normal mgl.Vec3) {

}

func MakeRenderRect(r mgl.Vec2, depth float32, texImg string) *RenderComponent {
	vertices, indices := VertexifyRect(r, depth)
	vao, vbo, indexBuffer := makeVertexArrayObject(vertices, indices)
	program := GetDefaultShaderProgram()
	tex, err := createTexture(texImg)
	if err != nil {
		panic(err)
	}
	comp := MakeRenderComponent(vao, vbo, indexBuffer, len(indices), tex, program)
	return &comp
}

func MakeRenderCube(size float32, texImg string) *RenderComponent {
	vertices, indices := VertexifyCube(size)
	vao, vbo, indexBuffer := makeVertexArrayObject(vertices, indices)
	program := GetDefaultShaderProgram()
	tex, err := createTexture(texImg)
	if err != nil {
		panic(err)
	}
	comp := MakeRenderComponent(vao, vbo, indexBuffer, len(indices), tex, program)
	return &comp
}

var gDefaultProgram *gl.Program = nil

func GetDefaultShaderProgram() *gl.Program {
	if gDefaultProgram == nil {
		vSrc := getFileAsString("vert_cylinder.glsl")
		fSrc := getFileAsString("frag_normal.glsl")
		gDefaultProgram = makeProgram(vSrc, fSrc)
	}
	return gDefaultProgram
}

type RenderComponent struct {
	vao         gl.VertexArray
	vbo         gl.Buffer
	indexBuffer gl.Buffer
	numIndices  int
	program     *gl.Program
	tex         gl.Texture
}

func MakeRenderComponent(vao gl.VertexArray, vbo gl.Buffer, indexBuffer gl.Buffer, numIndices int,
	tex gl.Texture, program *gl.Program) RenderComponent {
	return RenderComponent{vao, vbo, indexBuffer, numIndices, program, tex}
}

func (r RenderComponent) Draw(pos mgl.Vec2, VP mgl32.Mat4) {
	// global shader
	gl.Enable(gl.BLEND)
	defer gl.Disable(gl.BLEND)

	r.program.Use()
	defer r.program.Unuse()

	positionAttrib := r.program.GetAttribLocation("position")
	uOffsetLoc := r.program.GetUniformLocation("offset")
	uProjLoc := r.program.GetUniformLocation("VP")
	texCoordAttrib := r.program.GetAttribLocation("texCoord")
	samplerLoc := r.program.GetUniformLocation("Sampler")
	//offset uniform
	uOffsetLoc.Uniform2fv(1, []float32{float32(pos[0]), float32(pos[1])})
	//View-Projection uniform (VP)
	vpp := [16]float32(VP)
	uProjLoc.UniformMatrix4f(false, &vpp)
	//cylinderRadius uniform
	crLoc := r.program.GetUniformLocation("cylinderRadius")
	crLoc.Uniform1f(3)
	//cylinderHeight uniform
	chLoc := r.program.GetUniformLocation("cylinderHeight")
	chLoc.Uniform1f(0.8)
	//levelWidth uniform
	lwLoc := r.program.GetUniformLocation("levelWidth")
	lwLoc.Uniform1f(float32(gLevelWidth))
	//levelHeight uniform
	lhLoc := r.program.GetUniformLocation("levelHeight")
	lhLoc.Uniform1f(3)

	//texture sampler uniform
	samplerLoc.Uniform1i(0)

	gl.ActiveTexture(gl.TEXTURE0)
	r.tex.Bind(gl.TEXTURE_2D)
	defer r.tex.Unbind(gl.TEXTURE_2D)

	// Setup array buffer stuffer
	r.vao.Bind()
	defer r.vao.Unbind()

	r.vbo.Bind(gl.ARRAY_BUFFER)
	defer r.vbo.Unbind(gl.ARRAY_BUFFER)

	r.indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	defer r.indexBuffer.Unbind(gl.ELEMENT_ARRAY_BUFFER)

	//5*4 == 5 floats per vert (x,y,z,u,v) by 4 bytes per float
	positionAttrib.AttribPointer(3, gl.FLOAT, false, 5*4, nil)
	positionAttrib.EnableArray()
	defer positionAttrib.DisableArray()

	//offset by x,y,z floats
	texCoordAttrib.AttribPointer(2, gl.FLOAT, false, 5*4, uintptr(3*4))
	texCoordAttrib.EnableArray()
	defer texCoordAttrib.DisableArray()

	// Actual draw call
	//gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.DrawElements(gl.TRIANGLES, r.numIndices, gl.UNSIGNED_SHORT, nil)
}

func InitGL() {
	// Initialize OpenGL, and print version number to console
	gl.Init()
	version := gl.GetString(gl.VERSION)
	fmt.Println("OpenGL version", version)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.9, 0.9, 0.9, 1.0)
}

func ClearScreen() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func getFileAsString(filename string) string {
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(text)
}

func makeProgram(vertSource, fragSource string) *gl.Program {
	var infoLog string

	vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)
	vertex_shader.Source(vertSource)
	vertex_shader.Compile()
	infoLog = vertex_shader.GetInfoLog()
	if len(infoLog) > 0 {
		panic(infoLog)
	}
	defer vertex_shader.Delete()

	fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragment_shader.Source(fragSource)
	fragment_shader.Compile()
	infoLog = fragment_shader.GetInfoLog()
	if len(infoLog) > 0 {
		panic(infoLog)
	}
	defer fragment_shader.Delete()

	program := gl.CreateProgram()
	program.AttachShader(vertex_shader)
	program.AttachShader(fragment_shader)

	program.Link()
	infoLog = program.GetInfoLog()
	if len(infoLog) > 0 {
		panic(infoLog)
	}

	program.Validate()
	infoLog = program.GetInfoLog()
	if len(infoLog) > 0 {
		panic(infoLog)
	}

	return &program
}

func makeVertexArrayObject(vertices []float32, indices []uint16) (gl.VertexArray, gl.Buffer, gl.Buffer) {
	vao := gl.GenVertexArray()
	vao.Bind()
	defer vao.Unbind()

	indexBuffer := gl.GenBuffer()
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	defer indexBuffer.Unbind(gl.ELEMENT_ARRAY_BUFFER)

	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*2, indices, gl.STATIC_DRAW)

	vbo := gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)
	defer vbo.Unbind(gl.ARRAY_BUFFER)

	// 4 is sizeof float32
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, vertices, gl.STATIC_DRAW)

	return vao, vbo, indexBuffer
}

func createTexture(texFile string) (gl.Texture, error) {
	r, err := os.Open(texFile)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	img, err := png.Decode(r)
	if err != nil {
		return gl.Texture(0), err
	}
	rgbaImg, ok := img.(*image.NRGBA)
	if !ok {
		return gl.Texture(0), errors.New("texture must be an NRGBA image")
	}

	textureId := gl.GenTexture()
	textureId.Bind(gl.TEXTURE_2D)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()
	data := make([]byte, imgWidth*imgHeight*4)
	lineLen := imgWidth * 4
	dest := len(data) - lineLen
	for src := 0; src < len(rgbaImg.Pix); src += rgbaImg.Stride {
		copy(data[dest:dest+lineLen], rgbaImg.Pix[src:src+rgbaImg.Stride])
		dest -= lineLen
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, 4,
		imgWidth, imgHeight, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)
	return textureId, nil
}

func checkGLerror() {
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		string, _ := glu.ErrorString(glerr)
		panic(string)
	}
}
