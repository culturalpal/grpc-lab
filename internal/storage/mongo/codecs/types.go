package codecs

import (
	"errors"
	bookv1 "github.com/ppal31/grpc-lab/generated/book/v1"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"reflect"
)

var bookType = reflect.TypeOf(&bookv1.Book{})

func CreateObject(val reflect.Value) (interface{}, error) {
	switch val.Type() {
	case bookType:
		return val.Interface().(*bookv1.Book), nil
	default:
		return nil, errors.New("value type is not suppoerted")
	}
}

func EncoderMap() map[reflect.Type]bsoncodec.ValueEncoder {
	em := make(map[reflect.Type]bsoncodec.ValueEncoder)
	em[bookType] = &ProtoEncoder{}
	return em
}
