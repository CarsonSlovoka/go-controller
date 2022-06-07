@echo off
cd ..
go env -w CC=C:\mingw64\bin\gcc
go env -w CXX=C:\mingw64\bin\g++
go build -ldflags "-s -w" -o ..\bin\go-controller.exe
:: move go-controller.exe ..\bin
start ..\bin
echo restore & echo.
go env -w CC=gcc
go env -w CXX=g++
