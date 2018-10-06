# logging
这是一个日志模块，简单封装了`logrus`，提供的功能：
- `logrus`初始化
- 捕获`panic`输出调用栈并退出程序
- `context`支持
- `gin`框架的访问日志和异常捕获
- 输出日志到`fluentd`

### 为什么会有这个东西
1. 我们的程序要输出日志
1. 我们希望将所有程序的日志集中起来方便查看和管理
1. 我们选择了`elasticsearch`, `kibana` 和 `fluentd` 的组合来解决这个问题
1. 为了方便日志处理，我们需要输出格式化(`JSON`)的日志
1. 刚好`logrus`提供了这个功能，并且用的人也挺多
1. 为了输出`JSON`格式的日志需要初始化`logrus`的设置
1. 很久之后我们又需要把日志直接输出到`fluentd`(某云厂商功能不足导致的)

### 日志是如何收集的
简单来说就是格式化的日志输出到`fluentd`，`fluentd`再将日志输入到`elasticsearch`，然后就可以通过`kibana`查询了。

但是程序日志输出到`fluentd`却有两种方式，分别适用于阿里云部署的程序和京东云部署的程序。

方式一：
```text
    service --> stdout --> docker log driver --> fluentd --> elasticsearch
```
这是最理想的方式，程序不需要关心日志要输出到哪里，而是通过容器来解决。阿里云部署的程序采用的是这种方式。

方式二：
```text
    service --> fluentd(jclould) --> fluentd(aliyun) --> elasticsearch
```
这种方式是为了解决京东云无法使用`log driver`的问题。

### 怎么使用

```go
package main

import (
	"errors"
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zaoshu/hardcore/logging"
	"github.com/zaoshu/hardcore/logging/fluentd"
)

func test(ctx context.Context)  {
    // 从 context 中获取 logger
    logger := logging.FromContext(ctx)
    logger.Info("get logger from context")
}

func main() {
    defer logging.GodCatchPanic() // 捕获 panic 并退出程序
    
	// 初始化 logrus
	c := logging.Config{
		Level: "debug",
		Fluent: &fluentd.Config{ // 如果设置了将会把日志输出到 fluentd
		    Host: "127.0.0.1",
		    Port: 24424,
		    Tag:  "debug.test", // tag 的第一部分是运行环境, 第二部分是服务名称
		    ContainerName: "aaaa", // 京东云的容器无法从内部获取 container id, 使用 container name 替代 id，通过环境变量设置
		},
	}
	logging.InitFromConfig(c)

	logrus.Debug("debug")
	logrus.Info("info")
	logrus.Warn("warn")
	logrus.Error("this is error")
	
	test(context.Background())
}
```
