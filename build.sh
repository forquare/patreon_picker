#!/bin/sh

echo "Don't forget to run 'go mod tidy'."
if [ -f picker.txz ]; then
	rm picker.txz
fi

env GOOS=freebsd GOARCH=amd64 go build
tar cJf picker.txz patreon_picker creds.json cookie.txt templates static
rm patreon_picker
