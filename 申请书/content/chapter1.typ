#import "../template.typ": *

= 主要目标

eKuiper 中有以下几个重要的基本概念：规则，流，表。其中主要的与本项⽬相关的概念是规则。

一个规则代表了⼀个流处理流程，定义了从将数据输⼊流的数据源到各种处理逻辑，再到将数据输⼊到外部系统的动作。

#box(
figure(caption: [#text("eKuiper 架构")])[#image("../figures/ekuiperArch.png")],

)

在 eKuiper 中，可以部署多个流式规则来进行数据的处理。但是当 eKuiper 中部署多条规则时，我们无法通过目前的 go profile 来快速得知每条规则的 CPU 用量，以及规则与规则之间的 CPU 用量对比，从而来快速定位 eKuiper 资源消耗问题。

本项目的主要目标便是解决上述问题，通过以下方法完成：

- 开发兼容的、接口友好的、可独立的 profile 库
	- 兼容：引入这一库不对项目本身造成冲突
	- 接口友好：提供的 profile 接口易使用、易理解
	- 可独立：该库可以独立成一个小的开源库，为其他项目提供更加高级的 profile 功能
- 对于 eKuiper 项目中流式规则的代码使用 `pprof.SetGoroutineLabels` 埋点采集性能数据，使用开发的 profile 库和编写的性能数据处理函数提取关键数据，prometheus 开源库采集提取出的关键数据。用 prometheus提供的可视化工具或者其他开源可视化工具进行单条规则的性能观察和多条规则之间的性能对比


