module example.com/grpc-protobuf

go 1.25.0

// A real version of this module, once .pb.go files are generated for
// real with protoc, would also require:
//
//   require (
//       google.golang.org/grpc v1.6x.x
//       google.golang.org/protobuf v1.3x.x
//   )
//
// Those are deliberately omitted here -- see code/README.md. This
// go.mod exists so the directory has the right shape, not because
// this code actually builds in this sandbox.

require (
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/grpc v1.83.2 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)
