// CORS 対応 API ラッパー

package main

import (
	"fmt"
	"net/http"
	"strings"
)

// handlerAPI : CORS 対応 API ラッパー
func handlerAPIWithCORS(w http.ResponseWriter, r *http.Request) {
	// リクエストのダンプ
	dumpRequest(r)

	// Originヘッダのチェック
	if !IsAllowableOrigin(r) {
		// Originを許容できない場合は 403 を返す
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// OPTIONSメソッドのときは preflight request を処理して終了。
	if r.Method == http.MethodOptions {
		processPreFlightRequest(w, r)
		return
	}

	// handlerAPI のレスポンスヘッダに、
	// クロスオリジンを許容するヘッダを追加する。
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	// [MEMO] 上のヘッダがないとブラウザによってブロックされる。
	// 以下は Chrome のデベロッパーツールのコンソールに出力されたメッセージ
	// Cross-Origin Read Blocking (CORB) blocked cross-origin
	// response http://aaa.jp:8080/api with MIME type text/plain.
	// See https://www.chromestatus.com/feature/5629709824032768
	// for more details.

	// [MEMO]
	// 任意のクロスオリジンを許容する場合はワイルドカードでもいいらしい
	// Access-Control-Allow-Origin: *

	// API 処理の呼び出し
	handlerAPI(w, r)
}

// processPreFlightRequest : クロスオリジンでのアクセスを許容する条件を返答する関数
func processPreFlightRequest(w http.ResponseWriter, r *http.Request) {
	// 許容するオリジン
	//ここではクライアントからの Origin の値をそのまま返す
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	// 許容するメソッド
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	// 許容するヘッダ
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// preflight request の応答をキャッシュしてよい時間（秒数）
	w.Header().Set("Access-Control-Allow-Max-Age", "86400")
}

// IsAllowableOrigin : Originヘッダをチェックする関数
func IsAllowableOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	// Origin ヘッダの有無のみチェックする例
	// つまり、任意の Origin を許容する場合
	return origin != ""

	// [XXX] 上は同一オリジンでも Origin ヘッダが付与される前提。
	// Chrome ではそのように動作することを確認したが、他のブラウザは未確認。
	// 場合によっては Origin ヘッダがないときのことも想定すべきかもしれない。

	/*
		// 所定の Origin のみ許容する例
		return origin == "http://example.jp:8080"
	*/
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
