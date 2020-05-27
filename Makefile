build:
	@rm -rf target
	@go mod tidy
	@go build -o ./target/subscan -v github.com/itering/subscan/cmd
