package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/apis"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/env"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/interfaces"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/utility"
)

// fiber middleware for jwt
func TokenDecrypter(c *fiber.Ctx) error {
	reqHeaders := c.GetReqHeaders()
	tokenString, foundToken := reqHeaders[env.RequestTokenHeaderKey]
	if !foundToken || len(tokenString) != 1 || tokenString[0] == "" {
		c.Locals(interfaces.REQ_LOCAL_ERROR_KEY, &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_TOKEN_NOT_BEFORE,
			Message:    "Please Pass Valid Token",
			Name:       "ERROR_TOKEN_NOT_BEFORE",
		})
		return c.Next()
	}
	userRolesCustomClaim, localErr := apis.AccessTokenService.VerifyTokenGeneralPurpose(tokenString[0])
	if localErr != nil {
		c.Locals(interfaces.REQ_LOCAL_ERROR_KEY, localErr)
		return c.Next()
	}
	if e := utility.ValidateStructAndReturnReqError(userRolesCustomClaim, &interfaces.RequestError{
		StatusCode: 401,
		Code:       interfaces.ERROR_INVALID_TOKEN_SIGNATURE,
		Message:    "Invalid Token Structure",
		Name:       "ERROR_INVALID_TOKEN_SIGNATURE",
	}); e != nil {
		c.Locals(interfaces.REQ_LOCAL_ERROR_KEY, e)
		return c.Next()
	}
	c.Locals(interfaces.REQ_LOCAL_KEY_TOKEN_RAW_DATA, userRolesCustomClaim)
	return c.Next()
}

func AllowOnlyValidTokenMiddleWare(c *fiber.Ctx) error {
	jwtRawFromLocal := c.Locals(interfaces.REQ_LOCAL_KEY_TOKEN_RAW_DATA)
	localError := c.Locals(interfaces.REQ_LOCAL_ERROR_KEY)

	if jwtRawFromLocal == nil {
		if localError != nil {
			err, ok := localError.(*interfaces.RequestError)
			if !ok {
				return &interfaces.RequestError{
					StatusCode: 403,
					Code:       interfaces.ERROR_TOKEN_EXPIRED,
					Message:    "Cannot Cast Error Token",
					Name:       "NOT_VALID_DECRYPTED_TOKEN",
				}
			}
			return err
		}
		return &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_TOKEN_EXPIRED,
			Message:    "Invalid Token or token expired",
			Name:       "NOT_VALID_DECRYPTED_TOKEN",
		}
	}

	jwtRaw, ok := jwtRawFromLocal.(*utility.GeneralPurposeTokenGeneration)
	if !ok {
		return &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_TOKEN_EXPIRED,
			Message:    "Invalid Token or token expired",
			Name:       "NOT_VALID_DECRYPTED_TOKEN",
		}
	}
	if jwtRaw.Id == "" {
		return &interfaces.RequestError{
			StatusCode: 403,
			Code:       interfaces.ERROR_INVALID_TOKEN,
			Message:    "Invalid Bullion Id in Token",
			Name:       "INVALID_BULLION_ID_IN_TOKEN",
		}
	}

	return c.Next()
}
