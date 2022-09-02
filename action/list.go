package action

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"trojan-api/model"
	"trojan-api/utils"
)

func List() func(c *gin.Context) {
	return func(c *gin.Context) {
		var users []model.TrojanUserDetail

		err, stdout, stderr := utils.RunTrojanCLI("list")

		if err != nil {
			fmt.Println("get list error: ", err.Error()+"\n"+stderr)
			utils.RespondWithError(500, err.Error(), c)
			return
		}

		if err = json.Unmarshal([]byte(stdout), &users); err != nil {
			fmt.Println("parse json error: ", err.Error())
			utils.RespondWithError(500, err.Error(), c)
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    users,
		})
	}
}

func GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user model.TrojanSingleUser

		pass := c.Query("pass")

		if pass == "" {
			utils.RespondWithError(500, "Please provide password", c)
			return
		}

		err, stdout, stderr := utils.RunTrojanCLI("get -target-password " + pass)

		if err != nil {
			fmt.Println("get user error: ", err.Error()+"\n"+stderr)
			utils.RespondWithError(500, err.Error(), c)
			return
		}

		if err = json.Unmarshal([]byte(stdout), &user); err != nil {
			fmt.Println("parse json error: ", err.Error())
			utils.RespondWithError(500, err.Error(), c)
			return
		}

		c.JSON(200, gin.H{
			"success": user.Success,
			"msg":     user.Info,
			"data":    user.Status,
		})
	}
}
