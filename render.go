package main

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glu"
	mgl "github.com/go-gl/mathgl/mgl32"
	"image"
	"image/png"
	"io"
)

const (
	VertexShaderSource = `#version 120
  
  attribute vec2 position;
  attribute vec2 texCoord;

  uniform vec2 transform;

  varying vec2 TextureCoordOut;

  void main()
  {
    gl_Position = vec4(position+transform, 0.0, 1.0);
    TextureCoordOut = texCoord;
  }`

	FragmentShaderSource = `#version 120
  varying vec2 TextureCoordOut;

  uniform sampler2D Sampler;

  void main()
  {
    gl_FragColor = texture2D(Sampler, TextureCoordOut);
    //gl_FragColor = vec4(0.0, 1.0, 0.0, 1.0);
  }
  `
)

type Rect struct {
	size mgl.Vec2
	pos  mgl.Vec2 //not needed any more
}

func NewRect(sizeX, sizeY float32) Rect {
	return Rect{size: mgl.Vec2{sizeX, sizeY}}
}

func VertexifyRect(r Rect) []float32 {
	lX := float32(0)
	hX := r.size[0]
	lY := float32(0)
	hY := r.size[1]
	return []float32{
		lX, lY, 0, hX, lY, 0, lX, hY, 0,
		hX, lY, 0, hX, hY, 0, lX, hY, 0,
	}
}

func MakeRenderRect(r Rect) *RenderComponent {
	vertices := VertexifyRect(r)
	vao, vbo := makeVertexArrayObject(vertices)
	program := GetDefaultShaderProgram()
	comp := MakeRenderComponent(vao, vbo, program)
	return &comp
}

var gDefaultProgram *gl.Program = nil

func GetDefaultShaderProgram() *gl.Program {
	if gDefaultProgram == nil {
		gDefaultProgram = makeProgram(VertexShaderSource, FragmentShaderSource)
	}
	return gDefaultProgram
}

type RenderComponent struct {
	vao            gl.VertexArray
	vbo            gl.Buffer
	tex            gl.Texture
	program        *gl.Program
	positionAttrib gl.AttribLocation
	uTransformLoc  gl.UniformLocation
}

func MakeRenderComponent(vao gl.VertexArray, vbo gl.Buffer, program *gl.Program) RenderComponent {
	positionAttrib := program.GetAttribLocation("position")
	uTransformLoc := program.GetUniformLocation("transform")
	return RenderComponent{vao, vbo, program, positionAttrib, uTransformLoc}
}

func (r RenderComponent) Draw(pos mgl.Vec2) {
	// global shader
	r.program.Use()
	defer r.program.Unuse()

	// Setup array buffer stuffer
	r.vao.Bind()
	defer r.vao.Unbind()

	r.vbo.Bind(gl.ARRAY_BUFFER)
	defer r.vbo.Unbind(gl.ARRAY_BUFFER)

	r.positionAttrib.AttribPointer(3, gl.FLOAT, false, 0, nil)
	r.positionAttrib.EnableArray()
	defer r.positionAttrib.DisableArray()

	r.uTransformLoc.Uniform2fv(1, []float32{pos[0], pos[1]})

	// Actual draw call
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func InitGL() {
	// Initialize OpenGL, and print version number to console
	gl.Init()
	version := gl.GetString(gl.VERSION)
	fmt.Println("OpenGL version", version)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func ClearScreen() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
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

func createTexture(r io.Reader) (gl.Texture, error) {
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
