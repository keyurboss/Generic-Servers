package apis

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/interfaces"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/utility"
	// "github.com/keyurboss/Generic-Servers/go/lint-shortner/firebase"
	// "github.com/keyurboss/Generic-Servers/go/lint-shortner/interfaces"
)

type bodyBulkShortning struct {
	Urls []string `json:"urls" validate:"required,min=1"`
}

type bodySingleShortning struct {
	Urls string `json:"url" validate:"required"`
}

func BulkShortning(c *fiber.Ctx) error {
	body := new(bodyBulkShortning)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	resp := make(map[string]string)
	for _, url := range body.Urls {
		if !utility.ValidateUrl(url) {
			continue
		}
		id := utility.UniqueIdConst.GetUniqueId()
		resp[url] = utility.ServerUrl + id
		firebaseDb.SetPublicData(id, map[string]interface{}{"url": url, "created_at": time.Now().Unix()})
	}
	return c.JSON(fiber.Map{
		"success": 1,
		"data":    resp,
	})
}

func SignleShortning(c *fiber.Ctx) error {
	body := new(bodySingleShortning)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	// resp := make(map[string]string)

	// for _, url := range body.Urls {
	if !utility.ValidateUrl(body.Urls) {
		return c.Status(400).JSON(interfaces.RequestError{
			Code:    interfaces.ERROR_INVALID_INPUT,
			Message: "Please Pass Valid URL",
			Name:    "INVALID URL",
		})
	}
	id := utility.UniqueIdConst.GetUniqueId()
	newUrl := utility.ServerUrl + id
	firebaseDb.SetPublicData(id, map[string]interface{}{"url": body.Urls, "created_at": time.Now().Unix()})
	// }
	return c.JSON(fiber.Map{
		"success": 1,
		"data":    newUrl,
	})
	// return c.JSON(resp)
}
