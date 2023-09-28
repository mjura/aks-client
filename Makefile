.PHONY: client
client:
	go build -o bin/aks-client main.go

.PHONY: clean
clean:
	rm -rf build bin dist
