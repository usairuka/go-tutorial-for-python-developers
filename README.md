# Python開発者のためのGo言語チュートリアル

Pythonでの開発経験を足がかりに、Goの基本からテスト・並行処理・HTTP APIまでを学ぶサンプルコード集です。小さなプログラムを動かしながら、最後に総合演習「Todo HTTP API」の仕組みを確認します。

- [最初に実行する](#最初に実行する)
- [Pythonとの違いをつかむ](#pythonとの違いをつかむ)
- [学習の進め方](#学習の進め方)
- [章とファイルの対応](#章とファイルの対応)
- [演習に取り組む](#演習に取り組む)
- [Todo HTTP APIを動かす](#todo-http-apiを動かす)
- [全体の確認と理解度チェック](#全体の確認と理解度チェック)
- [困ったとき](#困ったとき)

## 最初に実行する

必要なものは **Go 1.22以降**、エディター、ターミナルです。Goの導入方法は[公式インストールガイド](https://go.dev/doc/install)を参照してください。第6章のAPIは、[Go 1.22で追加されたHTTPメソッド付きルーティング](https://go.dev/doc/go1.22#net/http)を使います。

Gitを使ってサンプルを取得します。Gitがない場合は、GitHubの **Code → Download ZIP** からダウンロードして展開できます。

```sh
git clone https://github.com/usairuka/go-tutorial-for-python-developers.git
cd go-tutorial-for-python-developers/go-lab
```

取得済みの場合は、このREADMEがあるディレクトリから `cd go-lab` で移動してください。

**以降のGoコマンドは、すべて `go-lab/` 内で実行します。**

```sh
go version
go env GOMOD
go run ./cmd/hello
go test ./...
```

次の結果になれば、学習を始める準備は完了です。

| コマンド | 確認すること |
| --- | --- |
| `go version` | Go 1.22以降がインストールされている |
| `go env GOMOD` | このリポジトリの `go-lab/go.mod` へのパスが表示される |
| `go run ./cmd/hello` | `Hello, Go!` が表示される |
| `go test ./...` | `internal/task` と `internal/httpapi` が `ok` になる。`cmd/` の `[no test files]` は正常 |

[go.mod](go-lab/go.mod) は用意済みなので、`go mod init` の実行は不要です。module名の `example.com/go-lab` は教材用の識別名で、ドメインの取得や公開は必要ありません。Goのmoduleは関連するpackageと依存関係を管理する単位です。この教材では `venv` の作成や `pip install` に相当する準備はありません。

```text
go-tutorial-for-python-developers/
├── README.md
├── LICENSE
└── go-lab/                  # Goコマンドはここで実行
    ├── go.mod               # module名とGoのバージョン要件
    ├── cmd/                 # 章ごとに独立した実行プログラム
    │   ├── hello/
    │   ├── basics/
    │   ├── collections/
    │   ├── model/
    │   ├── cleanup/
    │   ├── generics/
    │   ├── concurrency/
    │   └── api/
    └── internal/
        ├── task/            # Todoの型・検証・保存とテスト
        └── httpapi/         # HTTPハンドラーとテスト
```

## Pythonとの違いをつかむ

対応する概念を手がかりに、サンプルで次の違いを確かめてください。

| Pythonで馴染みのあるもの | Goで学ぶこと | 確認する章 |
| --- | --- | --- |
| 動的な型、型ヒント | 静的な型検査。`:=` でも型は決まり、宣言だけの変数には型ごとのゼロ値が入る | 1 |
| `for`・`while`、リスト内包表記 | 繰り返しは `for`。要素の抽出は `for` と `append` で書く | 1・2 |
| `list` とスライス | Goのsliceは要素型を持つ。Pythonのリストのスライスは新しいリストを作るが、Goのsliceは元の配列を共有する | 2 |
| `dict` | `map`。存在しないキーの読み取りはゼロ値を返し、`value, ok := m[key]` で有無を確認する。反復順序は保証されない | 2 |
| `str` と `len()` | Goの `len(string)` はバイト数。UTF-8文字列のコードポイント数は `utf8.RuneCountInString` で数える | 2 |
| クラス、オブジェクトへの参照 | `struct` とメソッド。structの代入は値のコピーになり、元の値を変更する場面でポインタを使う | 3 |
| ダックタイピング、`typing.Protocol` | `interface`。必要なメソッドを持つ型が暗黙に実装し、コンパイル時に適合性を検査する | 3 |
| `raise`・`try/except` | 失敗を戻り値の `error` で伝え、呼び出し側で `err != nil` を確認する | 3 |
| `with`・`try/finally` | 後処理の登録には `defer`。実行は関数の終了時で、登録と逆順になる | 3 |
| `unittest`・`pytest`、型パラメーター | 標準の `testing` と `go test`、テーブル駆動テスト、ジェネリクス | 4 |
| `asyncio`・`threading` | goroutine、channel、`context`、`sync.Mutex`。キャンセルへの応答と共有データの保護を明示する | 5 |

特に、sliceの共有、値のコピー、エラーの戻り値はPythonの感覚との違いが出やすい部分です。structのコピーも、sliceやmapなどをフィールドに持つ場合は参照先まで複製する「深いコピー」にはなりません。

言語仕様の詳細は[Goの仕様](https://go.dev/ref/spec)と[Pythonの組み込み型](https://docs.python.org/3/library/stdtypes.html)を参照してください。

## 学習の進め方

1. 第0章で実行環境を確認し、章の順にサンプルを開きます。
2. 出力を予想してから実行し、入力値や条件を変えて結果を比べます。
3. [演習](#演習に取り組む)の要件を読み、収録済みの実装を書き換えるか、別の作業用コピーで実装します。
4. 第6章でAPIへのリクエストとテストを実行し、最後に理解度チェックで振り返ります。

以下は学習計画の一例です。学習280分と、第3章・第5章の後の休憩各10分で計300分になります。

| 章 | 目安 | 学習の焦点 |
| --- | --- | --- |
| 0. 環境設定と最初の実行 | 15分 | module・package、実行・整形・ビルド |
| 1. 型・変数・制御構文・関数 | 35分 | 型推論、ゼロ値、複数の戻り値 |
| 2. slice・map・文字列 | 35分 | コレクションの操作、コピーと共有、UTF-8 |
| 3. struct・ポインタ・interface・error | 40分 | Todoのモデル、メソッド、エラー処理、後処理 |
| 4. package設計・テスト・ジェネリクス | 35分 | 責務の分割、正常系と境界条件のテスト |
| 5. goroutine・channel・context・排他制御 | 45分 | 完了待ち、結果の収集、キャンセル、競合の防止 |
| 6. 総合演習：Todo HTTP API | 65分 | JSON入出力、保存処理、HTTPテスト |
| 7. 理解度チェック | 10分 | 動作とその理由を自分の言葉で説明する |

## 章とファイルの対応

リンクはこのREADMEからの相対パス、実行コマンドは **`go-lab/` 内** を基準にしています。`cmd/` の各ディレクトリは別々のプログラムです。`go run ./cmd/hello` のように1つずつ実行してください。

| 章 | サンプル・テスト | 実行・確認コマンド |
| --- | --- | --- |
| 0 | [Hello World](go-lab/cmd/hello/main.go) | `go run ./cmd/hello` |
| 1 | [型・関数と偶数の合計](go-lab/cmd/basics/main.go) | `go run ./cmd/basics` |
| 2 | [slice・map・文字列](go-lab/cmd/collections/main.go) | `go run ./cmd/collections` |
| 3 | [Taskとメソッド](go-lab/internal/task/task.go)、[ポインタとinterface](go-lab/cmd/model/main.go) | `go run ./cmd/model` |
| 3 | [deferの評価と実行順](go-lab/cmd/cleanup/main.go) | `go run ./cmd/cleanup` |
| 4 | [Taskのテスト](go-lab/internal/task/task_test.go) | `go test ./internal/task -v` |
| 4 | [ジェネリクス](go-lab/cmd/generics/main.go) | `go run ./cmd/generics` |
| 5 | [並行処理とキャンセル](go-lab/cmd/concurrency/main.go) | `go run ./cmd/concurrency` |
| 6 | [メモリ上の保存処理](go-lab/internal/task/memory.go)、[並行アクセスのテスト](go-lab/internal/task/memory_test.go) | `go test -race ./internal/task` |
| 6 | [HTTPハンドラー](go-lab/internal/httpapi/handler.go)、[HTTPテスト](go-lab/internal/httpapi/handler_test.go) | `go test ./internal/httpapi -v` |
| 6 | [APIサーバーの起動](go-lab/cmd/api/main.go) | `go run ./cmd/api` |

`-race` の実行条件は[全体の確認](#全体の確認と理解度チェック)を参照してください。第7章はこのREADMEの理解度チェックで振り返ります。

## 演習に取り組む

**収録コードには、以下の演習の解答例が組み込まれています。** 自分で実装する場合は、対応する関数やテストを書き換えてください。同名の関数を追加すると重複定義になります。元のコードを比較用に残したい場合は、`go-lab/` を別の作業用ディレクトリにコピーして取り組めます。

| 章 | 演習の要件 |
| --- | --- |
| 1 | `sumEven(limit int) int`：1から `limit` 以下の偶数を合計する。`6` なら `12`、`0` なら `0` を返す |
| 2 | `filterAtLeast(values []int, min int) []int`：`min` 以上の要素を元の順序で取り出す。入力sliceの要素を変更しない |
| 3・4 | `Task.Rename`：前後の空白を除去して名前を更新する。空になる入力では `ErrEmptyTitle` を返し、元の名前を保つ。成功と失敗をテストする |
| 5 | 期限切れの `context` を `doubleAll` に渡し、返るエラーが `context.DeadlineExceeded` か確認する |
| 6 | `GET /health` でHTTP 200と `{"status":"ok"}` を返し、`httptest` で確認する |

<details>
<summary>解答例の場所と動作確認を見る</summary>

コード内の `掲載済みの演習コード` コメントが目印です。

| 演習 | 解答例の場所（`go-lab/` 内） | 確認方法 |
| --- | --- | --- |
| 偶数の合計 | `cmd/basics/main.go` の `sumEven` | `go run ./cmd/basics` の最後に `12`、`0` が表示される |
| 要素の抽出 | `cmd/collections/main.go` の `filterAtLeast` | 下記の呼び出しを追加して実行する |
| 名前の変更とテスト | `internal/task/task.go` の `Rename`、`task_test.go` の `TestRename` | `go test ./internal/task -run TestRename -v` |
| 期限切れ | `cmd/concurrency/main.go` の `main` 末尾 | `go run ./cmd/concurrency` の最後に `deadline: true` が表示される |
| 稼働確認 | `internal/httpapi/handler.go` のルート登録、`handler_test.go` の `TestHealth` | `go test ./internal/httpapi -run TestHealth -v` |

`filterAtLeast` は関数の定義だけが収録されています。`cmd/collections/main.go` の既存の `main` 関数内に次の行を追加し、`go run ./cmd/collections` で `[30 20]` が表示されることを確かめてください。

```go
fmt.Println(filterAtLeast([]int{10, 30, 20}, 20))
```

</details>

## Todo HTTP APIを動かす

APIの処理は、`internal/task`（Todoの検証と保存）、`internal/httpapi`（リクエストとレスポンス）、`cmd/api`（サーバーの組み立てと起動）に分かれています。

`go-lab/` 内で起動します。

```sh
go run ./cmd/api
```

`listening on http://127.0.0.1:8080` が表示されたら、サーバーを動かしたまま**別のターミナル**を開き、以下のコマンドを順番に実行します。

### Bashで試す

```sh
curl -i http://127.0.0.1:8080/tasks
curl -i -H 'Content-Type: application/json' --data '{"title":"Learn Go"}' http://127.0.0.1:8080/tasks
curl -i http://127.0.0.1:8080/tasks
curl -i http://127.0.0.1:8080/health
```

### PowerShellで試す

```powershell
Invoke-RestMethod -Uri 'http://127.0.0.1:8080/tasks'
Invoke-RestMethod -Uri 'http://127.0.0.1:8080/tasks' -Method Post -ContentType 'application/json' -Body '{"title":"Learn Go"}'
Invoke-RestMethod -Uri 'http://127.0.0.1:8080/tasks'
Invoke-RestMethod -Uri 'http://127.0.0.1:8080/health'
```

### レスポンスを確認する

起動直後に上の4つのリクエストを送ると、次の結果になります。`curl -i` はHTTPヘッダーも表示します。`Invoke-RestMethod` はJSONをオブジェクトに変換するため、表示形式は異なります。

| 順序 | リクエスト | ステータス | レスポンス本文（JSON） |
| --- | --- | --- | --- |
| 1 | `GET /tasks` | 200 | `[]` |
| 2 | `POST /tasks` | 201 | `{"id":1,"title":"Learn Go","done":false}` |
| 3 | `GET /tasks` | 200 | `[{"id":1,"title":"Learn Go","done":false}]` |
| 4 | `GET /health` | 200 | `{"status":"ok"}` |

不正な入力も確認できます。Bashで次のリクエストを送ると、HTTP 400と本文 `title is required` が返り、Todoは追加されません。

```sh
curl -i -H 'Content-Type: application/json' --data '{"title":"   "}' http://127.0.0.1:8080/tasks
```

終了するときはサーバーを起動したターミナルで **Ctrl+C** を押します。保存先はメモリなので、再起動するとTodoは消えます。追加を繰り返すとIDが増えるため、表と同じ結果を確認するときはサーバーを再起動してください。

## 全体の確認と理解度チェック

コードを変更したら、`go-lab/` 内で整形・ビルド・テスト・静的解析を実行します。

```sh
go fmt ./...
go build ./...
go test ./...
go vet ./...
```

`go fmt` はソースの書式を整え、`go build` はコンパイルを確認します。`go test` はPythonの `unittest` や `pytest` のようにテストを実行し、`go vet` は誤りの疑いがある記述を検出します。HTTPテストは `httptest` を使うため、APIサーバーの起動は不要です。

対応する環境では、データ競合も検査します。

```sh
go test -race ./...
go run -race ./cmd/concurrency
```

`-race` には対応OS・アーキテクチャ、cgoの有効化、Cコンパイラーが必要です。[公式の実行要件](https://go.dev/doc/articles/race_detector#Requirements)を参照してください。利用できない環境では、通常のテストと `go vet` まで実行します。競合検出は実行された処理が対象なので、成功してもすべての実行経路の安全性を保証するものではありません。

第7章では、次の項目をコードや実行結果と合わせて説明できるか確認します。

- [ ] PythonのリストのスライスとGoのsliceで、要素の変更が元データに及ぼす影響を説明できる。
- [ ] `Task` のコピーを変更した場合と、ポインタを通じて元の値を変更した場合の違いを説明できる。
- [ ] `err != nil` と `errors.Is` の使い分け、`defer` の実行タイミングを説明できる。
- [ ] goroutineの完了待ち、channelを閉じる役割、キャンセル、Mutexが必要な理由を説明できる。
- [ ] APIでTodoを追加・取得でき、不正な入力が保存されないことをテストで確認できる。
- [ ] 保存処理が失敗した場合のHTTP 500、内部エラー情報の非公開、`GET /health`、並行アクセスのテストを確認できる。

次の課題として、Todoを完了させるHTTPエンドポイント、DBへの永続化、graceful shutdown、同時実行数の制限に取り組めます。`Task.Complete` 自体は実装済みですが、これらの発展課題の完成実装は含まれていません。

## 困ったとき

| 症状 | 確認・対処 |
| --- | --- |
| `go` が見つからない | Goをインストールし、ターミナルを開き直して `go version` を確認する |
| `go.mod file not found`、`directory prefix ... does not contain main module` | リポジトリ直下から `cd go-lab` で移動し、`go env GOMOD` のパスを確認する |
| `redeclared` が出る | 演習の関数や `main` が重複していないか確認する。解答例は追記せず、既存の実装を書き換える |
| `declared and not used`、`imported and not used` | 未使用のローカル変数やimportを削除するか、実際に使用する |
| API起動時に `address already in use` が出る | 同じポートで起動済みのサーバーがないか確認する。自分が起動したサーバーは、そのターミナルでCtrl+Cを押して終了する |
| APIに接続できない | サーバー側の起動エラーを確認する。同じマシンから `http://127.0.0.1:8080` にアクセスする |
| `-race` を実行できない | [実行要件](https://go.dev/doc/articles/race_detector#Requirements)を確認し、通常の `go test ./...` と `go vet ./...` を先に実行する |
| `error obtaining VCS status` が出る | Gitの作業ディレクトリを確認する。一時的にVCS情報の埋め込みを省くには `go build -buildvcs=false ./...` を使う |
