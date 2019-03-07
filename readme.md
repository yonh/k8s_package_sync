# k8s包同步工具

这个是在安装k8s的时候因为墙的原因写的一个用于同步k8s包的工具，目前只同步`kubernetes-xenial`里面的包，也许以后增加其他的包吧

软件包会自动下载到当前目录下的mirror目录


### 使用
```
git clone https://github.com/yonh/k8s_package_sync.git
cd k8s_package_sync
go build .

./k8s_package_sync
```

docker run -it --rm -v $PWD/mirror:/usr/share/nginx/html -p 80:80 nginx


