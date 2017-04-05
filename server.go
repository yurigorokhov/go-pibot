package main

import (
	"flag"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/megapi"
	"log"

	"golang.org/x/net/websocket"
)

var robotType string

func init() {
	flag.StringVar(&robotType, "robot", "megapi", "the robot to use (megapi, mock)")
}

func initializeRobot(robotCommandChannel <-chan RobotMoveCommand) {
	flag.Parse()
	if robotType == "mock" {
		log.Println("Starting with a Mock bot")
		mockBot := NewMockBot()
		robotController := NewRobotController(mockBot, robotCommandChannel)
		robotController.Start()
	} else if robotType == "megapi" {

		// initialize motors
		megaPiAdaptor := megapi.NewAdaptor("/dev/ttyS0")
		leftMotor := megapi.NewMotorDriver(megaPiAdaptor, 4)
		rightMotor := megapi.NewMotorDriver(megaPiAdaptor, 1)

		// create robot
		megaPiRobot := NewMegaPiBot(leftMotor, rightMotor)
		robotController := NewRobotController(megaPiRobot, robotCommandChannel)

		// define work
		work := func() {
			robotController.Start()
		}

		devices := make([]gobot.Device, 2)
		devices[0] = gobot.Device(leftMotor)
		devices[1] = rightMotor
		robot := gobot.NewRobot("megaPiBot",
			[]gobot.Connection{megaPiAdaptor},
			[]gobot.Device{leftMotor, rightMotor},
			work,
		)
		robot.Start()
	} else {
		panic("Unsupported robot type")
	}
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
	e.GET("/ws", echo.WrapHandler(webSocketHandler(webSocketPool)))
	e.Logger.Fatal(e.Start(":8080"))
}

func webSocketHandler(pool *WebSocketConnectionPool) websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		newConnection := pool.OpenConnection(ws)
		pool.ProcessCommands(newConnection)
	})
}
