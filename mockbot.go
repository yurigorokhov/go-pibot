package main

var _ Robot = (*MockBot)(nil)

// MockBot is a bot that does not do anything
type MockBot struct{}

// NewMockBot creates a new MockBot
func NewMockBot() Robot {
	return &MockBot{}
}

func (robot *MockBot) GetMaxMotorSpeed() float64         { return 100 }
func (robot *MockBot) GetMinMotorSpeed() float64         { return 0 }
func (robot *MockBot) GetTurnSpeedDifferential() float64 { return 40 }
func (robot *MockBot) SetLeftMotorSpeed(speed int16)     {}
func (robot *MockBot) SetRightMotorSpeed(speed int16)    {}
func (robot *MockBot) StopAllMotors()                    {}
