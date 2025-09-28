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

// firstNonEmpty returns the first non-empty trimmed string.
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

func main() {
	// Global flags: memo-cli [GLOBAL FLAGS] <command> [FLAGS...]
	root := flag.NewFlagSet("memo-cli", flag.ExitOnError)
	addrFlag := root.String("addr", "", "gRPC server address (e.g. localhost:50051). Also reads MEMO_ADDR")
	deadlineFlag := root.Duration("deadline", 4*time.Second, "per-RPC deadline (e.g. 4s, 1500ms)")
	_ = root.Parse(os.Args[1:])

	args := root.Args()
	if len(args) < 1 {
		os.Exit(2)
	}
	cmd := args[0]
	cmdArgs := args[1:]

	addr := firstNonEmpty(*addrFlag, os.Getenv("MEMO_ADDR"), "localhost:50051")
	rpcDeadline := *deadlineFlag

	// Establish connection (block until ready or timeout)
	dialCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		dialCtx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("dial %s failed: %v", addr, err)
	}
	defer conn.Close()

	client := memo.NewMemoServiceClient(conn)

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

func createCmd(client memo.MemoServiceClient, deadline time.Duration, argv []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	title := fs.String("title", "", "Title of the memo (1-100 chars)")
	content := fs.String("content", "", "Content of the memo (0-2000 chars)")
	_ = fs.Parse(argv)

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	res, err := client.CreateMemo(ctx, &memo.CreateMemoRequest{
		Title:   *title,
		Content: *content,
	})
	if err != nil {
		handleRPCError("CreateMemo", err)
	}
	fmt.Printf("Created memo: ID=%s\n", res.GetMemo().GetId())
}

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

func handleRPCError(op string, err error) {
	if st, ok := status.FromError(err); ok {
		log.Fatalf("%s failed: code=%s msg=%s", op, st.Code(), st.Message())
	}
	log.Fatalf("%s failed: %v", op, err)
}
