package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tronprotocol/grpc-gateway/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/quiknode-labs/tron-proto/proto/protocol"
)

type tokenAuth struct {
	token string
}

// GetRequestMetadata adds the token to the request metadata
func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{"x-token": t.token}, nil
}

// RequireTransportSecurity indicates whether the credentials require transport security
func (tokenAuth) RequireTransportSecurity() bool {
	return true
}

func main() {
	// QuikNode gRPC endpoint configuration
	target := "my-endpoint-name.tron-mainnet.quiknode.pro:50051"
	token := "YOUR_TOKEN"
	creds := tokenAuth{token: token}

	// Define gRPC options
	grpcOpts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(creds),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	}

	// Establish connection
	grpcConn, err := grpc.Dial(target, grpcOpts...)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcConn.Close()

	// Create Wallet client
	client := api.NewWalletClient(grpcConn)

	// Call GetNowBlock
	block, err := client.GetNowBlock(context.Background(), &api.EmptyMessage{})
	if err != nil {
		log.Fatalf("Failed to get latest block: %v", err)
	}

	blockJSON, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		log.Fatalf("Failed to format block data: %v", err)
	}

	// Print the formatted JSON block data
	fmt.Println("Latest Block Information:")
	fmt.Println(string(blockJSON))
}