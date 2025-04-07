package apis

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type demoApiResponse struct {
	Register []demoApiResponse1 `json:"register"`
}
type demoApiResponse1 struct {
	Msg string `json:"msg"`
}

func DemoApi(c *fiber.Ctx) error {
	fmt.Printf("%v", c.Queries())

	return c.JSON(demoApiResponse{
		Register: []demoApiResponse1{
			{
				Msg: "success",
			},
		},
	})
}
