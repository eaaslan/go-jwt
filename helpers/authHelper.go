package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func MatchUserTypeToUid(c *gin.Context, userId string) error {

	userType := c.GetString("user_type")
	uid := c.GetString("uid")

	if userType == "USER" && uid != userId {
		err := errors.New("unauthorized to access this resource")
		return err
	}

	return nil
}
func CheckUserType(c *gin.Context, role string) error {
	userType := c.GetString("user_type")
	if userType != role {
		return errors.New("unauthorized to access")
	}
	return nil
}
