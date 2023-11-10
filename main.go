package main

import (
    "server/configs"
    "server/routes" //add this

    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()

    //run database
    configs.ConnectDB()

    //routes
    routes.UserRoute(e) //add this
    routes.CityRoute(e) //add this

    e.Logger.Fatal(e.Start(":6000"))
}