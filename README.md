# Easy-Online-Judge

项目描述：go在线练习算法系统的后端

## 技术栈

- Go：用于编写后端服务。
- Redis：用于缓存和数据存储。
- GORM：用于数据库操作。
- Swagger：用于API文档生成和测试。
- Gin：用于构建Web应用程序。
- go-uuid：用于生成UUID。

## 先决条件

在运行项目之前，请确保您的开发环境满足以下要求：

- Go 1.16或更高版本
- Redis服务器
- mysql
- 以及go语言基础

## 安装和运行

1. 克隆项目到本地：

   ```
   colne https://github.com/gonggonggnog/gin-oj.git
   cd gin-oj
   ```

2. 安装依赖：

   ```
   go mod tidy
   ```

3. 配置项目：

   复制示例配置文件，并根据您的需求进行配置：

   ```
   cp .env.example .env
   ```

   编辑`.env`文件，配置Redis和其他必要的环境变量。

4. 运行项目：

   ```
   go run main.go
   ```

   项目将在默认端口上运行（通常是8080）。

## API文档

您可以使用Swagger生成的API文档来查看和测试API。访问以下链接：

```
http://localhost:8080/swagger/index.html
```
