package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/env"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/firebase"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/interfaces"
)

func deferMainFunc() {
	println("Closing...")
	firebase.DeferFunction()
	// mongodb.DeferFunction()
	// redis.DeferFunction()
}
func main() {
	defer deferMainFunc()
	app := fiber.New(fiber.Config{
		ServerHeader: "Bullion Server V1.0.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			mappedError, ok := err.(*interfaces.RequestError)
			if !ok {
				println(err.Error())
				return c.Status(500).JSON(interfaces.RequestError{
					Code:    interfaces.ERROR_INTERNAL_SERVER,
					Message: "Some Internal Error",
					Name:    "Global Error Handler Function",
				})
			}
			return c.Status(mappedError.StatusCode).JSON(mappedError)
		},
	})
	firebaseDb := firebase.GetFirebaseFirestoreDatabase()
	app.Use("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if len(id) == 0 || strings.Contains(id, ".") {
			return c.Status(400).JSON(interfaces.RequestError{
				Code:    interfaces.ERROR_INVALID_INPUT,
				Message: "Please Pass Valid short URL",
				Name:    "INVALID SHORT URL",
			})
		}
		res, err := firebaseDb.GetShortUrl(id)
		if err != nil || res == nil {
			return c.Status(400).JSON(interfaces.RequestError{
				Code:    interfaces.ERROR_INVALID_INPUT,
				Message: "Please Pass Valid short URL",
				Name:    "INVALID SHORT URL",
			})
		}
		url := res["url"]
		if url == nil {
			return c.Status(400).JSON(interfaces.RequestError{
				Code:    interfaces.ERROR_INVALID_INPUT,
				Message: "Redirect URL Not Found",
				Name:    "INVALID SHORT URL",
			})
		}
		if url, ok := url.(string); !ok {
			return c.Status(400).JSON(interfaces.RequestError{
				Code:    interfaces.ERROR_INVALID_INPUT,
				Message: "Redirect URL Not Found Fount JSON",
				Name:    "INVALID SHORT URL",
			})
		} else {
			return c.Redirect(url)
		}
	})
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Sorry can't find that!")
	})
	hostAndPort := ""
	if env.Env.APP_ENV == env.APP_ENV_LOCAL || env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		hostAndPort = "127.0.0.1"
	}

	hostAndPort = hostAndPort + ":" + strconv.Itoa(env.Env.PORT)
	app.Listen(hostAndPort)
	// app.Use(logger.New())
	// app.Use(middleware.TokenDecrypter)
	// fmt.Println(Hello("LinkShortner"))
}
