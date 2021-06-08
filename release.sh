export GOOS=linux

printf "Enter the version : "

read -r version

echo "$version" > static/version

go generate ./...

go build -ldflags="-s -w"

gh release create "$version" ./hez -t "Release $version" -d --notes "Automated release of $version."