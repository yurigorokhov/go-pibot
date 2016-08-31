package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/megapi"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	"golang.org/x/net/websocket"
)

func initializeRobot(robotCommandChannel <-chan RobotMoveCommand) {
	gbot := gobot.NewGobot()

	// initialize motors
	megaPiAdaptor := megapi.NewMegaPiAdaptor("megapi", "/dev/ttyS0")
	leftMotor := megapi.NewMotorDriver(megaPiAdaptor, "motor left", 1)
	rightMotor := megapi.NewMotorDriver(megaPiAdaptor, "motor right", 4)

	// create robot
	megaPiRobot := NewMegaPiBot(leftMotor, rightMotor)
	robotController := NewRobotController(megaPiRobot, robotCommandChannel)

	// define work
	work := func() {
		robotController.Start()
	}

	robot := gobot.NewRobot("megaPiBot",
		[]gobot.Connection{megaPiAdaptor},
		[]gobot.Device{leftMotor, rightMotor},
		work,
	)
	gbot.AddRobot(robot)
	gbot.Start()
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	robotCommandChannel := make(chan RobotMoveCommand)
	go initializeRobot(robotCommandChannel)

	// start websocket handler
	webSocketPool := NewWebSocketConnectionPool(robotCommandChannel)

	// routes
	e.Static("/public", "public")
	e.File("/", "public/index.html")
	e.File("/logo.png", "public/raspberry-pi-logo.png")
	e.GET("/ws", standard.WrapHandler(webSocketHandler(webSocketPool)))
	e.Run(standard.New(":8080"))
}

func webSocketHandler(pool *WebSocketConnectionPool) websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		newConnection := pool.OpenConnection(ws)
		pool.ProcessCommands(newConnection)
	})
}
