# gocache

企业级软件开发大作业，用go实现一个分布式缓存

**github地址**：[go-cache](https://github.com/zkyoma/go-cache)

## 项目启动

go 版本：go 1.18

下载项目后，cd 到 go-cache 目录下，执行 sh run.sh，控制台如下图所示则启动成功：

![image-20230527160255904](https://blog-htz.oss-cn-hangzhou.aliyuncs.com/gihub/start.png)

默认运行的结构如下图所示：

![image-20230527160359871](https://blog-htz.oss-cn-hangzhou.aliyuncs.com/gihub/QQ%E6%88%AA%E5%9B%BE20230527162219.png)
启动了三个 Cache Server，端口分别为 8001、8002、8003，其中 API Server 用来接收外界请求，DB Server 是模拟订单假数据，只有如上图所示的这些数据。

可以用`curl http://localhost:9999/api?key=Tom` 测试，代表根据 key：Tom去查询缓存，第一次是缓存未命中，则去DB Server 拿到对应的值 “630”返回，并添加到缓存中，再请求一次，则缓存命中。对应的两次控制台日志如下：

![image-20230527161459881](https://blog-htz.oss-cn-hangzhou.aliyuncs.com/gihub/1.png)

![image-20230527161520836](https://blog-htz.oss-cn-hangzhou.aliyuncs.com/gihub/2.png)

**gocache是一个缓存库，也就是说不是一个完整的软件，需要自己实现main函数（main.go）**，若想搭建自己的缓存服务，则可以修改现有的 main.go 的逻辑。

## 线上地址

按照上述方式，我在自己的服务器上已经运行，以`http://42.192.227.166:9999/api?key=Tom` 这种形式可以访问缓存服务。可以修改 key 的值来获取对应的value。**注意：只能更改 key 的值其他格式不能改变，否则访问不到**
