cd proto

protoc --go_out=plugins=grpc:. prod.proto
protoc --go_out=plugins=grpc:. --validate_out=lang=go:. models.proto
protoc --go_out=plugins=grpc:. orders.proto
protoc --go_out=plugins=grpc:. users.proto

protoc --grpc-gateway_out=logtostderr=true:. prod.proto
protoc --grpc-gateway_out=logtostderr=true:. orders.proto
protoc --grpc-gateway_out=logtostderr=true:. users.proto