// GoLang による CORS 対応 Web API の実装例
// Web API をクロスオリジンで提供する場合のメモ
// Usage: sample.exe
// Web API への入力: {"x":整数,"y":整数}
// Web API の出力: {"status":ステータスコード("OK" or "NG"),"value":計算結果の整数,"message":エラー理由}

package main

import (
	"fmt"
	"net/http"
)

const (
	// 使用するポート番号
	port = ":8080"
	// 使用する静的コンテンツのパス
	htdocsDir = "htdocs"
)

func main() {
	fmt.Println(VERSION)

	// 静的コンテンツ用
	http.Handle("/", http.FileServer(http.Dir(htdocsDir)))

	// API 用
	// http.HandleFunc("aaa.jp/api", handlerAPI)
	http.HandleFunc("aaa.jp/api", handlerAPIWithCORS)

	// サーバ起動
	http.ListenAndServe(port, nil)
}
