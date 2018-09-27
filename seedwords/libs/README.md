# Linux:
go build -o ../../ui/libs/xmnseedwords.so -buildmode=c-shared lib.go

# Mac:
go build -o ../../ui/libs/xmnseedwords.dylib -buildmode=c-shared lib.go
