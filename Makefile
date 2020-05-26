build:
	@rm -rf target
	@go mod tidy
	@go build -o ./target/dargo -v github.com/itering/subscan/cmd
