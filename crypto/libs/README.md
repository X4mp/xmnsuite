# Linux:
go build -o ../../ui/libs/xmncrypto.so -buildmode=c-shared lib.go

# Mac:
go build -o ../../ui/libs/xmncrypto.dylib -buildmode=c-shared lib.go
