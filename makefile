.PHONY: genproto
genproto:
	protoc -I. --go_out=plugins=grpc:./pd/fight ./pd/fight/fight.proto
	protoc -I. --go_out=plugins=grpc:./pd/auth ./pd/auth/auth.proto
