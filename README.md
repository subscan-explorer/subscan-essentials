# subscan-end
 
### run
    go run main.go -conf ../configs

### substrate
     go run main.go -conf ../configs start substrate
 
### protoc
    cd libs/substrate/protos
    protoc -I codec_protos  codec_protos/rpc.proto --go_out=plugins=grpc:codec_protos

