package main

import (
	"fmt"
	"github.com/hybridgroup/gobot/platforms/megapi"
)

const MAX_MOTOR_SPEED = 250
const MIN_MOTOR_SPEED = 50

type MegaPiBot struct {
	leftMotor  *megapi.MotorDriver
	rightMotor *megapi.MotorDriver
}

func NewMegaPiBot(leftMotor *megapi.MotorDriver, rightMotor *megapi.MotorDriver) Robot {
	return &MegaPiBot{
		leftMotor:  leftMotor,
		rightMotor: rightMotor,
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
		robot.leftMotor.Speed(int16(adjustedSpeed))
		robot.rightMotor.Speed(-int16(adjustedSpeed))
	case Backwards:
		fmt.Printf("ROBOT: backwards at speed %+v\n", adjustedSpeed)
		robot.leftMotor.Speed(-int16(adjustedSpeed))
		robot.rightMotor.Speed(int16(adjustedSpeed))
	case Left:
		fmt.Printf("ROBOT: left at speed %+v\n", adjustedSpeed)
		robot.leftMotor.Speed(-int16(adjustedSpeed))
		robot.rightMotor.Speed(-int16(adjustedSpeed))
	case Right:
		fmt.Printf("ROBOT: right at speed %+v\n", adjustedSpeed)
		robot.leftMotor.Speed(+int16(adjustedSpeed))
		robot.rightMotor.Speed(+int16(adjustedSpeed))
	case Stop:
		fmt.Println("ROBOT: stopping")
		robot.leftMotor.Speed(1)
		robot.leftMotor.Speed(0)
		robot.rightMotor.Speed(1)
		robot.rightMotor.Speed(0)
	default:
		fmt.Println("PiBot: unknown move command")
	}
}
