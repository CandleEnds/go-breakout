package main

import (
	//"errors"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	//"github.com/go-gl/glu"
	mgl32 "github.com/go-gl/mathgl/mgl32"
	mgl "github.com/go-gl/mathgl/mgl64"
	"image"
	"image/draw"
	_ "image/png"
	"io/ioutil"
	"math"
	"os"
	"strings"
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

//make a unit circle centered at the origin with normal Z
func MakeRenderCircle(resolution uint16) (vertices []float32, indices []uint16) {
	vertices = make([]float32, resolution)
	indices = make([]uint16, (resolution-2)*3)
	increment := 2 * math.Pi / float64(resolution)
	for angle := float64(0); angle < float64(2)*math.Pi; angle += increment {
		x := float32(math.Cos(angle))
		y := float32(math.Sin(angle))
		vertices = append(vertices, x)
		vertices = append(vertices, y)
		vertices = append(vertices, 0)
		//tex coords, translate into 0->1 instead of -1->1
		vertices = append(vertices, (x+1)/2)
		vertices = append(vertices, (y+1)/2)
	}
	for idx := uint16(1); idx < resolution-2; idx++ {
		indices = append(indices, 0)
		indices = append(indices, idx)
		indices = append(indices, idx+1)
	}
	return
}

func MakeRenderRect(r mgl.Vec2, depth float32, texImg string) *RenderComponent {
	vertices, indices := VertexifyRect(r, depth)
	vao, vbo, indexBuffer := makeVertexArrayObject(vertices, indices)
	program := GetDefaultShaderProgram()
	tex, err := createTexture(texImg)
	if err != nil {
		panic(err)
	}
	comp := MakeRenderComponent(vao, vbo, indexBuffer, int32(len(indices)), tex, program)
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
	comp := MakeRenderComponent(vao, vbo, indexBuffer, int32(len(indices)), tex, program)
	return &comp
}

var gDefaultProgram uint32 = 0

func GetDefaultShaderProgram() uint32 {
	if gDefaultProgram == 0 {
		vSrc := getFileAsString("vert_cylinder_330.glsl")
		fSrc := getFileAsString("frag_normal_330.glsl")
		var err error
		gDefaultProgram, err = makeProgram(vSrc, fSrc)
		if err != nil {
			panic(err)
		}
	}
	return gDefaultProgram
}

type RenderComponent struct {
	vao         uint32
	vbo         uint32
	indexBuffer uint32
	numIndices  int32
	program     uint32
	tex         uint32
}

func MakeRenderComponent(vao uint32, vbo uint32, indexBuffer uint32, numIndices int32,
	tex uint32, program uint32) RenderComponent {
	return RenderComponent{vao, vbo, indexBuffer, numIndices, program, tex}
}

func glStr(s string) *byte {
	return gl.Str(fmt.Sprintf("%v\x00", s))
}

func glAttribLoc(program uint32, name string) uint32 {
	return uint32(gl.GetAttribLocation(program, glStr(name)))
}

func glUniformLoc(program uint32, name string) int32 {
	return gl.GetUniformLocation(program, glStr(name))
}

func (r RenderComponent) Draw(pos mgl.Vec2, VP mgl32.Mat4) {
	// global shader
	gl.Enable(gl.BLEND)
	defer gl.Disable(gl.BLEND)

	gl.UseProgram(r.program)

	positionAttrib := glAttribLoc(r.program, "position")
	texCoordAttrib := glAttribLoc(r.program, "texCoord")

	//offset uniform
	uOffsetLoc := glUniformLoc(r.program, "offset")
	offset := []float32{ float32(pos[0]), float32(pos[1]) }
	gl.Uniform2fv(uOffsetLoc, 1, &offset[0])

	//View-Projection uniform (VP)
	uProjLoc := glUniformLoc(r.program, "VP")
	vpp := [16]float32(VP)
	gl.UniformMatrix4fv(uProjLoc, 1, false, &vpp[0])

	//cylinderRadius uniform
	crLoc := glUniformLoc(r.program, "cylinderRadius")
	gl.Uniform1f(crLoc, 3)

	//cylinderHeight uniform
	chLoc := glUniformLoc(r.program, "cylinderHeight")
	gl.Uniform1f(chLoc, 0.8)

	//levelWidth uniform
	lwLoc := glUniformLoc(r.program, "levelWidth")
	gl.Uniform1f(lwLoc, float32(gLevelWidth))

	//levelHeight uniform
	lhLoc := glUniformLoc(r.program, "levelHeight")
	gl.Uniform1f(lhLoc, 3)

	//texture sampler uniform
	samplerLoc := glUniformLoc(r.program, "Sampler")
	gl.Uniform1i(samplerLoc, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, r.tex)

	// Setup array buffer stuffer
	gl.BindVertexArray(r.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.indexBuffer)

	//5*4 == 5 floats per vert (x,y,z,u,v) by 4 bytes per float
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	//offset by x,y,z floats
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	// Finally!
	gl.DrawElements(gl.TRIANGLES, r.numIndices, gl.UNSIGNED_SHORT, nil)
}

func InitGL() {
	// Initialize OpenGL, and print version number to console
	gl.Init()
	version := gl.GoStr(gl.GetString(gl.VERSION))
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

func makeProgram(vertSource, fragSource string) (uint32, error) {
	vertexShader, err := compileShader(vertSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := compileShader(fragSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(fmt.Sprintf("%v\x00", source))
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile shader program\n%v\n%v", log, source)
	}
	return shader, nil
}

func makeVertexArrayObject(vertices []float32, indices []uint16) (uint32, uint32, uint32) {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var indexBuffer uint32
	gl.GenBuffers(1, &indexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*2, gl.Ptr(indices), gl.STATIC_DRAW)

	return vao, vbo, indexBuffer
}

func createTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("Texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0,0}, draw.Src)

	var textureId uint32
	gl.GenTextures(1, &textureId)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, textureId)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	return textureId, nil
}

func checkGLerror() {
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		// string, _ := glu.ErrorString(glerr)
		panic("GL Error")
	}
}
