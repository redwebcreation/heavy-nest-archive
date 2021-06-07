export GOOS=linux

printf "Enter the version : "

read -r version

echo "$version" > .version

go build -ldflags="-s -w"

gh release create "$version" ./hez -t "Release $version"