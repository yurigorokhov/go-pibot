package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	"golang.org/x/net/websocket"
	"mime/multipart"
	"net/textproto"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// start camera
	piCamera := NewPiCamFileSystem("./public/walle.jpg")
	quitChan := make(chan bool)
	piCamera.Start(quitChan)

	// start websocket handler
	webSocketPool := NewWebSocketConnectionPool()

	// routes
	e.GET("/picam.mjpeg", func(c echo.Context) error {

		// return image as multipart file transfer
		mimeWriter := multipart.NewWriter(c.Response())

		contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
		c.Response().Header().Set("Content-Type", contentType)

		// start listening for images
		source := piCamera.Subscribe()
		for {
			partHeader := make(textproto.MIMEHeader)
			partHeader.Add("Content-Type", "image/jpeg")
			partWriter, partErr := mimeWriter.CreatePart(partHeader)
			if nil != partErr {
				break
			}
			imageData := <-source
			if _, writeErr := partWriter.Write(imageData); nil != writeErr {
				break
			}
			c.Response().(http.Flusher).Flush()
		}
		return nil
	})

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
