package apis

import (
	"fmt"
	"time"

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

//	{
//	    "register": [
//	        {
//	            "otp": "",
//	            "msg": "9428393489"
//	        }
//	    ]
//	}
func RegisterDeviceApi(c *fiber.Ctx) error {
	c.Response().Header.Set("Content-Type", "application/json")
	c.Response().Header.Set("X-Powered-By", "PHP/7.4.33")
	c.Response().Header.Set("Server", "Apache/2.4.62 () OpenSSL/1.0.2k-fips PHP/7.4.33")
	c.Response().Header.Set("Upgrade", "h2,h2c")

	return c.JSON(fiber.Map{
		"register": []fiber.Map{
			{
				"otp": "",
				"msg": "9428393489",
			},
		},
	})
}

// softwarekey/getenqoption.php
func GetenqoptionApi(c *fiber.Ctx) error {
	// c.Response().Header.Set("Accept", "application/json")
	return c.SendString("11111")
}

func RegClient(c *fiber.Ctx) error {
	//		{
	//	    "register": [
	//	        "9428393489"
	//	    ]
	//	}
	return c.JSON(fiber.Map{
		"register": []string{
			"9428393489",
		},
	})
}

func GetAmcApi(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"details": []fiber.Map{
			{
				"success": "1",
				"amcdue":  time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
				"amcext":  "0",
			},
		},
	})
}

// "details": [
// 	{
// 		"success": "1",
// 		"amcdue": "2025-04-10",
// 		"amcext": "0"
// 	}
// ]
