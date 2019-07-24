package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Cat struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}
type Dog struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}
type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func home(c echo.Context) error {
	return c.String(http.StatusOK, "home page")
}
func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "you are in admin page")
}
func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Custom Header")
		c.Response().Header().Set("Location", "BLR")
		return next(c)
	}
}
func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	if username == "gokul" && password == "12345" {
		cookie := &http.Cookie{}
		cookie.Name = "SessionID"
		cookie.Value = "some_value"
		cookie.Expires = time.Now().Add(10 * time.Minute)
		c.SetCookie(cookie)
		token, err := createJwtToken()
		if err != nil {
			log.Println("error in creating JWT tokens")
			return c.String(http.StatusInternalServerError, "something went wrong")
		}
		return c.JSON(http.StatusOK, map[string]string{
			"message": "you are logged in",
			"token":   token,
		})
	}
	return c.String(http.StatusUnauthorized, "enter the correct credentials")
}
func createJwtToken() (string, error) {
	claims := JwtClaims{
		"jack",
		jwt.StandardClaims{
			Id:        "user_id",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := rawToken.SignedString([]byte("secretstring"))
	if err != nil {
		return "", nil
	}
	return token, nil
}
func mainJwt(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claim := token.Claims.(jwt.MapClaims)
	log.Println("username : ", claim["name"], " expires in : ", claim["exp"], " user_id : ", claim["jti"])
	return c.String(http.StatusOK, "you are on secret Jwt page")
}
func main() {
	fmt.Println("hellllo world")
	e := echo.New()
	g := e.Group("/admin")
	jwtgroup := e.Group("/jwt")
	jwtgroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("secretstring"),
	}))
	//cookiegroup := e.Group("/cookie")
	g.Use(serverHeader)
	//g.Use(middleware.Logger()) ->simple logger
	//this is a logger with customized output
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host} ${latency_human}` + "\n",
	}))
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "gokul" && password == "12345" {
			return true, nil
		}
		return false, nil
	}))
	jwtgroup.GET("/main", mainJwt)
	e.GET("/login", login)
	g.GET("/main", mainAdmin)
	e.GET("/", home)
	e.Start(":8000")
}
