# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=WsLogger
BINARY_FOLDER_WIN=$(PWD)/dist/win
BINARY_FOLDER_MAC=$(PWD)/dist/mac
BINARY_FOLDER_UBUNTU=$(PWD)/dist/ubuntu
BINARY_FOLDER_ANDROID=$(PWD)/dist/android

dist:build-win build-mac build-ubuntu build-android
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_FOLDER_WIN)/$(BINARY_NAME)
	rm -f $(BINARY_FOLDER_MAC)/$(BINARY_NAME)
	rm -f $(BINARY_FOLDER_UBUNTU)/$(BINARY_NAME)
	rm -f $(BINARY_FOLDER_ANDROID)/$(BINARY_NAME)

# Cross compilation
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_FOLDER_WIN)/$(BINARY_NAME) -v

build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_FOLDER_MAC)/$(BINARY_NAME) -v

build-ubuntu:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_FOLDER_UBUNTU)/$(BINARY_NAME) -v

build-android:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(BINARY_FOLDER_ANDROID)/$(BINARY_NAME) -v