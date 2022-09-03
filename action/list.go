package action

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"trojan-api/model"
	"trojan-api/utils"
)

type UserStore struct {
	OldTrafficUp   int
	OldTrafficDown int
	Id             int
}

var (
	Users  = make(map[string]UserStore)
	IdHash = make(map[int]string)
)

func List() func(c *gin.Context) {
	return func(c *gin.Context) {
		var tUsers []model.TrojanUserDetail

		err, stdout, stderr := utils.RunTrojanCLI("list")

		if err != nil {
			fmt.Println("get list error: ", err.Error()+"\n"+stderr)
			utils.RespondWithError(500, err.Error(), c)
			return
		}

		if err = json.Unmarshal([]byte(stdout), &tUsers); err != nil {
			fmt.Println("parse json error: ", err.Error())
			utils.RespondWithError(500, err.Error(), c)
			return
		}

		// trojan-go提供的信息不能直接拿来用，而是需要我们进行转换
		var users = make([]model.UserDetail, 0)
		for _, tUser := range tUsers {
			var user model.UserDetail
			transformTrojanUserToApiUser(&tUser.Status, &user)
			users = append(users, user)
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    users,
		})
	}
}

func GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var tUser model.TrojanSingleUser

		pass := c.Query("pass")
		id := c.Query("id")

		if pass == "" && id == "" {
			utils.RespondWithError(500, "Please provide password or id", c)
			return
		}

		var (
			ok     = false
			errMsg string
			hash   string
		)
		if pass == "" {
			iId, err := strconv.Atoi(id)
			if err != nil {
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
				"data":    nil,
			})
			return
		}

		var user model.UserDetail
		transformTrojanUserToApiUser(tUser.Status, &user)

		c.JSON(200, gin.H{
			"success": tUser.Success,
			"msg":     tUser.Info,
			"data":    user,
		})
	}
}

func getTrojanUser(pass string, user *model.TrojanSingleUser) (bool, string) {
	err, stdout, stderr := utils.RunTrojanCLI("get -target-password " + pass)

	if err != nil {
		fmt.Println("get user error: ", err.Error()+"\n"+stderr)
		return false, err.Error() + "\n" + stderr
	}

	if err = json.Unmarshal([]byte(stdout), user); err != nil {
		fmt.Println("parse json error: ", err.Error())
		return false, err.Error()
	}

	return true, ""
}

func getTrojanUserByHash(hash string, user *model.TrojanSingleUser) (bool, string) {
	err, stdout, stderr := utils.RunTrojanCLI("get -target-hash " + hash)

	if err != nil {
		fmt.Println("get user error: ", err.Error()+"\n"+stderr)
		return false, err.Error() + "\n" + stderr
	}

	if err = json.Unmarshal([]byte(stdout), user); err != nil {
		fmt.Println("parse json error: ", err.Error())
		return false, err.Error()
	}

	return true, ""
}

func transformTrojanUserToApiUser(status *model.TrojanStatus, user *model.UserDetail) {
	hash := (*status).User.Hash
	(*user).Id = Users[hash].Id
	(*user).TargetHash = hash
	(*user).DownloadTraffic = (*status).TrafficTotal.DownloadTraffic
	(*user).UploadTraffic = (*status).TrafficTotal.UploadTraffic
	if (*user).UploadTraffic > Users[hash].OldTrafficUp {
		(*user).UploadTraffic -= Users[hash].OldTrafficUp
	}
	if (*user).DownloadTraffic > Users[hash].OldTrafficDown {
		(*user).DownloadTraffic -= Users[hash].OldTrafficDown
	}
	(*user).DownloadSpeed = (*status).SpeedCurrent.DownloadSpeed
	(*user).UploadSpeed = (*status).SpeedCurrent.UploadSpeed
	(*user).DownloadSpeedLimit = (*status).SpeedLimit.DownloadSpeed
	(*user).UploadSpeedLimit = (*status).SpeedLimit.UploadSpeed
	(*user).IPLimit = (*status).IPLimit
	var us = Users[hash]
	us.OldTrafficUp = (*status).TrafficTotal.UploadTraffic
	us.OldTrafficDown = (*status).TrafficTotal.DownloadTraffic
	Users[hash] = us
}
