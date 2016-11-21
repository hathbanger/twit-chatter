package server

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/engine/standard"
)



func Run() {
// ECHO
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	 AllowOrigins: []string{"*"},
	 AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

// ROUTES
	e.GET("/ws", standard.WrapHandler(twitChatter()))
	fmt.Println("Server now running on port: 1323")
	e.Run(standard.New(":1323"))

}

