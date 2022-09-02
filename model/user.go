package model

type User struct {
	Id                 int    `form:"id"`
	TargetPassword     string `form:"pass" binding:"required"`
	IpLimit            *int   `form:"ip_limit" binding:"required"`
	UploadSpeedLimit   *int   `form:"upload_speed_limit" binding:"required"`
	DownloadSpeedLimit *int   `form:"download_speed_limit" binding:"required"`
}

type UserDetail struct {
	Id                 int    `json:"id"`
	TargetHash         string `json:"target_hash"`
	UploadTraffic      int    `json:"upload_traffic"`
	DownloadTraffic    int    `json:"download_traffic"`
	UploadSpeed        int    `json:"upload_speed"`
	DownloadSpeed      int    `json:"download_speed"`
	UploadSpeedLimit   int    `json:"upload_speed_limit"`
	DownloadSpeedLimit int    `json:"download_speed_limit"`
	IPLimit            int    `json:"ip_limit"`
}
