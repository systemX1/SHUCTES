# 计算思维实训项目 课程评价系统 后端部分

**SHUCTES(Shanghai University Course and Teaching Evaluation System)**

### 主要架构

go + mysql

### 主要第三方库，框架

gin, viper, logrus, go-sql-driver/mysql

### 接口文档

https://documenter.getpostman.com/view/13925655/TzCFir9R 

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

### TO DO

日志分割, swagger, 自动重启, 表单验证

权限系统, 评论点赞点踩