if [ -z "$1" ]; then
  echo "Usage: release [version]"
  exit 1
fi

echo "Releasing version: $1"

GOOS=linux go build -ldflags="-w -s -X github.com/redwebcreation/hez/globals.Version=$1" -gcflags=all="-l"

gh release create "$1" ./hez -t "Release $1" -d --notes "Automated release of $1." $2
