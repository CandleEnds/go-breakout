package main

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glu"
	mgl "github.com/go-gl/mathgl/mgl32"
	"image"
	"image/png"
	"os"
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
		lX, lY, 0,
		0, 0,
		hX, lY, 0,
		1, 0,
		lX, hY, 0,
		0, 1,
		hX, lY, 0,
		1, 0,
		hX, hY, 0,
		1, 1,
		lX, hY, 0,
		0, 1,
	}
}

func MakeRenderRect(r Rect, texImg string) *RenderComponent {
	vertices := VertexifyRect(r)
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
		gDefaultProgram = makeProgram(VertexShaderSource, FragmentShaderSource)
	}
	return gDefaultProgram
}

type RenderComponent struct {
	vao            gl.VertexArray
	vbo            gl.Buffer
	program        *gl.Program
	positionAttrib gl.AttribLocation
	uTransformLoc  gl.UniformLocation
	tex            gl.Texture
	uSamplerLoc    gl.UniformLocation
	texCoordAttrib gl.AttribLocation
}

func MakeRenderComponent(vao gl.VertexArray,
	vbo gl.Buffer, tex gl.Texture,
	program *gl.Program) RenderComponent {
	positionAttrib := program.GetAttribLocation("position")
	uTransformLoc := program.GetUniformLocation("transform")
	texCoordAttrib := program.GetAttribLocation("texCoord")
	samplerLoc := program.GetUniformLocation("Sampler")
	return RenderComponent{vao, vbo, program, positionAttrib,
		uTransformLoc, tex, samplerLoc, texCoordAttrib}
}

func (r RenderComponent) Draw(pos mgl.Vec2) {
	// global shader
	gl.Enable(gl.BLEND)
	defer gl.Disable(gl.BLEND)

	r.program.Use()
	defer r.program.Unuse()

	r.uTransformLoc.Uniform2fv(1, []float32{pos[0], pos[1]})
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
	gl.DepthFunc(gl.LESS)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.2, 0.2, 0.2, 1.0)
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
