echo package main > version.go
echo const >> version.go

type ..\VERSION.txt >> version.go
go build -o ../sample.exe
