package utiles

import (
	"github.com/bilibili/kratos/pkg/log"
	"google.golang.org/grpc"
)

func GrpcFromPythonClient() (*grpc.ClientConn, error) {
	rpcHost := GetEnv("GRPC_PYTHON_HOST", "127.0.0.1")
	rpcPort := GetEnv("GRPC_PYTHON_PORT", "50051")
	conn, err := grpc.Dial(rpcHost+":"+rpcPort, grpc.WithInsecure())
	if err != nil {
		log.Error("did not connect: %v", err)
	}
	return conn, err
}
