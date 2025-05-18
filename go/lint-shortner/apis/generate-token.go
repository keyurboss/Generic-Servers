package apis

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/interfaces"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/utility"
)

type bodyGenerateToken struct {
	Uname    string `json:"uname" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func GenerateToken(c *fiber.Ctx) error {
	body := new(bodyGenerateToken)
	c.BodyParser(body)
	if err := utility.ValidateReqInput(body); err != nil {
		return err
	}
	if body.Uname != utility.UName {
		return c.Status(403).JSON(interfaces.RequestError{
			Code:    interfaces.ERROR_INVALID_INPUT,
			Message: "Invalid Uname",
			Name:    "INVALID Uname",
		})
	}
	if body.Password != utility.Password {
		return c.Status(403).JSON(interfaces.RequestError{
			Code:    interfaces.ERROR_INVALID_PASSWORD,
			Message: "Invalid Password",
			Name:    "INVALID PASSWORD",
		})
	}
	now := time.Now()
	signed, err := AccessTokenService.GenerateToken(&utility.GeneralPurposeTokenGeneration{
		RegisteredClaims: &jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Minute * 2)},
		},
		Id: strings.ReplaceAll(uuid.New().String(), "-", ""),
	})
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, )
	if err != nil {
		return c.JSON(interfaces.RequestError{
			Code:    interfaces.ERROR_INTERNAL_SERVER,
			Message: "Some Internal Error",
			Name:    "Generate Token",
			Extra:   err,
		})
	}
	return c.JSON(fiber.Map{
		"token": signed,
	})
}
