# Image merge

要发旅游攻略给女票，下载回来的攻略是多张图片的形式，写了这个程序垂直拼接图片，方便查看。

## 使用

```shell
go get -u github.com/dujigui/imgmerge
imgmerge
```

### 用例
```shell
1. imgmerge -od ~/Desktop/ -i ~/Desktop/imgs
2. imgmerge -of ~/Desktop/imgmerge.png ~/Desktop/1.jpg ~/Desktop/2.jpg
3. imgmerge -od ~/Desktop -m min -i ~/Desktop/imgs
4. imgmerge -od ~/Desktop -i ~/Desktop/imgs -s 1.5
```

## 编译

### 要求
1. go1.11.2
2. 启动 go module

### 步骤
```shell
git clone git@github.com:dujigui/imgmerge.git
cd imgmerge
go mod tidy
go build .
```

## Todo

1. 增加test*.go
2. 优化程序出错时的提示