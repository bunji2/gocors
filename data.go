// API 用データ定義

package main

import (
	"encoding/json"
	"fmt"
)

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

// JSONString : JSON形式の文字列を生成
func (rd RetData) JSONString() string {
	bytes, err := json.Marshal(rd)
	if err != nil {
		return fmt.Sprintf(`{status:"NG", message:"%s"}`, err.Error())
	}
	return string(bytes)
}
