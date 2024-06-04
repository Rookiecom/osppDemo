= 实现

== profile 库

设计的profile库架构图如下：

#box(
  figure(caption: [#text("profile 库架构")])[#image("../figures/架构图.png")],

)

当程序启动时，通过profile库提供的`StartCPUProfiler`启动CPU数据采集者协程，该协程定时采集程序的profile数据（可配置采集间隔时间和每次采集的时间窗口），每次采集之后会将数据通过channel发送给注册的消费者。

项目通过引入profile库，添加消费者实例来使用profile库的功能。添加消费者实例必须通过设置回调函数来设置消费者处理数据的方式。添加一个新的消费者实例并使用profile库提供的`StartConsume`启动消费者首先会在CPU数据采集者处注册，然后启动消费者协程等待CPU数据采集者。成功收到数据后，消费者会调用回调函数来处理数据。（对应图2的下部）。

profile库中的CPU数据采集者是通过`runtime/pprof`库提供的`StopCPUProfile`和`StartCPUProfile`来实现指定时间段性能采集。为了保证兼容性，防止引入该库的项目调用`runtime/pprof`库提供的`StopCPUProfile`和`StartCPUProfile`而产生的混乱，profile库应该提供替代的接口（名称不变），这两个接口的实现仍是通过消费者来实现（对应图2的上部）。

另外，go还提供了web查看程序性能的方法，可以通过下面的代码来使用这一功能：

#box(
```go
import _ "net/http/pprof"

func main() {
	go func() {
		log.Println(http.ListenAndServe(addr, mux))
	}()
  // 程序逻辑
}
```
)

大致的实现原理是引入`net/http/pprof`意味着程序启动会先执行库中的init函数，该库的init函数中的`http.HandleFunc("/debug/pprof/profile", Profile)`设置的网络接口的处理函数同样使用了`runtime/pprof`库提供的`StopCPUProfile`和`StartCPUProfile`。所以需要在profile库提供同样的功能，但是需要将`/debug/pprof/profile`的处理函数用上文实现的替代接口实现。

可以在线查看profile库的#link("https://github.com/Rookiecom/osppDemo", "demo")实现。

== eKuiper rule CPU用量统计

在eKuiper中，每条规则有规则id，在go中的上下文（contex）中记录规则id，这样做能够让每个协程都能通过上下文获取到启动自身的规则id。

在整个规则的代码逻辑中，通过`pprof.SetGoroutineLabels(ctx context.Context)`埋点设置来采集每条规则的性能数据，以下是埋点的代码：
#box(
```go
ctx = pprof.WithLabels(ctx, pprof.Labels("rule", ruleID))
pprof.SetGoroutineLabels(ctx)
```
)

eKuiper rule CPU用量统计首先引入2.1实现的profile库，然后添加消费者实例，为消费者实例提供数据处理回调函数。对于eKuiper项目来说，处理数据要关注标签为rule的数据，并且需要对每条规则的CPU用量单独统计。对于统计出的数据，可以定义对应的prometheus收集器并注册到prometheus中，这样便能够在prometheus中查看每条规则的CPU用量情况，并且能够通过编写prometheus的查询语句PromQl来对比多条规则之间的CPU用量。

具体的实现逻辑、埋点的位置考量需要在项目完成期间完成。但是我编写了一个小的demo来测试我的profile库，这一demo会不断的并发执行两个任务——并行筛法求素数和并行归并排序。设置两个任务都是并行的目的是更加贴近项目（一条规则的执行流程包含多个协程）。demo中的并行归并排序埋点代码如下：

// #box(  clip: true,
```go
func ParallelMergeSort(ctx context.Context, array []int, enableProfile bool) {
	if enableProfile {
		pprof.SetGoroutineLabels(ctx)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
  // split array code ellipsis
	go func() {
		ParallelMergeSort(ctx, leftArray, enableProfile)
		wg.Done()
	}()
	go func() {
		ParallelMergeSort(ctx, rightArray, enableProfile)
		wg.Done()
	}()
	wg.Wait()
  // merged array code ellipsis
}
```

#box(
"       可以看到每个协程都有调用pprof.SetGoroutineLabels(ctx)来埋点统计性能信息。对应并行筛法求素数也做了同样的处理。demo调用了2.1中的profile库并按照规则编写了数据处理函数统计了两个任务的CPU量。单个任务的数据查看和两个任务之间的数据对比，我使用了go中的plot库，使用两个任务的CPU占用量(百分比)生成折线图，并且可以通过访问网络接口实时监控折线图的变化，效果图如下："
)


#align(center, figure(caption: [#text("demo 监控图")])[#image("../figures/监控图.png", height: 43%)])

#box(
"       另外，我还做了负载测试。"
)


测试说明：以并行筛法求解 2~10000 素数为例，N 倍负载代表同时启动 N 个这样的任务（这 N 个任务将被打上同一个标签）。在测试中，将同时启动 20 倍负载和 200 倍负载，并分别统计总 CPU 用量，测试结果如下：

#align(center, figure(caption: [#text("demo 负载测试图")])[#image("../figures/负载测试.png", height: 25%)])


#box(
"       可以看到：200 倍负载的 CPU 占用量几乎是 20 倍负载的 CPU 占用量的 10倍，符合预期。"
)
