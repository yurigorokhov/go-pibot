package main

import (
	"fmt"
	"github.com/yurigorokhov/go-megapi"
)

const MAX_MOTOR_SPEED = 150
const MIN_MOTOR_SPEED = 50

type MegaPiBot struct {
	megaPi *megapi.MegaPi
}

func NewMegaPiBot(megaPi *megapi.MegaPi) Robot {
	return &MegaPiBot{
		megaPi: megaPi,
	}
}

func (robot *MegaPiBot) HandleCommand(command RobotMoveCommand) {
	if command.Direction == Right || command.Direction == Left {
		command.Speed /= 4
	}
	adjustedSpeed := command.Speed*(MAX_MOTOR_SPEED-MIN_MOTOR_SPEED)/100 + MIN_MOTOR_SPEED
	switch command.Direction {
	case Forward:
		fmt.Printf("ROBOT: forward at speed %+v\n", adjustedSpeed)
		robot.megaPi.DcMotorRun(1, int16(adjustedSpeed))
		robot.megaPi.DcMotorRun(2, -int16(adjustedSpeed))
	case Backwards:
		fmt.Printf("ROBOT: backwards at speed %+v\n", adjustedSpeed)
		robot.megaPi.DcMotorRun(1, -int16(adjustedSpeed))
		robot.megaPi.DcMotorRun(2, +int16(adjustedSpeed))
	case Left:
		fmt.Printf("ROBOT: left at speed %+v\n", adjustedSpeed)
		robot.megaPi.DcMotorRun(1, int16(adjustedSpeed)+1)
		robot.megaPi.DcMotorRun(2, int16(adjustedSpeed))
	case Right:
		fmt.Printf("ROBOT: right at speed %+v\n", adjustedSpeed)
		robot.megaPi.DcMotorRun(1, -int16(adjustedSpeed)+1)
		robot.megaPi.DcMotorRun(2, -int16(adjustedSpeed))
	case Stop:
		fmt.Println("ROBOT: stopping")
		robot.megaPi.DcMotorStop(1)
		robot.megaPi.DcMotorStop(2)
	default:
		fmt.Println("PiBot: unknown move command")
	}
}
