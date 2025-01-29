package main

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"

	rd "math/rand"
	"net"
	"time"

	pb "github.com/temaniarpit27/aggchain_prover_mock/aggchain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedAggchainProofServiceServer
}

func (s *server) GenerateAggchainProof(ctx context.Context, req *pb.GenerateAggchainProofRequest) (*pb.GenerateAggchainProofResponse, error) {
	rnd := rd.New(rd.NewSource(time.Now().UnixNano()))
	if rnd.Float64() < 0.2 { // 20% chance of timeout
		log.Println("Simulating timeout...")
		time.Sleep(130 * time.Second)
	}

	startBlock := req.StartBlock
	endBlock, _ := rand.Int(rand.Reader, big.NewInt(int64(req.MaxEndBlock-startBlock+1)))
	endBlock.Add(endBlock, big.NewInt(int64(startBlock)))

	proof := make([]byte, 32)
	rand.Read(proof)

	return &pb.GenerateAggchainProofResponse{
		AggchainProof: proof,
		StartBlock:    startBlock,
		EndBlock:      endBlock.Uint64(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAggchainProofServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
