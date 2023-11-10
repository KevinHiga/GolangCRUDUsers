package routes

import (
    "server/controllers"
    "github.com/labstack/echo/v4"
)

func CityRoute(e *echo.Echo)  {
    e.POST("/city", controllers.CreateCity)
    e.GET("/city/:cityId", controllers.GetACity)
    e.PUT("/city/:cityId", controllers.EditACity)
    e.DELETE("/city/:cityId", controllers.DeleteACity)
    e.GET("/citys", controllers.GetAllCitys)
    e.GET("/citys/:city", controllers.GetAllCitysByKey)
}