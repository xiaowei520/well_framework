## well_framework

### 编译
上线系统最终上线的是output目录，因此编译后的可执行文件以及配置文件都会放到output目录中
执行 `./build.sh`  进行编译
看到done表示编译成功

### 运行
* 开发环境：make run
* 线上环境：`cd output && ./bin/well_framework --rpcconf=yourpath1 --logconf=yourpath2

线上因为没有Makefile所以没法make run，但线上环境的运行方式也能在线下使用

### 日志
运行日志保存在`output/log`中

### 代码结构
* main.go: 用来进行整个进程的初始化，包括读取配置文件和启动logger，然后监听端口分对http提供服务
* vendor: 依赖包，一般不要动。vendor.json不能手动改，必须使用[govendor](https://github.com/kardianos/govendor)命令行工具来更新和操作依赖
* conf: conf中没有代码，都是配置文件，打包时会copy到output/conf
* server: server是个adaptor层：
  * http: http主要是配置路由，以及每个路由的回调方法的实现
* logic: 业务逻辑
* models: 数据模型

> 需要注意的是，在http中，并不实现对应接口的业务逻辑，做的工作只是填充数据结构，调用logic层的函数。目的是实现 解耦业务逻辑实现函数 和 请求方法 的解耦 
