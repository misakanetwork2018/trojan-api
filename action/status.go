package action

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	"trojan-api/utils"
)

func Status() func(c *gin.Context) {
	return func(c *gin.Context) {
		var loadavg, uptime string

		data, err := ioutil.ReadFile("/proc/loadavg")
		if err != nil {
			fmt.Println("File reading error", err)
			utils.RespondWithError(500, "Cannot get load", c)
			return
		}
		loadavg = strings.Replace(string(data), "\n", "", 1)

		data, err = ioutil.ReadFile("/proc/uptime")
		if err != nil {
			fmt.Println("File reading error", err)
			utils.RespondWithError(500, "Cannot get uptime", c)
			return
		}
		uptime = strings.Split(string(data), " ")[0]

		c.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"loadavg": loadavg,
				"uptime":  uptime,
			},
		})
	}
}
