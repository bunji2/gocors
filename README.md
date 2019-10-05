# gocors --- CORS (Cross-Origin Resource Sharing) の GoLang によるサンプル

あるコンテンツから "XMLHttpRequest" から WebAPI にアクセスすることを考える。

![fig0](images/fig0.png)

これは "Same-Origin Policy" の場合に成立する。つまり最初のコンテンツと、WebAPI が同じ生成元の場合である。

で、例えばこの WebAPI が有用で別のオリジンのコンテンツからも使えるようにすることを考える。

通常は同一でないオリジンで WebAPI にアクセスしようとするとブラウザにブロックされてしまう。その際にブラウザのコンソールには次のようなメッセージが表示される。

```
Access to XMLHttpRequest at 'http://aaa.jp:8080/api' from origin 'http://example.jp:8080' has been blocked by CORS policy: Response to preflight request doesn't pass access control check: No 'Access-Control-Allow-Origin' header is present on the requested resource.
```

オリジンが 'example.jp' でここから XMLHttpRequest で 'aaa.jp' の WebAPI にアクセスする前に、"preflight request" のレスポンスが "CORS policy" を満たさずブロックされてしまう。

![fig1](images/fig1.png)

ここでクロスオリジンな状況でも 'aaa.jp' が WebAPI を提供するには "OPTIONS" メソッドからなる "preflight request" に対してオリジンを許可するレスポンスを返す必要があり、具体的には以下のレスポンスヘッダを追加しなければならない。

|レスポンスヘッダ|概要|
|:--|:--|
|Access-Control-Allow-Origin|アクセスを許容するオリジン|
|Access-Control-Allow-Methods|アクセスを許容するメソッド群|
|Access-Control-Allow-Headers|アクセスを許容するヘッダ群|

以下に例を示す。

```
Access-Control-Allow-Origin: http://example.jp:8080
Access-Control-Allow-Methods: POST, OPTIONS
Access-Control-Allow-Headers: Content-Type
```

"Access-Control-Allow-Methods" で "OPTIONS" メソッドを追加している。また、WebAPI で使う "POST" メソッドで入力データのコンテンツタイプに JSON データを想定しているため、"Access-Control-Allow-Headers" で "Content-Type" を指定する。

まず、用意する WebAPI を示す。想定として JSON 形式のデータを入力として受け取り、結果を JSON 形式で返答するものとする。

```
func main() {
	http.HandleFunc("aaa.jp/api", handlerAPI)
	http.ListenAndServe(port, nil)
}

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
	}

	w.Header().Set("Content-Type", "application/json")

	// JSON 形式でクライアントに返信
	fmt.Fprintf(w, `%s`, retData.JSONString())
}
```

通常、同一オリジンで使用する WebAPI ならば上記のように "POST" メソッドのみ実装すればよいが、クロスオリジンで使用するには "preflight request" に応答できなければならない。以下に "preflight request" への対応例を示す。

```
func main() {
	// http.HandleFunc("aaa.jp/api", handlerAPI)
	http.HandleFunc("aaa.jp/api", handlerAPIWithCORS)
	http.ListenAndServe(port, nil)
}

func handlerAPI(w http.ResponseWriter, r *http.Request) {
    // 省略
}

func handlerAPIWithCORS(w http.ResponseWriter, r *http.Request) {
	// Origin ヘッダのチェック
	if !IsAllowableOrigin(r) {
		// Origin を許容できない場合は 403 を返す
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// OPTIONS メソッドのときは preflight request を処理して終了。
	if r.Method == http.MethodOptions {
		processPreFlightRequest(w, r)
		return
	}

	// handlerAPI のレスポンスにオリジンを許容するヘッダを追加する。
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	// API 処理の呼び出し
	handlerAPI(w, r)
}

func processPreFlightRequest(w http.ResponseWriter, r *http.Request) {
	// クライアントからの Origin を許容する
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

```

ポイントは以下の二点である。

* "processPreFlightRequest" 関数により、"preflight request" に対してオリジンを許可するレスポンスを返す。

* "handlerAPIWithCORS" 関数により　"handlerAPI" をラップし、オリジンを許可するレスポンスを追加する。

![fig2](images/fig2.png)
