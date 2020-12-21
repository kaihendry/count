bin/local/main: main.go
	mkdir -p bin/local
	CGO_ENABLED=0 go build -o bin/local/main main.go

localdev:
	func start --verbose

deploy:
	func azure functionapp publish counttesting
