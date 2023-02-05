proto-chat: 
	protoc --go_out=. --go_opt=paths=source_relative \
    ./chat/chat.proto

run:
	go run ./main.go