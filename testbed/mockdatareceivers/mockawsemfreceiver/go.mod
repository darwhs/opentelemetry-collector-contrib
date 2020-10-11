module github.com/open-telemetry/opentelemetry-collector-contrib/testbed/mockdatareceivers/mockawsxrayreceiver

go 1.14

require (
	github.com/gorilla/mux v1.8.0
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/collector v0.9.1-0.20200903224024-3eb3b664a832
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.16.0
)

replace go.opentelemetry.io/collector => ../../../../opentelemetry-collector
