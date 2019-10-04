package main

import (
	"fmt"
	"github.com/DanielFallon/gobot/platforms/dji/tello"
	"github.com/simulatedsimian/joystick"
	"log"
	"runtime"
	"time"
)

var (
	jsConfig joystickConfig
	err      error
)

//Axis List
const (
	axLeftX = iota
	axLeftY
	axRightX
	axRightY
	axL1
	axL2
	axR1
	axR2
)

//Button List
const (
	btnX = iota
	btnCircle
	btnTriangle
	btnSquare
	btnL1
	btnL2
	btnL3
	btnR1
	btnR2
	btnR3
	btnUnknown
)

const deadZone = 2000

type joystickConfig struct {
	axes    []int
	buttons []uint
}

var dualShock4Config = joystickConfig{
	axes: []int{
		axLeftX: 0, axLeftY: 1, axRightX: 3, axRightY: 4,
	},
	buttons: []uint{
		btnX: 0, btnCircle: 1, btnTriangle: 2, btnSquare: 3, btnL1: 4,
		btnL2: 6, btnR1: 5, btnR2: 7,
	},
}

var dualShock4ConfigWin = joystickConfig{
	axes: []int{
		axLeftX: 0, axLeftY: 1, axRightX: 2, axRightY: 3,
	},
	buttons: []uint{
		btnX: 1, btnCircle: 2, btnTriangle: 3, btnSquare: 0, btnL1: 4,
		btnL2: 6, btnR1: 5, btnR2: 7,
	},
}

// hotas mapping seems the same on windows and linux
var tflightHotasXConfig = joystickConfig{
	axes: []int{
		axLeftX: 4, axLeftY: 2, axRightX: 0, axRightY: 1,
	},
	buttons: []uint{
		btnR1: 0, btnL1: 1, btnR3: 2, btnL3: 3, btnSquare: 4, btnX: 5,
		btnCircle: 6, btnTriangle: 7, btnR2: 8, btnL2: 9,
	},
}

func printJoystickHelp() {
	fmt.Print(
		`TelloTerm Joystick Control Mapping

Right Stick  Forward/Backward/Left/Right
Left Stick   Up/Down/Turn
Triangle     Takeoff
X            Land
Circle       
Square       Take Photo
L1           Bounce (on/off)
L2           Palm Land
`)
}

func listJoysticks() {
	for jsid := 0; jsid < 10; jsid++ {
		js, err := joystick.Open(jsid)
		if err != nil {
			if jsid == 0 {
				fmt.Println("No joysticks detected")
			}
			return
		}
		fmt.Printf("Joystick ID: %d: Name: %s, Axes: %d, Buttons: %d\n", jsid, js.Name(), js.AxisCount(), js.ButtonCount())
		js.Close()
	}
}

func setupJoystickConfig() bool {
	switch runtime.GOOS {
	case "windows":
		jsConfig = dualShock4ConfigWin
	default:
		jsConfig = dualShock4Config
	}
	return true
}

func intAbs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

func droneCommand(jsid int, tello *tello.Driver) {
	js, err := joystick.Open(jsid)
	if err != nil {
		log.Printf("Error reading joystick: %v\n", err)
	}

	//run this function for a combination of a joystick and tello
	var jsState, prevState joystick.State
	var x, y, z, rotation float32
	var s, e float32
	s = 32767 //int16.MinValue
	e = 65534 //int.MaxValue*2
	for {
		jsState, err = js.Read()

		if err != nil {
			log.Printf("Error reading joystick: %v\n", err)
		}
		// rescaling int16 (-32767 to +32767) to float32 0 to 1,
		// value + 32767 / 65534

		// Axis Updates
		//left X == rotation cc or ccw
		var axLeftX = float32(jsState.AxisData[jsConfig.axes[axLeftX]])
		rotation = (axLeftX - s) / e
		if rotation > 1 || axLeftX == 32768 {
			rotation = 1
		}
		//left y == set y
		var axLeftY = float32(jsState.AxisData[jsConfig.axes[axLeftY]])
		y = (axLeftY - s) / e
		if y > 1 || axLeftY == 32768 {
			y = 1
		}
		//right x == set X
		var axRightX = float32(jsState.AxisData[jsConfig.axes[axRightX]])
		x = (axRightX - s) / e
		if x > 1 || axRightX == 32768 {
			x = 1
		}
		//right y == set z
		var axRightY = float32(jsState.AxisData[jsConfig.axes[axRightY]])
		z = (axRightY - s) / e
		if z > 1 || axRightY == 32768 {
			z = 1
		}
		tello.SetVector(x, y, z, rotation)
		if jsState.Buttons&(1<<jsConfig.buttons[btnTriangle]) != 0 && prevState.Buttons&(1<<jsConfig.buttons[btnTriangle]) == 0 {
			log.Println("Y pressed")
			tello.TakeOff()
		}
		if jsState.Buttons&(1<<jsConfig.buttons[btnX]) != 0 && prevState.Buttons&(1<<jsConfig.buttons[btnX]) == 0 {
			log.Println("X pressed")
			tello.Land()
		}
		time.Sleep(time.Millisecond * 50)
	}
}
