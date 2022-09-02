package action

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"trojan-api/model"
	"trojan-api/utils"
)

func DelUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			msg  string
			ok   bool
			stat int
		)

		pass := c.PostForm("pass")

		// get last traffic info before delete
		err, stdout, stderr := utils.RunTrojanCLI("get -target-password " + pass)

		if err == nil {
			var user model.TrojanUserDetail
			if err = json.Unmarshal([]byte(stdout), &user); err != nil {
				fmt.Println("parse json error: ", err.Error())
				return
			}

			stat = user.Status.TrafficTotal.DownloadTraffic + user.Status.TrafficTotal.UploadTraffic

			err, stdout, stderr = utils.RunTrojanCLI("set -delete-profile -target-password " + pass)

			if err != nil {
				msg = err.Error() + "\n" + stderr
			} else {
				if stdout == "Done\n" {
					ok = true
				} else {
					msg = stdout
				}
			}
		} else {
			msg = err.Error() + "\n" + stderr
		}

		c.JSON(200, gin.H{
			"success":  ok,
			"msg":      msg,
			"transfer": stat,
		})
	}
}
