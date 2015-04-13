package utils

//处理mac表的结构体 ->正向 ,公司 : mac address
type MacItem struct {
	Hex     string `json:"hex"`     //hex模式的mac段 00-00-01
	Base16  string `json:"base"`    //base 16模式的mac段 000000
	Company string `json:"company"` //公司
}

type ComMac struct {
	Company string  `json:"company"` //公司
	Mac     MacItem `json:"mac"`     //mac结构体
}
