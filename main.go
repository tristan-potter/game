package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/DrTeePot/game/entity"
	"github.com/DrTeePot/game/light"
	"github.com/DrTeePot/game/loader"
	"github.com/DrTeePot/game/maths"
	"github.com/DrTeePot/game/model"
	"github.com/DrTeePot/game/render"
	"github.com/DrTeePot/game/shaders"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 1260
	windowHeight = 720
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to
	//	be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// TODO manage all this in a display thing
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "GAME", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	vertexShader := "shaders/vertexShader.glsl"
	fragmentShader := "shaders/fragmentShader.glsl"

	shader, err := shaders.NewBasicShader(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	render.Initialize(shader)

	rawModel, err := loader.LoadObjModel("assets/dragon.obj")
	if err != nil {
		fmt.Println("problem loading model")
		panic(err)
	}
	textureID, err := loader.LoadTexture("assets/blank.png")
	if err != nil {
		panic(err)
	}
	texture := model.NewTexture(textureID)
	texture.SetShine(10)
	texture.SetReflectivity(1)
	model := model.NewTexturedModel(rawModel, texture)

	var camera entity.Camera // new 0'd camera

	entity := entity.Entity{
		Model:    model,
		Position: mgl32.Vec3{0, -5, -20},
		RotX:     0,
		RotY:     0,
		RotZ:     0,
		Scale:    1,
	}
	coolLight := light.Create(
		mgl32.Vec3{5, 5, -15},
		mgl32.Vec3{1, 1, 1},
	)

	// Configure global settings
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0.11, 0.545, 0.765, 0.0)

	// this is disgusting. TODO make a keyboard package
	_ = window.SetKeyCallback(
		func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			if key == glfw.KeyW && action == glfw.Press {
				camera.Position = camera.Position.Sub(mgl32.Vec3{0, 0, 0.2})
			}
			if key == glfw.KeyS && action == glfw.Press {
				camera.Position = camera.Position.Add(mgl32.Vec3{0, 0, 0.2})
			}
			if key == glfw.KeyD && action == glfw.Press {
				camera.Position = camera.Position.Add(mgl32.Vec3{0.2, 0, 0})
			}
			if key == glfw.KeyA && action == glfw.Press {
				camera.Position = camera.Position.Sub(mgl32.Vec3{0.2, 0, 0})
			}
		})

	for !window.ShouldClose() {
		// entity.IncreasePosition(0.0, 0, -0.002)
		entity.IncreaseRotation(0, 1, 0)
		render.Prepare()
		shader.Start()
		shader.LoadLightPosition(coolLight.Position())
		shader.LoadLightColour(coolLight.Colour())

		// move the world
		// could put this somewhere else?
		viewMatrix := maths.CreateViewMatrix(camera)
		shader.LoadViewMatrix(viewMatrix)

		render.Render(entity, shader)

		shader.Stop()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

	shader.Delete()
	loader.CleanUp()
}
