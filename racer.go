package main

import (
	"log"
	"math"
	"time"

	"github.com/DanielFallon/gobot/platforms/dji/tello"
	"github.com/simulatedsimian/joystick"
	"gobot.io/x/gobot"
)

func main() {
	// port 8889 is the one on which it recieves instructions, 8888 is where it replies
	drone0 := tello.NewEDUDriver("10.0.0.196:8889", "8888")
	work := func() {
		drone0.SendCommand("command")
		readJoystick(0, drone0)
	}
	racer0 := gobot.NewRobot("Racer1",
		[]gobot.Connection{},
		[]gobot.Device{drone0},
		work,
	)
	racer0.Start()

}

//Axis List
const (
	axLeftX  = 0
	axLeftY  = 1
	axL1     = 2
	axRightY = 3
	axRightX = 4
	axL2     = 5
	axR1     = 6
	axR2     = 7
	deadzone = 2000
)

//Button List
const (
	btnA  = 0
	btnB  = 1
	btnX  = 2
	btnY  = 3
	btnL1 = 4
	btnL2 = 5
	btnR1 = 6
	btnR2 = 7
)

func readJoystick(jsid int, tello *tello.Driver) {
	js, err := joystick.Open(jsid)
	if err != nil {
		log.Printf("Error  reading joystick: %v\n", err)
	}

	var jsState, prevState joystick.State
	var x, y, z, rotation float32

	for {
		jsState, err = js.Read()
		if err != nil {
			log.Printf("Error reading joystick: %v\n", err)
		}

		/// BUTTON HANDLERS
		// if A is pressed and WASN'T pressed in previous state
		if jsState.Buttons&(1<<btnA) != 0 && prevState.Buttons&(1<<btnA) == 0 {
			log.Println("A pressed")
			tello.TakeOff()
		}
		if jsState.Buttons&(1<<btnX) != 0 && prevState.Buttons&(1<<btnX) == 0 {
			log.Println("X pressed")
			tello.Land()
		}
		if jsState.Buttons&(1<<btnB) != 0 && prevState.Buttons&(1<<btnB) == 0 {
			log.Println("B pressed")
			tello.Land()
		}
		if jsState.Buttons&(1<<btnY) != 0 && prevState.Buttons&(1<<btnY) == 0 {
			log.Println("Y pressed")
			tello.Land()
		}
		//log.Println(jsState.Buttons)
		// END BUTTON HANDLER
		// AXIS HANDLER

		//Left Stick Axis (should control rotation)
		rotation = getAxisValue(float64(jsState.AxisData[axLeftX]))
		z = getAxisValue(float64(jsState.AxisData[axLeftY])) * -1 //for whatever reason this axis is flipped

		x = getAxisValue(float64(jsState.AxisData[axRightX])) * -1 //for whatever reason this axis is flipped
		y = getAxisValue(float64(jsState.AxisData[axRightY]))

		tello.SetVector(x, y, z, rotation)
		log.Println(tello.Vector())
		// Set State & Poll
		prevState = jsState
		time.Sleep(time.Millisecond)
	}
}

func getAxisValue(rawValue float64) (i float32) {
	//log.Println("Raw Value: ", rawValue)
	if math.Abs(rawValue) < 2000 {
		//log.Println("Deadzone")
		return 0
	}

	var val = rawValue / 32768 //-1 to 1
	return float32(val)
}
