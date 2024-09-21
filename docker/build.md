## 镜像构建

``` shell
# build 
docker build -t backup:0.1.0 . 
# tag
docker tag backup:0.1.0 calmw/backup:0.1.0
# push
docker push calmw/backup:0.1.0
```