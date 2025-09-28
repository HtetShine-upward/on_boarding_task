package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	memo "example.com/memo/api/memo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// firstNonEmpty は、引数の中で最初に空でない文字列を返すユーティリティ関数
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

func main() {
	// CLIのグローバルフラグを定義（例: memo-cli -addr localhost:50051 create ...）
	root := flag.NewFlagSet("memo-cli", flag.ExitOnError)
	addrFlag := root.String("addr", "", "gRPC server address (e.g. localhost:50051). Also reads MEMO_ADDR")
	deadlineFlag := root.Duration("deadline", 4*time.Second, "per-RPC deadline (e.g. 4s, 1500ms)")
	_ = root.Parse(os.Args[1:]) // 引数をパース

	args := root.Args()
	if len(args) < 1 {
		os.Exit(2) // コマンドが指定されていない場合は終了
	}
	cmd := args[0]      // 実行するコマンド（create, get, list, delete）
	cmdArgs := args[1:] // コマンドに渡す引数

	// gRPCサーバーのアドレスを取得（優先順位: フラグ > 環境変数 > デフォルト）
	addr := firstNonEmpty(*addrFlag, os.Getenv("MEMO_ADDR"), "localhost:50051")
	rpcDeadline := *deadlineFlag // RPCのタイムアウト時間

	// conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// gRPC接続を確立（5秒以内に接続できなければ失敗）
	dialCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		dialCtx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // TLSなしで接続
		grpc.WithBlock(), // 接続が確立するまでブロック
	)
	if err != nil {
		log.Fatalf("dial %s failed: %v", addr, err)
	}
	defer conn.Close()

	// MemoServiceのクライアントを作成
	client := memo.NewMemoServiceClient(conn)

	// コマンドに応じて処理を分岐
	switch cmd {
	case "create":
		createCmd(client, rpcDeadline, cmdArgs)
	case "get":
		getCmd(client, rpcDeadline, cmdArgs)
	case "list":
		listCmd(client, rpcDeadline, cmdArgs)
	case "delete":
		deleteCmd(client, rpcDeadline, cmdArgs)
	default:
		fmt.Println("Unknown command:", cmd)
		os.Exit(2)
	}
}

// createCmd はメモの作成処理を行う
func createCmd(client memo.MemoServiceClient, deadline time.Duration, argv []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	title := fs.String("title", "", "Title of the memo (1-100 chars)")
	content := fs.String("content", "", "Content of the memo (0-2000 chars)")
	_ = fs.Parse(argv)

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	// CreateMemo RPCを呼び出す
	res, err := client.CreateMemo(ctx, &memo.CreateMemoRequest{
		Title:   *title,
		Content: *content,
	})
	if err != nil {
		handleRPCError("CreateMemo", err)
	}
	fmt.Printf("Created memo: ID=%s\n", res.GetMemo().GetId())
}

// getCmd は指定されたIDのメモを取得する
func getCmd(client memo.MemoServiceClient, deadline time.Duration, argv []string) {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	id := fs.String("id", "", "Memo ID")
	_ = fs.Parse(argv)

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	res, err := client.GetMemo(ctx, &memo.GetMemoRequest{Id: *id})
	if err != nil {
		handleRPCError("GetMemo", err)
	}
	m := res.GetMemo()
	fmt.Printf("Memo: ID=%s Title=%q Content=%q CreatedAt=%d\n",
		m.GetId(), m.GetTitle(), m.GetContent(), m.GetCreatedAt())
}

// listCmd はすべてのメモを一覧表示する
func listCmd(client memo.MemoServiceClient, deadline time.Duration, argv []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	_ = fs.Parse(argv)

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	res, err := client.ListMemos(ctx, &memo.ListMemosRequest{})
	if err != nil {
		handleRPCError("ListMemos", err)
	}
	for _, m := range res.GetMemos() {
		fmt.Printf("ID=%s title=%q created_at=%d\n", m.GetId(), m.GetTitle(), m.GetCreatedAt())
	}
}

// deleteCmd は指定されたIDのメモを削除する
func deleteCmd(client memo.MemoServiceClient, deadline time.Duration, argv []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	id := fs.String("id", "", "Memo ID")
	_ = fs.Parse(argv)

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	_, err := client.DeleteMemo(ctx, &memo.DeleteMemoRequest{Id: *id})
	if err != nil {
		handleRPCError("DeleteMemo", err)
	}
	fmt.Println("Memo deleted:", *id)
}

// handleRPCError は gRPC エラーを整形して表示する
func handleRPCError(op string, err error) {
	if st, ok := status.FromError(err); ok {
		log.Fatalf("%s failed: code=%s msg=%s", op, st.Code(), st.Message())
	}
	log.Fatalf("%s failed: %v", op, err)
}
