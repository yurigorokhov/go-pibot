package main

import (
	"flag"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/megapi"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
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
		gbot := gobot.NewGobot()

		// initialize motors
		megaPiAdaptor := megapi.NewMegaPiAdaptor("megapi", "/dev/ttyS0")
		leftMotor := megapi.NewMotorDriver(megaPiAdaptor, "motor left", 4)
		rightMotor := megapi.NewMotorDriver(megaPiAdaptor, "motor right", 1)

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
	e.GET("/ws", standard.WrapHandler(webSocketHandler(webSocketPool)))
	e.Run(standard.New(":8080"))
}

func webSocketHandler(pool *WebSocketConnectionPool) websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		newConnection := pool.OpenConnection(ws)
		pool.ProcessCommands(newConnection)
	})
}
