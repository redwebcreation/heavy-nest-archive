export GOOS=linux

printf "Enter the version : "

read -r version

mkdir static

echo "$version" > static/version

go build -ldflags="-s -w"

gh release create "$version" ./hez -t "Release $version" -d --notes "Automated release of $version."