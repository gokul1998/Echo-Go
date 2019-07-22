package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Cat struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}
type Dog struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

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

//faster method
func addCats(c echo.Context) error {
	cat := Cat{}
	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("failed to read request body %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("failed to unmarshall the json  %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("This is your cat %#v", cat)
	return c.String(http.StatusOK, "we got your cat")
}

//slower method
func addDogs(c echo.Context) error {
	dog := Dog{}
	err := c.Bind(&dog)
	if err != nil {
		log.Printf("failed processing Dogs request : %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("This is your Dog %#v", dog)
	return c.String(http.StatusOK, "we got your dog")
}
func main() {
	fmt.Println("hellllo world")
	e := echo.New()
	e.GET("/", home)
	e.GET("/cats/:type", getCats)
	e.POST("/cats", addCats)
	e.POST("/dogs", addDogs)
	e.Start(":8000")
}
