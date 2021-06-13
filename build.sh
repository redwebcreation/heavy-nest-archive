go build -ldflags="-w -s -X github.com/redwebcreation/hez/globals.Version=$(git describe --tags)" -gcflags=all="-l"
