# Online Judge
基于gin+grpc+gorm+consul的在线判题系统

## 模块划分
* 题目服务 - 后端主服务
* DB服务：mysql和redis
* 判题服务

## 目录build-env
`builder-env.dockerfile`：基于debian:trixie-slim搭建一个包含gcc、g++、cmake、make的基础镜像。

`builder.dockerfile`：基于上诉的基础镜像，添加服务编译所需要的依赖并打包成新的镜像。

`create_builder.sh`：执行脚本构建具有最终编译环境的容器。

`build.sh`：进入容器，可以开始编译，参数(目标路径)

## 编译步骤
根据Makefile提示将依赖库编译成链接文件放在指定位置，并导入头文件路径。
1. muduo
2. protobuf

依赖准备好后，在judge-service目录下执行make。
