//TODO: send commands to control the robot
//TODO: allow only one connection to control the robot, and time it out if no commands are received in a period of time

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"sync"
)

type RobotMoveCommand struct {
	Direction RobotDirection
	Speed     uint
}

// Directon to move the robot
type RobotDirection int

const (
	Forward RobotDirection = iota
	Backwards
	Left
	Right
	Stop
)

// A type of command
type ClientCommandType int

const (

	// For controlling robot movement (Forward, Backwards, Left, Right)
	MoveRobot ClientCommandType = iota
	TakeControl
)

// Represents a command coming from the client
type ClientCommand struct {
	Type ClientCommandType
	Data map[string]interface{}
}

// Represents an open connection
type WebSocketConnection struct {
	id   uint64
	conn *websocket.Conn
}

// Keeps track of all open websocket connections
type WebSocketConnectionPool struct {
	openConnections     []*WebSocketConnection
	connectionCounter   uint64
	connectionsLock     *sync.Mutex
	robotCommandChannel chan<- RobotMoveCommand

	// the connection that is currently controlling the robot
	controllingConnection *WebSocketConnection
}

func NewWebSocketConnectionPool(robotCommandChanel chan<- RobotMoveCommand) *WebSocketConnectionPool {
	return &WebSocketConnectionPool{
		openConnections:       make([]*WebSocketConnection, 1),
		connectionCounter:     0,
		connectionsLock:       &sync.Mutex{},
		robotCommandChannel:   robotCommandChanel,
		controllingConnection: nil,
	}
}

func (pool *WebSocketConnectionPool) OpenConnection(conn *websocket.Conn) *WebSocketConnection {
	pool.connectionsLock.Lock()
	defer pool.connectionsLock.Unlock()
	pool.connectionCounter += 1
	newConnection := &WebSocketConnection{
		id:   pool.connectionCounter,
		conn: conn,
	}
	pool.openConnections = append(pool.openConnections, newConnection)
	return newConnection
}

func (pool *WebSocketConnectionPool) CloseConnection(connection *WebSocketConnection) {
	connection.conn.Close()
	pool.connectionsLock.Lock()
	defer pool.connectionsLock.Unlock()
	for i, c := range pool.openConnections {
		if c.id == connection.id {
			pool.openConnections[i] = pool.openConnections[len(pool.openConnections)-1]
			pool.openConnections[len(pool.openConnections)-1] = nil
			pool.openConnections = pool.openConnections[:len(pool.openConnections)-1]
			break
		}
	}

	// if this is the controlling connection, remove it
	if pool.controllingConnection.id == connection.id {
		pool.controllingConnection = nil
	}

}

func (pool *WebSocketConnectionPool) HasControlLock(connection *WebSocketConnection) bool {
	pool.connectionsLock.Lock()
	defer pool.connectionsLock.Unlock()
	return pool.controllingConnection != nil && pool.controllingConnection.id == connection.id
}

func (pool *WebSocketConnectionPool) TakeControl(connection *WebSocketConnection) error {
	pool.connectionsLock.Lock()
	defer pool.connectionsLock.Unlock()
	if pool.controllingConnection != nil && connection.id != pool.controllingConnection.id {
		return errors.New(fmt.Sprintln("Somebody is currently controlling the robot"))
	}
	pool.controllingConnection = connection
	return nil
}

func (pool *WebSocketConnectionPool) ProcessCommands(connection *WebSocketConnection) {
	for {

		// Read command from websocket
		msg := ""
		err := websocket.Message.Receive(connection.conn, &msg)
		if err != nil {
			fmt.Println(err)
			pool.CloseConnection(connection)
			return
		}

		// parse message
		var command ClientCommand
		err = json.Unmarshal([]byte(msg), &command)
		if err != nil {
			fmt.Println(err)
			pool.CloseConnection(connection)
			return
		}

		// execute command!
		switch command.Type {
		case MoveRobot:

			// check if we are currently holding the lock on the robot
			if !pool.HasControlLock(connection) {
				continue
			}

			// process command
			var direction RobotDirection
			speed, valid := command.Data["speed"].(float64)
			if !valid {
				fmt.Println(err)
				pool.CloseConnection(connection)
				return
			}
			directionFloat, valid := command.Data["direction"].(float64)
			if !valid {
				fmt.Println(err)
				pool.CloseConnection(connection)
				return
			}
			direction = RobotDirection(directionFloat)
			pool.robotCommandChannel <- RobotMoveCommand{
				Direction: direction,
				Speed:     uint(speed),
			}
		case TakeControl:
			pool.TakeControl(connection)
		default:
			panic("Unhandled command!")
		}
	}
}
