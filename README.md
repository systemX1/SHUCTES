# 计算思维实训项目 课程评价系统后端部分

#### 介绍

SHUCTES(Shanghai University Course and Teaching Evaluation System)

go+mysql

#### 主要第三方框架、库

gin, logrus, viper, go-sql-driver/mysql

#### 接口文档

[SHUCTES Aliyun (getpostman.com)](https://documenter.getpostman.com/view/13925655/TzCFir9R)

#### 编译

go build -o ctes.exe ./src/Main.go

#### 线上部署

更改conf/config.yml

```shell
#local
docker build . -f DockerFile -t ctes
docker tag ctes NAME[:TAG]
docker push NAME[:TAG]

#remote
docker pull NAME[:TAG]

docker run \
-p 8000:8000 \
-d \
--add-host=host.docker.internal:host-gateway \
NAME[:TAG]
```

