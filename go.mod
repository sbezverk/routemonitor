module github.com/sbezverk/routemonitor

go 1.14

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/sbezverk/gobmp v0.0.0-20200623233057-23947e63d1b6
	google.golang.org/grpc v1.30.0
)

replace github.com/sbezverk/gobmp => ../gobmp
