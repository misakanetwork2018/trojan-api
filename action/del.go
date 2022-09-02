package action

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"trojan-api/model"
	"trojan-api/utils"
)

func DelUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			msg         string
			ok          bool
			err         error
			errMsg      string
			hash        string
			upTraffic   int
			downTraffic int
		)

		var tUser model.TrojanSingleUser

		pass := c.PostForm("pass")
		id := c.PostForm("id")
		var iId = 0

		if pass == "" && id == "" {
			utils.RespondWithError(500, "Please provide password or id", c)
			return
		}

		if pass == "" {
			iId, err = strconv.Atoi(id)
			if err != nil {
				utils.RespondWithError(500, err.Error(), c)
				return
			}
			if hash, ok = IdHash[iId]; ok {
				ok, errMsg = getTrojanUserByHash(hash, &tUser)
			} else {
				c.JSON(200, gin.H{
					"success": false,
					"msg":     "id not exist",
				})
				return
			}
		} else {
			ok, errMsg = getTrojanUser(pass, &tUser)
		}

		if !ok {
			utils.RespondWithError(500, errMsg, c)
			return
		}

		if tUser.Status == nil {
			c.JSON(200, gin.H{
				"success": tUser.Success,
				"msg":     tUser.Info,
			})
			return
		}

		upTraffic = tUser.Status.TrafficTotal.UploadTraffic
		downTraffic = tUser.Status.TrafficTotal.DownloadTraffic

		hash = tUser.Status.User.Hash
		if us, ok := Users[hash]; ok {
			if upTraffic >= us.OldTrafficUp {
				upTraffic = upTraffic - us.OldTrafficUp
			}
			if downTraffic >= us.OldTrafficDown {
				downTraffic = downTraffic - us.OldTrafficUp
			}
		}

		err, stdout, stderr := utils.RunTrojanCLI("set -delete-profile -target-hash " + hash)

		if err != nil {
			msg = err.Error() + "\n" + stderr
		} else {
			if stdout == "Done\n" {
				ok = true
				delete(Users, hash)
				delete(IdHash, iId)
			} else {
				msg = stdout
			}
		}

		c.JSON(200, gin.H{
			"success": ok,
			"msg":     msg,
			"transfer": gin.H{
				"upload":   upTraffic,
				"download": downTraffic,
			},
		})
	}
}
