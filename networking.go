package main

import (
	"bytes"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/session"
	_ "github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"strconv"
	"time"
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
	app.Use(cors.New())

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

		p := new(userData)
		if err := c.BodyParser(p); err != nil {
			return err
		}

		salt := []byte("salt")
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
/*                     tcpSendPackage(lockMode) */

                    p := new(userData)
                    if err := c.BodyParser(p); err != nil {
                        return err
                    }

                    aqlNoReturn(
                        "INSERT {" +
                            "   name: " +
                            p.Username +
                            "," +
                            "    time: DATE_NOW()" +
                            "  } INTO LOCK_HISTORY OPTIONS { ignoreErrors: true }")

                    fmt.Println("added to lock history DB")

                    return c.SendStatus(200)
                }
                return c.SendStatus(400)


        })

	///TODO ---------------------------------------------
	app.Post("/logout", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
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
						log.Fatal("marshall error", err)

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
		log.Fatal("marshall error", err)

	}
	conn, err := net.DialTimeout("tcp", embeddedAddress+":"+embeddedPort, time.Second*30)
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}

	_, err = conn.Write(data)
}
