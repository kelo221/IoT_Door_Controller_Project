package main

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/session"
	_ "github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"strconv"
	"time"
)

type userData struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
}

func handleHTTP(lockMode *LOCK_STATUS) {

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./public", fiber.Static{
		Compress:      true,
		Index:         "login.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	store := session.New()

	app.Get("/home", func(c *fiber.Ctx) error {

		sess, err := store.Get(c)
		if err != nil {
			log.Println(err)
		}

		username := sess.Get("Username")
		isLogin := username != nil

		if isLogin {
			fmt.Println("success")
			return c.Render("index", fiber.Map{
				"User": username,
				"Mode": lockMode.GetLockStatus(),
			})
		}
		fmt.Println("not logged in")
		return c.Redirect("/")
	})

	app.Post("/auth/*", func(c *fiber.Ctx) error {

		sess, err := store.Get(c)
		if err != nil {
			log.Println(err)
		}

		salt := []byte("salt")
		p := new(userData)
		if err := c.BodyParser(p); err != nil {
			return err
		}

		passwordHash := hex.EncodeToString(HashPassword([]byte(p.Password), salt))
		expectedPasswordHash := aqlToString("FOR r IN DOOR_LOGIN FILTER r.username == \"" + p.Username + "\" RETURN r.hash")

		usernameMatch := subtle.ConstantTimeCompare([]byte(passwordHash[:]), []byte((expectedPasswordHash[:]))) == 1

		if usernameMatch {

			sess.Set("Username", p.Username)
			if err := sess.Save(); err != nil {
				panic(err)
			}

			fmt.Println("AUTH OK")
			return c.Redirect("/home")
		}
		return c.SendStatus(400)
	})

	///TODO send only when session is correct
	app.Put("/updateLock/:lockMode", func(c *fiber.Ctx) error {

		newLockMode := c.Params("lockMode")

		i, err := strconv.Atoi(newLockMode)
		if err != nil {
			// handle error
			fmt.Println(err)
			return c.SendStatus(400)
		}

		if i >= 0 && i <= 3 {
			lockMode.LockStatus = LOCK_STATUSLock(i)
			tcpSendPackage(lockMode)
			return c.SendStatus(200)
		}
		return c.SendStatus(400)

	})

	///TODO ---------------------------------------------
	app.Post("/logout", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	///TODO send only when session is correct
	app.Get("/statistic/modeChanged", func(c *fiber.Ctx) error {

		return c.Send(

			[]byte("[" +
				"{\"Who\":\"John\", \"When\":1012239724, \"Mode\":\"Soft-Lock\"}," +
				"{\"Who\":\"Max\", \"When\":1757699363, \"Mode\":\"Hard-Lock\"}," +
				"{\"Who\":\"Mona\", \"When\":1586899080 , \"Mode\":\"Open\"}," +
				"{\"Who\":\"Lisa\", \"When\":1327233733, \"Mode\":\"Hard-Lock\"}" +
				"]"),
		)
	})

	app.Get("/statistic/keycardUsed", func(c *fiber.Ctx) error {

		/// TODO
		///	ask database if hashed key is matched
		///	if matches get the username and time, store them in a database

		return c.Send(

			[]byte("[" +
				"{\"Who\":\"John\", \"When\":1012239724 }," +
				"{\"Who\":\"Max\", \"When\":1757699363}," +
				"{\"Who\":\"Mona\", \"When\":1586899080}," +
				"{\"Who\":\"Lisa\", \"When\":1327233733}" +
				"]"),
		)
	})

	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}

func tcpListenerLoop(rfidMessage *RFID_MESSAGE) {

	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go getTpcPackage(conn, rfidMessage)
		}
	}

}

func tcpSendPackage(lockMode *LOCK_STATUS) {

	//fmt.Println("check that data is correct: ", lockMode.GetLockStatus())

	data, err := proto.Marshal(lockMode)
	if err != nil {
		log.Fatal("marshall error", err)

	}
	conn, err := net.DialTimeout("tcp", "localhost:8082", time.Second*30)
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}

	_, err = conn.Write(data)
}
