judge: src/judge.go
	go build -o bin/judge src/judge.go

server: src/server.go
	go build -o bin/server src/server.go 
