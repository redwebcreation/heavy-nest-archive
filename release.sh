export GOOS=linux

printf "Enter the version : "

read -r version

go build -ldflags="-w -s -X github.com/redwebcreation/hez/globals.Version=$version" -gcflags=all="-l"

gh release create "$version" ./hez -t "Release $version" -d --notes "Automated release of $version."