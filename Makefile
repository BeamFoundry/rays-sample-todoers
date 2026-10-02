.PHONY: all
all: clean nocgo

.PHONY: tool
	go tool stereotyper
	
.PHONY: clean
clean:
	rm -f todolist 

.PHONY: cgo
cgo: tool
	go build -ldflags "-s -w" -o todolist cmd/main.go

.PHONY: nocgo
nocgo: tool
	CGO_ENABLED=0 go build -ldflags "-s -w" -o todolist cmd/main.go
