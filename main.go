// Sample of CORS using GoLang
// Usage: sample.exe
// Web API:
// INPUT: {"x":INTVALUE, "y":INTVALUE}
// OUTPUT: {"status":statuscode("OK" or "NG"), "value":Icalc(x, y), "message":"reason of error"}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	// 使用するポート番号
	port = ":8080"
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

func main() {
	fmt.Println(VERSION)
	http.HandleFunc("/api", handlerAPI)
	http.ListenAndServe(port, nil)
}

// handlerAPI : CORS 対応 API 処理関数
func handlerAPI(w http.ResponseWriter, r *http.Request) {
	dumpRequest(r)

	// Originヘッダのチェック
	if !checkOrigin(r) {
		// Originを許容できない場合は 403 を返す
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// メソッドによって分岐する
	switch r.Method {

	case http.MethodOptions: // preflight request 用
		processOPTIONS(w, r)
		return

	case http.MethodPost: // API用
		processPOST(w, r)
		return

	}

	// その他のメソッドの場合は 405 を返す
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// checkOrigin : Originヘッダをチェックする
func checkOrigin(r *http.Request) bool {
	// Origin の有無のみチェックする場合の例
	return r.Header.Get("Origin")!=""

	// Origin の値をチェックする場合の例
	//return r.Header.Get("Origin") == "http://example.jp:8080"
}

// processOPTIONS : preflight request 用 OPTIONS メソッドの処理
func processOPTIONS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Max-Age", "86400")
}

// processPOST : API 用 POST メソッドの処理
func processPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Content-Type", "application/json")

	// 返信データの用意。デフォルトは失敗のデータ。
	retData := RetData{Status: "NG"}

	// JSON形式のパラメータを取得
	param, err := getParam(r)

	// パラメータの取得に成功したとき
	if err == nil {
		// calc の計算結果を返信データに設定
		retData = RetData{Status: "OK", Value: calc(param.x, param.y)}
	} else {
		// 失敗理由を返信データに設定
		retData.Message = err.Error()
	}

	// JSON 形式でクライアントに返信
	retBytes, _ := json.Marshal(retData)
	fmt.Fprintf(w, `%s`, string(retBytes))
}

// calc : 与えられたパラメータで計算を行う関数
func calc(x, y int) (r int) {
	r = x + y // 計算例
	return
}

// getParam : リクエストからパラメータを取得する関数
func getParam(r *http.Request) (param Param, err error) {
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

// dumpRequest : 受信したリクエストのダンプ
func dumpRequest(r *http.Request) {
	fmt.Println("------------------------------------")
	fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
	fmt.Println("Host:", r.Host)
	for k, v := range r.Header {
		fmt.Println(k, ":", strings.Join(v, ", "))
	}
}
