package main

import (
	"gobot.io/x/gobot/platforms/megapi"
)

var _ Robot = (*MegaPiBot)(nil)

type MegaPiBot struct {
	leftMotor             *megapi.MotorDriver
	rightMotor            *megapi.MotorDriver
	maxMotorSpeed         float64
	minMotorSpeed         float64
	turnSpeedDifferential float64
}

func NewMegaPiBot(leftMotor *megapi.MotorDriver, rightMotor *megapi.MotorDriver) Robot {
	return &MegaPiBot{
		leftMotor:             leftMotor,
		rightMotor:            rightMotor,
		maxMotorSpeed:         250,
		minMotorSpeed:         50,
		turnSpeedDifferential: 70,
	}
}

func (robot *MegaPiBot) GetMaxMotorSpeed() float64         { return robot.maxMotorSpeed }
func (robot *MegaPiBot) GetMinMotorSpeed() float64         { return robot.minMotorSpeed }
func (robot *MegaPiBot) GetTurnSpeedDifferential() float64 { return robot.turnSpeedDifferential }

func (robot *MegaPiBot) SetLeftMotorSpeed(speed int16) {
	robot.leftMotor.Speed(speed)
}

func (robot *MegaPiBot) SetRightMotorSpeed(speed int16) {
	robot.rightMotor.Speed(speed)
}

func (robot *MegaPiBot) StopAllMotors() {
	robot.leftMotor.Speed(0)
	robot.rightMotor.Speed(0)
}
