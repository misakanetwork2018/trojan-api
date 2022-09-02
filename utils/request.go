package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ParseParams(c *gin.Context, model interface{}) error {
	err := c.ShouldBind(model)

	if err != nil {
		fmt.Println(err.Error())
		RespondWithError(500, "param error", c)
	}

	return err
}
