# 计算思维实训项目 课程评价系统 后端部分

**SHUCTES(Shanghai University Course and Teaching Evaluation System)**

### 主要架构

go + mysql

### 主要第三方库，框架

gin, viper, logrus, go-sql-driver/mysql

### 手动编译

go build -o ctes.exe ./src/Main.go

### 线上部署

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

