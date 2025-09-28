セットアップ手順

1.リポジトリのクローンと初期化
git clone <your-repo-url>
cd memo-service
go mod tidy

2.proto からコード生成
# プラグイン導入（未導入の場合）
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# PATH 設定（必要なら）
export PATH="$PATH:$(go env GOPATH)/bin"

# コード生成
protoc \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  api/memo.proto

3.ローカル開発
go run ./server/cmd/server

# CLI クライアント使用例
# メモ作成
go run ./client/cmd/cli create --title "First" --content "hello"
# メモ一覧
go run ./client/cmd/cli list
# メモ取得
go run ./client/cmd/cli get --id <ID>
# メモ削除
go run ./client/cmd/cli delete --id <ID>

4.Dockerビルド／起動
イメージビルド
docker build -f Dockerfile.server -t memo-server:local .
docker build -f Dockerfile.client -t memo-client:local .

サーバー起動
docker run --rm -p 50051:50051 memo-server:local
docker run --rm --network host memo-client:local list

5. Kubernetes デプロイ(kind)
kind create cluster
kind load docker-image memo-server:local

kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# ポートフォワード
kubectl port-forward svc/memo-server 50051:50051

CLI実行
go run ./client/cmd/cli list

