这是个练手的项目，当时好奇fan  墙软件是怎么实现的，而公司正在使用的工具是一个同事用Go语言[gotunnel](https://github.com/shahuwang/gotunnel)写的，
Go语言比较好懂，所以就找来代码来读了。读的过程中，感觉好像懂了，但是又好像不是很懂，于是乎自己临摹了一下他的代码，不过我去掉了抽象层，
从client -> tunnel -> server, 发现不行，于是变成了 client -> tube -> tunnel -> tube -> server 的形式，总算能跑通了，不过只能执行一次，
因为我发现这个抽象结构，无法处理连接的关闭，导致了长连接一直在等待。于是我发现，还是需要再增加两层，一层用于buffer，一层用做link，同
一个tunnel可以接收发送多个连接（以linkid做标示）的数据，数据存储在link.buffer里面。发送的数据需要包含一个header，指明数据的长度和linkid。
linkid 为 0 的数据，其实是指令，表明是关闭了，创建了等等一系列tcp连接事件。

搞明白之后，突然意兴阑珊，因为就基本上是抄gotunnel的代码了。 所以，就暂时先搁置了，主要是练手。不过我想把gotunnel里面的一些东西写成文章，
比如一个通道被多个客户使用，tcp数据怎么处理呢；比如客户端和服务器端如何进行加密认证，确保可信服务的呢；

把别人如何实现一个东西看懂，还是挺有成就感的。但是单纯抄别人的代码真的是好没有成就感的事情，所以就这样吧。
