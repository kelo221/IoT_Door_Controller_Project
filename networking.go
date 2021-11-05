package main

import (
	"crypto/subtle"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	_ "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

type userData struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
}

func handleHTTP() {

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Test App v1.0.1",
	})

	app.Static("/", "./public", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "login.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	app.Use(limiter.New(), encryptcookie.New(encryptcookie.Config{
		Key: "secret-thirty-2-character-string",
	}))

	app.Post("/auth/*", func(c *fiber.Ctx) error {
		salt := []byte("salt")
		p := new(userData)
		if err := c.BodyParser(p); err != nil {
			return err
		}

		passwordHash := hex.EncodeToString(HashPassword([]byte(p.Password), salt))
		expectedPasswordHash := aqlToString("FOR r IN DOOR_LOGIN FILTER r.username == \"" + p.Username + "\" RETURN r.hash")

		usernameMatch := subtle.ConstantTimeCompare([]byte(passwordHash[:]), []byte((expectedPasswordHash[:]))) == 1

		if usernameMatch {
			return c.SendStatus(200)
		}
		return c.SendStatus(400)
	})

	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
