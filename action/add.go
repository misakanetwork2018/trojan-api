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

		if user.TargetPassword == "" {
			c.JSON(200, gin.H{
				"success": false,
				"msg":     "Your should provide password",
			})
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
				// 查询其Hash
				var tUser model.TrojanSingleUser
				getTrojanUser(user.TargetPassword, &tUser)
				if tUser.Success {
					var us UserStore
					hash := tUser.Status.User.Hash
					us.OldTrafficDown = 0
					us.OldTrafficUp = 0
					us.Id = user.Id
					Users[hash] = us
					IdHash[user.Id] = hash
				}
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
