package main

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/keyurboss/Generic-Servers/go/jwelly-server/apis"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/interfaces"
)

func main() {
	app := fiber.New(fiber.Config{
		ServerHeader: "JWELLY MOCK SERVER V1.0.0",
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
	app.Get("/softwarekey/demo.php", apis.DemoApi)
	app.Get("/softwarekey/democheckpin.php", apis.DemoCheckPinApi)
	app.Get("/softwarekey/erpupg.php", apis.ErpupgradeApi)
	// app.Post("/softwarekey/reg_client_device.php", apis.DemoCheckPinApi1)
	// app.Get("/softwarekey/reg_client_device.php", apis.DemoCheckPinApi)
	// app.Post("/softwarekey/register_device.php", apis.DemoCheckPinApi1)
	// app.Get("/softwarekey/register_device.php", apis.DemoCheckPinApi1)
	// SOFTWARE SETTINGS API
	app.Post("/softwarekey/getenqoption.php", apis.GetenqoptionApi)
	// app.Get("/:id", apis.GetShortUrlAfterRedirect)
	// app.Post("/get_token", apis.GenerateToken)
	// g := app.Group("/", middleware.TokenDecrypter, middleware.AllowOnlyValidTokenMiddleWare)
	// g.Post("/short_url", apis.SignleShortning)
	// g.Post("/bulk_shortning", apis.SignleShortning)
	app.Use(func(c *fiber.Ctx) error {
		fmt.Println("*****----new req----*****")
		fmt.Printf("Path: %v\n", c.Request().URI().String())
		fmt.Printf("Path: %s\n", c.Request().Header.Method())
		fmt.Printf("Headers: %v\n", c.GetReqHeaders())
		fmt.Printf("URL PARAMS: %v\n", c.Queries())
		fmt.Printf("BODY: %v\n", string(c.Body()))

		return c.Next()
	})

	hostAndPort := ""
	// if env.Env.APP_ENV == env.APP_ENV_LOCAL || env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
	// hostAndPort = "127.0.0.1"
	// }

	hostAndPort = hostAndPort + ":" + strconv.Itoa(443)
	app.ListenTLS(hostAndPort, "./crt.crt", "./key.key")
	// app.Use(logger.New())
	// app.Use(middleware.TokenDecrypter)
	// fmt.Println(Hello("LinkShortner"))
}
