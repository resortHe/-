package serializer

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

func WriteProtobufToJSONFile(message proto.Message, filename string) error {
	data, err := ProtobufToJSON(message)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	err2 := ioutil.WriteFile(filename, []byte(data), 0644)
	if err2 != nil {
		return fmt.Errorf(err2.Error())
	}
	return nil

}
func WriteProtobufToBinaryFile(message proto.Message, filname string) error {
	b, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	err2 := ioutil.WriteFile(filname, b, 0644)
	if err2 != nil {
		return fmt.Errorf(err2.Error())
	}
	return nil
}
func ReadProtobufToBinaryFile(filname string, message proto.Message) error {
	b, err := ioutil.ReadFile(filname)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	err2 := proto.Unmarshal(b, message)
	if err2 != nil {
		return fmt.Errorf(err2.Error())
	}
	return nil
}
