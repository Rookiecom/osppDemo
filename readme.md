说明：demo的功能基本上都在申请书中写明，此处写每个文件的内容概括，方便理解

- demo/task 文件夹存放的是任务本身和CPU数据处理函数
  - prime.go 是并行筛法求素数
  - mergeSort.go 是并行归并排序
  - profile.go 主要是注册消费者（提供数据处理函数）
  - graph.go 主要是用plot库对多次采集的CPU数据进行可视化分析（申请书中有图）
- demo/cpu_profile 文件夹是profile库的实现
  - api.go是库向外部暴露的接口
  - profile.go是CPU数据采集者的结构定义和接口定义实现
  - consume.go是消费者的接口定义和接口定义实现
  - collect.go是使用消费者来实现对`runtime/pprof`库的`StartCPUProfile`和`StopCPUProfile`两个接口的替代
  - web.go是对`net/http/pprof`的兼容处理

我已经将demo打包成linux(./demo/demo)和windows(./demo/demo.ext)的可执行文件，可以不用下载库依赖。