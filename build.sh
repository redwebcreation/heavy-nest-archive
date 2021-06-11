go build -ldflags="-w -s -X github.com/redwebcreation/hez2/globals.Version=$(git describe --tags)" -gcflags=all="-l"
