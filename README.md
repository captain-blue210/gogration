# gogration

## 概要
- TimeCampからtogglにデータ移行するために個人用に作ったツールです。
- Goで開発しています。

## 使い方
- gogration直下にTimeCampから落としてきたCSVデータを`detailed.csv`として配置
- `go run main.go`を実行
- gogration直下にtoggl.csvができるので、それをtoggl上からアップロードする

## 注意事項
- 個人用に作ったので、データ移行までは確認していますが細かいバグ等までは考慮しておりません。もし問題点ありましたらプルリクお願いします。
