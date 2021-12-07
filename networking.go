package main

import (
	"bytes"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/template/html"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type userData struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
}

func handleHTTP(lockMode *Door_Request) {

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./public", fiber.Static{
		Compress:      false,
		Index:         "login.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	app.Use(cors.New())

	app.Post("/login", func(c *fiber.Ctx) error {

		user := c.FormValue("username")
		pass := c.FormValue("password")

		// Create the Claims
		claims := jwt.MapClaims{
			"name":  "John Doe",
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
	})

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	app.Get("/getInitData", func(c *fiber.Ctx) error {

		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)

		s := "{\"name\":\"" + name + "\", \"mode\":\"" + lockMode.GetLockStatus().String() + "\"}"

		return c.SendString(s)

	})

	/*		---		RESTRICTED		---		*/

	app.Put("/updateLock/:lockMode", func(c *fiber.Ctx) error {

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
			lockMode.LockStatus = LOCK_STATUSLock(i)
			/*                     tcpSendPackage(lockMode) */

			aqlNoReturn(
				fmt.Sprintf("INSERT {name: \"%s\" , time: DATE_NOW(),mode: \"%s\" } INTO LOCK_HISTORY", name, lockMode.GetLockStatus().String()))

			fmt.Println("added to lock history DB")

			return c.SendStatus(200)
		}
		return c.SendStatus(400)

	})

	app.Get("/statistics/modeChanged", func(c *fiber.Ctx) error {

		fmt.Println("Lock history requested")

		return c.JSON(aqlJSON("FOR x IN LOCK_HISTORY RETURN x"))

	})

	app.Get("/statistics/keycardUsed", func(c *fiber.Ctx) error {

		fmt.Println("Door history requested")

		return c.JSON(aqlJSON("FOR x IN DOOR_HISTORY RETURN x"))
	})

	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}

func tcpListenerLoop(lockMode *Door_Request) {

	listen, err := net.Listen("tcp", ":"+tcpListenPort)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
		} else {

			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {

				}
			}(conn)

			result := bytes.NewBuffer(nil)
			var buf [1024]byte
			for {
				n, err := conn.Read(buf[0:])
				result.Write(buf[0:n])
				if err != nil {
					if err == io.EOF {
						continue
					} else {
						//	fmt.Println(err)
						break
					}
				} else {
					newMessage := RFID_MESSAGE{}
					err = proto.Unmarshal(result.Bytes(), &newMessage)
					if err != nil {
						panic(err)
					}
					fmt.Println(newMessage.GetRFID_CODE())

					salt := []byte("salt")

					rfidHash := hex.EncodeToString(HashPassword([]byte(newMessage.GetRFID_CODE()), salt))
					expectedHash := aqlToString("FOR doc IN DOOR_RFID  FILTER doc.HASHED_RFID == \"" + rfidHash + "\"  RETURN doc.HASHED_RFID")

					hashFound := subtle.ConstantTimeCompare([]byte(rfidHash[:]), []byte((expectedHash[:]))) == 1

					if hashFound {
						fmt.Println("match")

						aqlNoReturn("FOR doc IN DOOR_RFID" +
							" FILTER doc.HASHED_RFID == \"" + rfidHash + "\"" +
							"   INSERT {" +
							"   name: doc.RFID_OWNER," +
							"    time: DATE_NOW()" +
							"  } INTO DOOR_HISTORY OPTIONS { ignoreErrors: true }")

						lockMode.DoorRequest = LOCK_STATUS_APPROVED
					} else {
						lockMode.DoorRequest = LOCK_STATUS_DISAPPROVED
					}

					data, err := proto.Marshal(lockMode)
					if err != nil {
						fmt.Println("marshall error", err)

					}

					conn, err := net.DialTimeout("tcp", embeddedAddress+":"+embeddedPort, time.Second*30)
					if err != nil {
						fmt.Printf("connect failed, err : %v\n", err.Error())
						return
					}

					_, err = conn.Write(data)

					lockMode.DoorRequest = LOCK_STATUS_NO_REQUEST
				}
				result.Reset()
			}
		}
	}

}

func tcpSendPackage(lockMode *Door_Request) {

	data, err := proto.Marshal(lockMode)
	if err != nil {
		fmt.Println("marshall error", err)

	}
	conn, err := net.DialTimeout("tcp", embeddedAddress+":"+embeddedPort, time.Second*30)
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}

	_, err = conn.Write(data)
}
