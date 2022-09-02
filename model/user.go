package model

type User struct {
	TargetPassword     string `form:"pass" binding:"required"`
	IpLimit            *int   `form:"ip_limit" binding:"required"`
	UploadSpeedLimit   *int   `form:"upload_speed_limit" binding:"required"`
	DownloadSpeedLimit *int   `form:"download_speed_limit" binding:"required"`
}
