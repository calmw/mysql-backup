## 镜像构建

``` shell
# build 
docker build -t backup:0.1.1 . 
# tag
docker tag backup:0.1.1 calmw/backup:0.1.1
# push
docker push calmw/backup:0.1.1
```