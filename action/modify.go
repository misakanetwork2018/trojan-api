package action

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"trojan-api/model"
	"trojan-api/utils"
)

func ModifyUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user model.User

		if utils.ParseParams(c, &user) != nil {
			return
		}

		if user.TargetPassword == "" && user.Id == 0 {
			utils.RespondWithError(500, "Please provide password or id", c)
			return
		}

		command := "set -modify-profile -target-password " + user.TargetPassword

		if user.IpLimit == nil && user.DownloadSpeedLimit == nil && user.UploadSpeedLimit == nil {
			utils.RespondWithError(500, "At least provide one modify item", c)
			return
		}

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
