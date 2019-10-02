// GoLang による CORS 対応 Web API の実装例
// Usage: sample.exe
// Web API への入力: {"x":整数,"y":整数}
// Web API の出力: {"status":ステータスコード("OK" or "NG"),"value":計算結果の整数,"message":エラー理由}

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
	// 使用する静的コンテンツのパス
	htdocsDir = "htdocs"
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
	// virtual host #1
	http.Handle("example.jp/", http.FileServer(http.Dir(htdocsDir)))

	// virtual host #2
	http.HandleFunc("aaa.jp/api", handlerAPI)
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
	// Origin をチェックしない場合の例
	return true

	// Origin をチェックする場合の例
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
		retData = RetData{Status: "OK", Value: calc(param)}
	} else {
		// 失敗理由を返信データに設定
		retData.Message = err.Error()
	}

	// JSON 形式でクライアントに返信
	retBytes, _ := json.Marshal(retData)
	fmt.Fprintf(w, `%s`, string(retBytes))
}

// calc : 与えられたパラメータで計算を行う関数
func calc(param Param) (r int) {
	r = param.X + param.Y // 計算例
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
