# Installation

## Client
- [Windows](https://github.com/Tradou/deskor/blob/windows/README.md)
- [Linux](https://github.com/Tradou/deskor/blob/linux/README.md)

## Server
[See here](https://github.com/Tradou/deskor/blob/server/README.md)

## Setup development environment (notes)
### Prerequisites

- Docker is required to build client platform from a custom linux image
- Go is required to build server side

### Build custom docker image
```shell
docker build -t fyne-custom -f build/Dockerfile .
```

### Build client platform
```shell
# Need to delete server.go in order to build client platform
rm server.go

# Platform should be the target platform (eg: linux)
fyne-cross PLATFORM -image fyne-custom
```

### Build server
```shell
go build -o bin/server server.go
```

### Build assets
Each assets need to be bundle to be used in the **client** version. 
```shell
fyne bundle -o ./assets/bundle/MY_ASSET.go ./assets/MY_FOLDER/MY_ASSET
```
Make sure to edit the bundle assets to make it exportable