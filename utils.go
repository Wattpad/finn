package finn

import (
	"fmt"
	"log"

	"github.com/ugorji/go/codec"
)

// LogError logs error messages to standard error
func LogError(err error) {
	log.Printf("Error: %s", err)
}

// LogInfo logs informational messages to standard out
func LogInfo(info string) {
	fmt.Printf("%s\n", info)
}

// LogInfoColour is LogInfo except it colours the info green
func LogInfoColour(info string) {
	fmt.Printf("\x1B[0;92m%s\x1B[0m\n", info)
}

// Unpack decodes a job into a worker object
func Unpack(value []byte, template GenericWorker) (GenericWorker, error) {
	var handle codec.MsgpackHandle
	decoder := codec.NewDecoderBytes(value, &handle)
	worker := template.NewInstance()
	return worker, decoder.Decode(worker)
}

// Pack encodes a worker object into a job to be put into a queue
func Pack(worker GenericWorker) []byte {
	var handle codec.MsgpackHandle
	var packed []byte
	encoder := codec.NewEncoderBytes(&packed, &handle)
	encoder.Encode(worker)
	return packed
}
