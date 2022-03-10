package main

//UserData 註冊用的使用者資料
type UserData struct {
	ID          string //使用者uid
	UserName    string //帳號
	Password    string //密碼(md5)
	ChineseName string //使用者姓名
	Email       string //使用者mail
	Address     string //使用者地址
	Phone       string //電話
	Point       int    //剩餘的$$
	State       int    //帳號狀態
	LastLogin   string //最後登入時間
}

//LoginData 登入用的struct
type LoginData struct {
	UserName string //帳號
	Password string // 密碼
}

// LoginResult 登入結果
type LoginResult struct {
	Token string `json:"token"`
}

//AndHostNowInfo 放置當前系統資訊
type AndHostNowInfo struct {
	CPUUse      float64 //cpu用量
	MemoryUse   uint64  //記憶體用量
	SwapUse     uint64  // swap用量
	RootDiskUse float64 //以使用總硬碟空間(root)
	UpdateTime  int64   //更新時間
}

//CPUUseInfo cpu使用率結構體
type CPUUseInfo struct {
	IdleTicks  float64
	TotalTicks float64
	CPUUsage   float64
}

//WebResult json的回傳資訊
type WebResult struct {
	Ret    int
	Reason string
	Data   interface{}
}
