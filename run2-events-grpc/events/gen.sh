protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf \
    ./events.proto
