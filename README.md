# Image merge

要发旅游攻略给女票，下载回来的攻略是多张图片的形式，写了这个程序垂直拼接图片，方便查看。

## 使用

```shell
go get -u github.com/pheynix/imgmerge
imgmerge -h
```

## 编译

### 要求
1. go1.11.2
2. 启动 go module

### 步骤
```shell
git clone git@github.com:pheynix/imgmerge.git
cd imgmerge
go mod tidy
go build .
```

## Todo

1. 增加test*.go
2. 优化程序出错时的提示