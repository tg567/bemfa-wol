GOOS=linux GOARCH=amd64 go build -o wol_amd64_linux -ldflags "-s -w" .
GOOS=linux GOARCH=arm64 go build -o wol_arm64_linux -ldflags "-s -w" .
GOOS=linux GOARCH=arm go build -o wol_arm_linux -ldflags "-s -w" .
GOOS=windows GOARCH=amd64 go build -o wol_amd64_windows.exe -ldflags "-s -w" .