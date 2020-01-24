APP=mpcifier

.PHONY: build
build: clean
	go build -o ${APP} .

.PHONY: run
run:
	go run -race .

.PHONY: clean
clean:
	go clean
	rm ${APP} | true