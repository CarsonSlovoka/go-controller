# go-controller

餵入json檔案，即可按照此腳本來執行動作{滑鼠移動、觸發熱鍵...}

## Build & USAGE

```
git clone https://github.com/CarsonSlovoka/go-controller.git
cd go-controller/v2
go build -o go-controller.exe
go-controller.exe -file="<test.json>"
```

You need to feed a file (JSON format) to the executable.

for example

> go-controller.exe -file="[./examples/hello-world.json](v2/examples/hello-world.json)"

Once the program is running, you can type `help` in the console to see other commands that you can use.

## 給開發人員

本套件主要倚靠[go-vgo/robotgo]所完成

- [使用Robotgo開發的前置作業](https://github.com/CarsonSlovoka/robotgo/blob/13a5c80/README_zh_TW.md)


[go-vgo/robotgo]: https://github.com/go-vgo/robotgo
