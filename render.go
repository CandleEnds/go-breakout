package main

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glu"
	mgl "github.com/go-gl/mathgl/mgl32"
	"image"
	"image/png"
	"io/ioutil"
	"os"
)

func VertexifyRect(r mgl.Vec2, depth float32) []float32 {
	lX := float32(0)
	hX := r[0]
	lY := float32(0)
	hY := r[1]
	one := float32(.95)
	zero := float32(.2)
	return []float32{
		lX, lY, depth,
		zero, zero,
		hX, lY, depth,
		one, zero,
		lX, hY, depth,
		zero, one,
		hX, lY, depth,
		one, zero,
		hX, hY, depth,
		one, one,
		lX, hY, depth,
		zero, one,
	}
}

func MakeRenderRect(r mgl.Vec2, depth float32, texImg string) *RenderComponent {
	vertices := VertexifyRect(r, depth)
	vao, vbo := makeVertexArrayObject(vertices)
	program := GetDefaultShaderProgram()
	tex, err := createTexture(texImg)
	if err != nil {
		panic(err)
	}
	comp := MakeRenderComponent(vao, vbo, tex, program)
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
	vao            gl.VertexArray
	vbo            gl.Buffer
	program        *gl.Program
	positionAttrib gl.AttribLocation
	uOffsetLoc     gl.UniformLocation
	uProjLoc       gl.UniformLocation
	tex            gl.Texture
	uSamplerLoc    gl.UniformLocation
	texCoordAttrib gl.AttribLocation
}

func MakeRenderComponent(vao gl.VertexArray,
	vbo gl.Buffer, tex gl.Texture,
	program *gl.Program) RenderComponent {
	positionAttrib := program.GetAttribLocation("position")
	uOffsetLoc := program.GetUniformLocation("offset")
	uProjLoc := program.GetUniformLocation("VP")
	texCoordAttrib := program.GetAttribLocation("texCoord")
	samplerLoc := program.GetUniformLocation("Sampler")
	return RenderComponent{vao, vbo, program, positionAttrib,
		uOffsetLoc, uProjLoc, tex, samplerLoc, texCoordAttrib}
}

func (r RenderComponent) Draw(pos mgl.Vec2, VP mgl.Mat4) {
	// global shader
	gl.Enable(gl.BLEND)
	defer gl.Disable(gl.BLEND)

	r.program.Use()
	defer r.program.Unuse()

	//offset uniform
	r.uOffsetLoc.Uniform2fv(1, []float32{pos[0], pos[1]})
	//View-Projection uniform (VP)
	vpp := [16]float32(VP)
	r.uProjLoc.UniformMatrix4f(false, &vpp)
	//cylinderRadius uniform
	crLoc := r.program.GetUniformLocation("cylinderRadius")
	crLoc.Uniform1f(2)
	//cylinderHeight uniform
	chLoc := r.program.GetUniformLocation("cylinderHeight")
	chLoc.Uniform1f(1)
	//levelWidth uniform
	lwLoc := r.program.GetUniformLocation("levelWidth")
	lwLoc.Uniform1f(gLevelWidth)
	//levelHeight uniform
	lhLoc := r.program.GetUniformLocation("levelHeight")
	lhLoc.Uniform1f(2)

	//texture sampler uniform
	r.uSamplerLoc.Uniform1i(0)

	gl.ActiveTexture(gl.TEXTURE0)
	r.tex.Bind(gl.TEXTURE_2D)
	defer r.tex.Unbind(gl.TEXTURE_2D)

	// Setup array buffer stuffer
	r.vao.Bind()
	defer r.vao.Unbind()

	r.vbo.Bind(gl.ARRAY_BUFFER)
	defer r.vbo.Unbind(gl.ARRAY_BUFFER)

	//5*4 == 5 floats per vert (x,y,z,u,v) by 4 bytes per float
	r.positionAttrib.AttribPointer(3, gl.FLOAT, false, 5*4, nil)
	r.positionAttrib.EnableArray()
	defer r.positionAttrib.DisableArray()

	//offset by x,y,z floats
	r.texCoordAttrib.AttribPointer(2, gl.FLOAT, false, 5*4, uintptr(3*4))
	r.texCoordAttrib.EnableArray()
	defer r.texCoordAttrib.DisableArray()

	// Actual draw call
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func InitGL() {
	// Initialize OpenGL, and print version number to console
	gl.Init()
	version := gl.GetString(gl.VERSION)
	fmt.Println("OpenGL version", version)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
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

func makeVertexArrayObject(vertices []float32) (gl.VertexArray, gl.Buffer) {
	vao := gl.GenVertexArray()
	vao.Bind()
	defer vao.Unbind()

	vbo := gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)
	defer vbo.Unbind(gl.ARRAY_BUFFER)

	// 4 is sizeof float32
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, vertices, gl.STATIC_DRAW)

	return vao, vbo
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
