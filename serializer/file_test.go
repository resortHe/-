package serializer_test

import (
	"proto_demo/pb"
	"proto_demo/sample"
	"proto_demo/serializer"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()
	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"
	laptop1 := sample.NewLaptop()
	err := serializer.WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err2 := serializer.ReadProtobufToBinaryFile(binaryFile, laptop2)
	require.NoError(t, err2)
	require.True(t, proto.Equal(laptop1, laptop2))
	err3 := serializer.WriteProtobufToJSONFile(laptop1, jsonFile)
	require.NoError(t, err3)
}
