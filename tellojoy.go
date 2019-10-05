package main

import (
	"github.com/DanielFallon/gobot/platforms/dji/tello"
	"gobot.io/x/gobot"
)

func main() {
	// port 8889 is the one on which it recieves instructions, 8888 is where it replies
	setupJoystickConfig()
	drone0 := tello.NewEDUDriver("10.0.0.196:8889", "8888")
	work := func() {
		drone0.SendCommand("command")
		//drone0.SendCommand("takeoff")
		drone0.TakeOff()
		//droneCommand(0, drone0)
	}
	racer0 := gobot.NewRobot("Racer1",
		[]gobot.Connection{},
		[]gobot.Device{drone0},
		work,
	)

	racer0.Start()
}

//func droneCommand(drone, cmd){}
