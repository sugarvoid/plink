
clean() {
    if [ -f plink ]; then
        rm gplink
        echo "old build is removed"
    fi
}

windows() {
    CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" .
}

linux() {
    go build -ldflags "-s -w" .
    ./gplink
}

case "$1" in
    "linux")
        clean
        linux
        ;;
    "windows")
        windows
        ;;
    *)  
        clean
        linux
        ;;
esac