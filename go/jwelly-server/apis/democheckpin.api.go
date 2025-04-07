package apis

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func DemoCheckPinApi(c *fiber.Ctx) error {
	fmt.Printf("%v", c.Queries())

	return c.JSON(demoApiResponse{
		Register: []demoApiResponse1{
			{
				Msg: "OTP is confirmed for Demo 39636",
			},
		},
	})
}
func DemoCheckPinApi1(c *fiber.Ctx) error {
	// var test interface{}
	// if err := json.Unmarshal(c.Body(), &test); err != nil {
	// 	return err
	// }
	return c.JSON(demoApiResponse{
		Register: []demoApiResponse1{
			{
				Msg: "success",
			},
		},
	})
}

func ErpupgradeApi(c *fiber.Ctx) error {
	return c.SendString("2024-07-15")
}

// softwarekey/getenqoption.php
func GetenqoptionApi(c *fiber.Ctx) error {
	// c.Response().Header.Set("Accept", "application/json")
	return c.SendString("11111")
}
