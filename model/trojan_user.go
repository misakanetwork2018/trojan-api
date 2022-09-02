package model

type TrojanUserDetail struct {
	Status TrojanStatus `json:"status"`
}

type TrojanSingleUser struct {
	Success bool          `json:"success"`
	Info    string        `json:"info"`
	Status  *TrojanStatus `json:"status,omitempty"`
}

type TrojanUser struct {
	Hash string `json:"hash"`
}

type TrojanStatus struct {
	User         TrojanUser         `json:"user"`
	TrafficTotal TrojanTraffic      `json:"traffic_total"`
	SpeedCurrent TrojanSpeedCurrent `json:"speed_current"`
	SpeedLimit   TrojanSpeedLimit   `json:"speed_limit"`
	IPLimit      int                `json:"ip_limit"`
}

type TrojanTraffic struct {
	UploadTraffic   int `json:"upload_traffic"`
	DownloadTraffic int `json:"download_traffic"`
}

type TrojanSpeedCurrent struct {
	UploadSpeed   int `json:"upload_speed"`
	DownloadSpeed int `json:"download_speed"`
}

type TrojanSpeedLimit struct {
	UploadSpeed   int `json:"upload_speed"`
	DownloadSpeed int `json:"download_speed"`
}
