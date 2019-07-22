package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func home(c echo.Context) error {
	return c.String(http.StatusOK, "home page")
}
func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catAge := c.QueryParam("age")
	dataType := c.Param("type")
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is %s \n your cat's age is %s", catName, catAge))
	}
	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"age":  catAge,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "please enter either string or json format",
	})
}
func main() {
	fmt.Println("hellllo world")
	e := echo.New()
	e.GET("/", home)
	e.GET("/cats/:type", getCats)
	e.Start(":8000")
}
