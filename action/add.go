package action

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"trojan-api/model"
	"trojan-api/utils"
)

func AddUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user model.User

		if utils.ParseParams(c, &user) != nil {
			return
		}

		command := "set -add-profile -target-password " + user.TargetPassword

		if user.IpLimit != nil {
			command += " -ip-limit " + strconv.Itoa(*user.IpLimit)
		}

		if user.DownloadSpeedLimit != nil {
			command += " -download-speed-limit " + strconv.Itoa(*user.DownloadSpeedLimit)
		}

		if user.UploadSpeedLimit != nil {
			command += " -upload-speed-limit " + strconv.Itoa(*user.UploadSpeedLimit)
		}

		err, stdout, stderr := utils.RunTrojanCLI(command)

		var msg string
		var ok = false

		if err != nil {
			msg = err.Error() + "\n" + stderr
		} else {
			if stdout == "Done\n" {
				ok = true
			} else {
				msg = stdout
			}
		}

		c.JSON(200, gin.H{
			"success": ok,
			"msg":     msg,
		})
	}
}