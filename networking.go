package main

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"strconv"
	"time"

	jwtware "github.com/gofiber/jwt/v3"
)

// @title IoT Door Controller Project
// @version 1.0
// @description User interface for handling the locking of the door.

// @license.name AGPL V3
// @license.url https://www.gnu.org/licenses/agpl-3.0.en.html

// @host localhost:8080
// @BasePath /
// @schemes http

func handleHTTP() {

	app := fiber.New(fiber.Config{})

	app.Static("/", "./public", fiber.Static{
		Compress:      false,
		Index:         "login.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	app.Use(cors.New())

	app.Post("/login", loginHandler)

	app.Get("/swagger/*", swagger.Handler) // default

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	/*		---		RESTRICTED APIS		---		*/
	app.Get("/getInitData", userDataHandler)
	app.Put("/manualOpen", forceOpen)
	app.Put("/updateLock/:lockMode", updateLock)
	app.Get("/statistics/modeChanged", lockHistoryHandler)
	app.Get("/statistics/keycardUsed", keycardHistoryHandler)

	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}

// keycardHistoryHandler godoc
// @Summary Returns the date when the RFID card was read
// @Description Queries the ArangoDB for RFID history and returns results as JSON
// @Tags root
// @Accept */*
// @Produce json
// @Success 200
// @Router /statistics/keycardUsed [get]
func keycardHistoryHandler(c *fiber.Ctx) error {
	//fmt.Println("Door history requested")
	return c.JSON(aqlJSON("FOR x IN DOOR_HISTORY RETURN x"))
}

// lockHistoryHandler godoc
// @Summary Returns the date when the lock mode was changed
// @Description Queries the ArangoDB for lock history and returns results as JSON
// @Tags root
// @Accept */*
// @Produce json
// @Success 200
// @Router /statistics/modeChanged [get]
func lockHistoryHandler(c *fiber.Ctx) error {
	//fmt.Println("Lock history requested")
	return c.JSON(aqlJSON("FOR x IN LOCK_HISTORY RETURN x"))
}

// updateLock godoc
// @Summary Set the lock mode
// @Description User sets the wanted lock mode from the interface
// @Tags root
// @Param        ?   query     int  true  "Set lock Mode"
// @Accept */*
// @Success 200
// @Router /updateLock/ [put]
func updateLock(c *fiber.Ctx) error {

	fmt.Println("new lockmode received")

	newLockMode := c.Params("lockMode")

	i, err := strconv.Atoi(newLockMode)
	if err != nil {
		// handle error
		fmt.Println(err)
		return c.SendStatus(400)
	}
	fmt.Println(newLockMode)

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	fmt.Println(name)

	if i >= 0 && i <= 3 {
		tcpPacketOut.LockStatus = LOCK_STATUSLock(i)
		handleMQTTOut()

		aqlNoReturn(
			fmt.Sprintf("INSERT {name: \"%s\" , time: DATE_NOW(),mode: \"%s\" } INTO LOCK_HISTORY", name, tcpPacketOut.GetLockStatus().String()))

		fmt.Println("added to lock history DB")

		return c.SendStatus(200)
	}
	return c.SendStatus(400)

}

// forceOpen godoc
// @Summary Commands the embedded device to open the door
// @Description When the lock state of the device is set to hard, the door can only be opened manually from the interface
// @Tags root
// @Accept */*
// @Success 200
// @Router /manualOpen [put]
func forceOpen(c *fiber.Ctx) error {
	tcpPacketOut.DoorRequest = LOCK_STATUS_APPROVED
	handleMQTTOut()
	fmt.Println("door manually opened")
	tcpPacketOut.DoorRequest = LOCK_STATUS_NO_REQUEST
	return c.SendStatus(200)
}

// userDataHandler godoc
// @Summary Returns the current users name and the lock status
// @Description Gets the currents users claim the name that is associated with it and check the status of the Protobuffer package
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} string
// @Failure	401
// @Router /getInitData [get]
func userDataHandler(c *fiber.Ctx) error {

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)

	s := "{\"name\":\"" + name + "\", \"mode\":\"" + tcpPacketOut.GetLockStatus().String() + "\"}"

	return c.SendString(s)
}

// loginHandler godoc
// @Summary Handles the authentication of the users.
// @Description If the user is found in the database the server returns a JWT token, which is used to access the other APIs
// @Tags root
// @Accept */*
// @Produce json
// @Param       username  query       string  true  "Username"
// @Param       password  query       string  true  "Password"
// @Success 200 {object} map[string]interface{}
// @Failure	401
// @Router /login [post]
func loginHandler(c *fiber.Ctx) error {
	user := c.FormValue("username")
	pass := c.FormValue("password")

	// Create the Claims
	///
	claims := jwt.MapClaims{
		"name":  aqlToString("FOR r IN DOOR_LOGIN FILTER r.username == \"" + user + "\" RETURN r.name"),
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	salt := []byte("salt")
	passwordHash := hex.EncodeToString(HashPassword([]byte(pass), salt))
	expectedPasswordHash := aqlToString("FOR r IN DOOR_LOGIN FILTER r.username == \"" + user + "\" RETURN r.hash")

	usernameMatch := subtle.ConstantTimeCompare([]byte(passwordHash[:]), []byte((expectedPasswordHash[:]))) == 1

	if usernameMatch {

		fmt.Println("AUTH OK")
		return c.JSON(fiber.Map{"token": t})
	}
	return c.SendStatus(fiber.StatusUnauthorized)
}
