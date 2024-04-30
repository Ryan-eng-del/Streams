package proxy

import (
	"fmt"

	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

// Codec returns a proxying grpc.Codec with the default protobuf codec as parent.
//
// See CodecWithParent.
// 构建输出函数
func Codec() encoding.Codec {
	return CodecWithParent(&protoCodec{})
}
// CodecWithParent returns a proxying grpc.Codec with a user provided codec as parent.
//
// This codec is *crucial* to the functioning of the proxy. It allows the proxy server to be oblivious
// to the schema of the forwarded messages. It basically treats a gRPC message frame as raw bytes.
// However, if the server handler, or the client caller are not proxy-internal functions it will fall back
// to trying to decode the message using a fallback codec.
func CodecWithParent(fallback encoding.Codec) encoding.Codec {
	return &rawCodec{fallback}
}

func init() {
	encoding.RegisterCodec(Codec())
}


// This our custom codec
type rawCodec struct {
	parentCodec encoding.Codec
}

type frame struct {
	payload []byte
}

func (c *rawCodec) Marshal(v interface{}) ([]byte, error) {
	// log.Printf("进来 Marshal 值为 %+v, 类型为 %+v", v, reflect.TypeOf(v))
	out, ok := v.(*frame)
	// log.Println("进来 Marshal", out, ok)
	// log.Printf("进来 Marshal v:%+v, ok: %t", v, ok)
	if !ok {
		return c.parentCodec.Marshal(v)
	}
	return out.payload, nil
}

func (c *rawCodec) Unmarshal(data []byte, v interface{}) error {
	// log.Printf("进来 Unmarshal 值为 data: %+v(type:%+v), v: %+v(type:%+v)", data, reflect.TypeOf(data), v, reflect.TypeOf(v))
	dst, ok := v.(*frame)
	// log.Printf("进来 Unmarshal 2 dst: %+v, ok: %t, v : %+v", data, ok, v)
	if !ok {
		return c.parentCodec.Unmarshal(data, v)
	}
	dst.payload = data
	// log.Println(dst.payload, "payload", "进来 Unmarshal 3")
	return nil
}

func (c *rawCodec) String() string {
	return fmt.Sprintf("proxy>%s", c.parentCodec.Name())
}

func (c *rawCodec) Name() string {
	return "proxy-proto"
}


// This is a callback codec implementation with google protoc
// protoCodec is a Codec implementation with protobuf. It is the default rawCodec for gRPC.
// 构建proto解码器
type protoCodec struct{}
// Name is the name registered for the proto compressor.
const Name = "proxy-proto"

// codec is a Codec implementation with protobuf. It is the default codec for gRPC.
func (protoCodec) Marshal(v any) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (protoCodec) Unmarshal(data []byte, v any) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

// func messageV2Of(v any) proto.Message {
// 	switch v := v.(type) {
// 	case protoadapt.MessageV1:
// 		return protoadapt.MessageV2Of(v)
// 	case protoadapt.MessageV2:
// 		return v
// 	}

// 	return nil
// }

func (protoCodec) Name() string {
	return Name
}

func (c *protoCodec) String() string {
	return fmt.Sprintf("proxy>%s", c.Name())
}