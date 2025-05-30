#!/bin/sh
env GOOS=linux GOARCH=arm64 go build -o ./build/gpsd_viewer_linux_arm64
env GOOS=linux GOARCH=amd64 go build -o ./build/gpsd_viewer_linux_amd64
env GOOS=darwin GOARCH=arm64 go build -o ./build/gpsd_viewer_darwin_arm64
env GOOS=darwin GOARCH=amd64 go build -o ./build/gpsd_viewer_darwin_amd64
env GOOS=windows GOARCH=amd64 go build -o ./build/gpsd_viewer_windows_amd64.exe
