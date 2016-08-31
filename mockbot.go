package main

import (
	"log"
)

var _ Robot = (*MockBot)(nil)

// MockBot is a bot that does not do anything
type MockBot struct{}

// NewMockBot creates a new MockBot
func NewMockBot() Robot {
	return &MockBot{}
}

// HandleCommand implements the Robot interface
func (robot *MockBot) HandleCommand(command RobotMoveCommand) {
	if command.Direction == Right || command.Direction == Left {
		command.Speed /= 4
	}
	adjustedSpeed := command.Speed*(MAX_MOTOR_SPEED-MIN_MOTOR_SPEED)/100 + MIN_MOTOR_SPEED
	switch command.Direction {
	case Forward:
		log.Printf("ROBOT: forward at speed %+v\n", adjustedSpeed)
	case Backwards:
		log.Printf("ROBOT: backwards at speed %+v\n", adjustedSpeed)
	case Left:
		log.Printf("ROBOT: left at speed %+v\n", adjustedSpeed)
	case Right:
		log.Printf("ROBOT: right at speed %+v\n", adjustedSpeed)
	case ForwardRight:
		log.Printf("ROBOT: forward-right at speed %+v\n", adjustedSpeed)
	case ForwardLeft:
		log.Printf("ROBOT: forward-left at speed %+v\n", adjustedSpeed)
	case BackwardsRight:
		log.Printf("ROBOT: backwards-right at speed %+v\n", adjustedSpeed)
	case BackwardsLeft:
		log.Printf("ROBOT: backwards-left at speed %+v\n", adjustedSpeed)
	case Stop:
		log.Println("ROBOT: stopping")
	default:
		log.Println("ROBOT: unknown move command")
	}
}
