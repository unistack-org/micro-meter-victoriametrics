module github.com/micro/go-plugins/wrapper/monitoring/victoriametrics/v2

go 1.13

require (
	github.com/VictoriaMetrics/metrics v1.9.3
	github.com/micro/go-micro/v2 v2.9.1-0.20200716123506-3627e47f04eb
	github.com/stretchr/testify v1.4.0
)

replace github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
