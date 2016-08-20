package main

import (
	"time"
)

type Robot interface {
	HandleCommand(command RobotMoveCommand)
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
	controller.robot.HandleCommand(RobotMoveCommand{
		Direction: Stop,
		Speed:     0,
	})
	go func() {
		for {
			select {
			case command := <-controller.commandChannel:
				controller.robot.HandleCommand(command)
			case <-time.After(time.Second * 2):
				controller.robot.HandleCommand(RobotMoveCommand{
					Direction: Stop,
					Speed:     0,
				})

			}
		}
	}()
}
