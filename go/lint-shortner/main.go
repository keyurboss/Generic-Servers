package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/apis"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/env"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/firebase"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/interfaces"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/middleware"
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

	// API STARTS HERE

	app.Get("/:id", apis.GetShortUrlAfterRedirect)
	app.Post("/get_token", apis.GenerateToken)
	g := app.Group("/", middleware.TokenDecrypter, middleware.AllowOnlyValidTokenMiddleWare)
	g.Post("/short_url", apis.SignleShortning)
	g.Post("/bulk_shortning", apis.SignleShortning)

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
