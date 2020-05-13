package pkg

import (
	"github.com/labstack/echo"
	"net/http"
)

type status struct {
	Version string `json:"version"`
}

func getStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, status{
		Version: "1.0",
	})
}

//func todoitemCreateHandler(service *Service) echo.HandlerFunc {
//	return func(c echo.Context) (err error) {
//		item := new(TodoItem)
//		if err = c.Bind(&item); err != nil {
//			return
//		}
//		err = service.CreateTodoItem(item)
//		if err != nil {
//			return
//		}
//
//		return c.JSON(http.StatusCreated, item)
//	}
//}

func getLastTodoItem(service *Service) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		item, err := service.LastTodoItem()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, item)
	}
}


//MakeHandlers define all routes to our microservice
func MakeHandlers(e *echo.Echo, service *Service) {
	e.GET("/status", getStatus).Name = "status"
	e.GET("/todo", getLastTodoItem(service)).Name = "item-create"
}

