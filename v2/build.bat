@echo off
go env -w CC="C:\mingw64\bin\gcc
go env -w CXX="C:\mingw64\bin\g++
go build
echo restore & echo.
go env -w CC=gcc
go env -w CXX=g++
