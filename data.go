// API 用データ定義

package main

// Param : クライアントから受信するパラメータの型
type Param struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// RetData : クライアントに返信するデータの型
type RetData struct {
	Status  string `json:"status"`
	Value   int    `json:"value"`
	Message string `json:"message"`
}

// NewRetData : 返信データの新規作成
func NewRetData() (r RetData) {
	r = RetData{Status: "NG"}
	return
}
