module example.com/grpc-protobuf

go 1.22

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
