# git-nfs

make git repo as a nfs server file storage

## mount option

```shell
mount -o "port=<port>,mountport=<port>,intr,noresvport,nolock,noacl" -t nfs localhost:/ /mount
```

## TODO

* optimize: 历史commit数量可设置
* optimize: file diff when push
* optimize: 追踪运行时间<https://github.com/jaegertracing/jaeger>
* optimize: 日志输出优化<https://github.com/uber-go/zap>
* optimize: OpenTelemetry <https://openTelemetry.io>
