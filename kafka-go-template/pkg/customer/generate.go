package customer

// auto-generated code.  Use this whenever changes are made to the
// avro schema and sink.go

//go:generate $GOPATH/bin/gogen-avro -containers --package customer . customer.avsc

//go:generate mkdir -p mocks
//go:generate $GOPATH/bin/mockgen -source=sink.go -destination=mocks/sink_mock.go
