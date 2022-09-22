.ONESHELL:

default:
	cd ./frontend; pnpm build;
	go build -o otp-panel *.go

clean:
	rm -rf ./public
	rm ./otp-panel*

multiarch:
	GOOS=linux GOARCH=amd64 go build -o otp-panel-linux-amd64 *.go
	GOOS=darwing GOARCH=arm64 go build -o otp-panel-darwin-aarch64 *.go
	GOOS=windows GOARCH=amd64 go build -o otp-panel-windows-amd64.exe *.go