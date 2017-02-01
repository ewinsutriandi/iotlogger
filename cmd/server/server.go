package main

import (
	"net/http"
	"time"

	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ewinsutriandi/iotlogger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const secretKey string = "ular lari lurus"

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	isvaliduser := iotlogger.Authenticate(username, password)
	if isvaliduser {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = username
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(secretKey))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	mdjwt := middleware.JWT([]byte(secretKey))

	// Login route
	e.POST("/login", login)

	// Unauthenticated route
	e.GET("/", accessible)

	//restricted route group for authenticated user
	ru := e.Group("/restricted")
	ru.Use(mdjwt)
	ru.GET("/user/:usr", getUser)
	ru.POST("/telem", putTelemReport)

	e.Logger.Fatal(e.Start(":1323"))
}

func getUser(c echo.Context) (err error) {
	username := c.Param("usr")
	ctl, err := iotlogger.NewController()
	user, err := ctl.GetUser(username)
	return c.JSON(http.StatusOK, user)
}

func putUser(c echo.Context) (err error) {
	username := c.FormValue("username")
	password := c.FormValue("password")
	ctl, err := iotlogger.NewController()
	if err != nil {
		return err
	}
	err = ctl.AddUser(username, password)
	return c.JSON(http.StatusOK, "user: "+username+" created")
}

func putTelemReport(c echo.Context) (err error) {
	ip := c.RealIP()
	devname := c.FormValue("devname")
	metric := c.FormValue("metric")
	unit := c.FormValue("unit")
	value, err := strconv.ParseFloat(c.FormValue("value"), 64)
	if err != nil {
		iotlogger.HandleError(err)
		return
	}
	telem := iotlogger.NewTelemetryReport(ip, devname, metric, unit, value)
	ctl, err := iotlogger.NewController()
	if err != nil {
		iotlogger.HandleError(err)
	}
	err = ctl.AddTelemetryReport(*telem)
	iotlogger.HandleError(err)
	return c.JSON(http.StatusCreated, telem)
}
