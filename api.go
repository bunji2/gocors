// API 処理

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// handlerAPI : API 処理関数
func handlerAPI(w http.ResponseWriter, r *http.Request) {

	// POSTメソッド意外の場合は 405 を返す
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 返信データの用意
	retData := NewRetData()

	// JSON形式のパラメータを取得
	param, err := getParamFromJSON(r)

	// パラメータの取得に成功したとき
	if err == nil {
		// calc の計算結果を返信データに設定
		retData.Status = "OK"
		retData.Value = calc(param.X, param.Y)
	} else {
		// 失敗理由を返信データに設定
		retData.Message = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")

	// JSON 形式でクライアントに返信
	fmt.Fprintf(w, `%s`, retData.JSONString())
}

// calc : 与えられたパラメータで計算を行う関数
func calc(x, y int) (r int) {
	r = x + y // 計算例
	return
}

// getParamFromJSON : リクエストからJSON形式のパラメータを取得する関数
func getParamFromJSON(r *http.Request) (param Param, err error) {

	// コンテンツタイプが application/json でないときのエラー
	if r.Header.Get("Content-Type") != "application/json" {
		err = fmt.Errorf("content-type is not application/json")
		return
	}

	param = Param{}

	var clen int

	// Content-Length の値を整数で取得
	tmp := r.Header.Get("Content-Length")
	clen, err = strconv.Atoi(tmp)

	// Atoi のエラー
	if err != nil {
		return
	}

	// コンテンツレングスのサイズが空
	if clen <= 0 {
		err = fmt.Errorf("content-length is empty")
		return
	}

	// POST データのバイト列を取得
	bytes := make([]byte, clen)
	var nRead int
	nRead, err = r.Body.Read(bytes)

	// EOF のとき(これ以上データがないとき)は err を nil に設定
	if err == io.EOF {
		err = nil
	}

	// Read のエラー
	if err != nil {
		return
	}

	// 読み出したサイズが足りないエラー
	if nRead < clen {
		err = fmt.Errorf("read data is too few")
		return
	}

	// バイト列を JSON 形式としてパース
	err = json.Unmarshal(bytes[:clen], &param)

	return
}
