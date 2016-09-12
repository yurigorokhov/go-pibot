package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

type Robot interface {
	SetLeftMotorSpeed(speed int16)
	SetRightMotorSpeed(speed int16)
	StopAllMotors()
	GetMaxMotorSpeed() float64
	GetMinMotorSpeed() float64
	GetTurnSpeedDifferential() float64
}

type RobotController struct {
	commandChannel <-chan RobotMoveCommand
	robot          Robot
}

func NewRobotController(robot Robot, commandChannel <-chan RobotMoveCommand) *RobotController {
	return &RobotController{
		commandChannel: commandChannel,
		robot:          robot,
	}
}

func (controller *RobotController) Start() {
	HandleCommand(RobotMoveCommand{
		Direction: Stop,
		Speed:     0,
	}, controller.robot)
	for {
		select {
		case command := <-controller.commandChannel:
			HandleCommand(command, controller.robot)
			log.Println(fmt.Sprintf("Received command: %+v", command))
		case <-time.After(time.Second * 2):
			HandleCommand(RobotMoveCommand{
				Direction: Stop,
				Speed:     0,
			}, controller.robot)
		}
	}
}

func HandleCommand(command RobotMoveCommand, robot Robot) {
	if command.Direction == Right || command.Direction == Left {
		command.Speed /= 4
	}
	adjustedSpeed := float64(command.Speed)*(robot.GetMaxMotorSpeed()-robot.GetMinMotorSpeed())/100 + robot.GetMinMotorSpeed()
	differential := math.Max(math.Log(float64(adjustedSpeed)*6/(robot.GetMaxMotorSpeed()-robot.GetMinMotorSpeed())), 0.3) * robot.GetTurnSpeedDifferential()
	switch command.Direction {
	case Forward:
		robot.SetLeftMotorSpeed(-int16(adjustedSpeed))
		robot.SetRightMotorSpeed(-int16(adjustedSpeed))
	case Backwards:
		robot.SetLeftMotorSpeed(int16(adjustedSpeed))
		robot.SetRightMotorSpeed(int16(adjustedSpeed))
	case Left:
		robot.SetLeftMotorSpeed(-int16(adjustedSpeed))
		robot.SetRightMotorSpeed(int16(adjustedSpeed))
	case Right:
		robot.SetLeftMotorSpeed(int16(adjustedSpeed))
		robot.SetRightMotorSpeed(-int16(adjustedSpeed))
	case Stop:
		robot.StopAllMotors()
	case ForwardRight:
		leftSpeed := math.Min(float64(adjustedSpeed)+differential, robot.GetMaxMotorSpeed())
		rightSpeed := math.Max(float64(adjustedSpeed)-differential, robot.GetMinMotorSpeed())
		robot.SetLeftMotorSpeed(-int16(leftSpeed))
		robot.SetRightMotorSpeed(-int16(rightSpeed))
	case ForwardLeft:
		leftSpeed := math.Min(float64(adjustedSpeed)-differential, robot.GetMaxMotorSpeed())
		rightSpeed := math.Max(float64(adjustedSpeed)+differential, robot.GetMinMotorSpeed())
		robot.SetLeftMotorSpeed(-int16(leftSpeed))
		robot.SetRightMotorSpeed(-int16(rightSpeed))
	case BackwardsRight:
		leftSpeed := math.Min(float64(adjustedSpeed)+differential, robot.GetMaxMotorSpeed())
		rightSpeed := math.Max(float64(adjustedSpeed)-differential, robot.GetMinMotorSpeed())
		robot.SetLeftMotorSpeed(int16(leftSpeed))
		robot.SetRightMotorSpeed(int16(rightSpeed))
	case BackwardsLeft:
		leftSpeed := math.Min(float64(adjustedSpeed)-differential, robot.GetMaxMotorSpeed())
		rightSpeed := math.Max(float64(adjustedSpeed)+differential, robot.GetMinMotorSpeed())
		robot.SetLeftMotorSpeed(int16(leftSpeed))
		robot.SetRightMotorSpeed(int16(rightSpeed))
	default:
		log.Println("ROBOT: unknown move command")
	}
}
