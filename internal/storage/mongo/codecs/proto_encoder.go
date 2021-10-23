package codecs

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
)

type ProtoEncoder struct{}

func (b *ProtoEncoder) EncodeValue(ec bsoncodec.EncodeContext, writer bsonrw.ValueWriter, val reflect.Value) error {
	var (
		result map[string]interface{}
		obj    interface{}
		err    error
	)

	if obj, err = CreateObject(val); err != nil {
		return err
	}

	js, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(js, &result)
	result["_id"] = result["id"]
	delete(result, "id")
	if err != nil {
		return err
	}
	codec := bsoncodec.MapCodec{}
	return codec.EncodeValue(ec, writer, reflect.ValueOf(result))
}
