package avro

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/linkedin/goavro.v1"
	"io/ioutil"
	"reflect"
	"testing"
)

/* Round-trip some primitive values through our serializer and goavro to verify */
const fixtureJson = `
[
{"IntField": {"small": 1, "min":-2147483647, "max":2147483647}, "LongField": {"small": 2, "min": 9223372036854775807, "max": -9223372036854775807}, "FloatField": {"small": 3.4, "verysmall": 3.402823e-38, "large": 3.402823e+38}, "DoubleField": {"small": 5.6, "verysmall": 2.2250738585072014e-308}, "StringField": {"short": "789", "longer": "a slightly longer string"}, "BoolField": {"true": true, "false":false}, "BytesField": {"small": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n", "longer": "VGhpcyBpcyBhIG11Y2ggbG9uZ2VyIHRlc3Qgc3RyaW5nIGxvbmcgbG9uZw=="}},
{"IntField": {}, "LongField": {}, "FloatField": {}, "DoubleField": {}, "StringField": {}, "BoolField": {"true": true}, "BytesField": {}}
]
`

func TestMapFixture(t *testing.T) {
	fixtures := make([]MapTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	schemaJson, err := ioutil.ReadFile("maps.avsc")
	if err != nil {
		t.Fatal(err)
	}
	codec, err := goavro.NewCodec(string(schemaJson))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		if err != nil {
			t.Fatal(err)
		}
		datum, err := codec.Decode(&buf)
		if err != nil {
			t.Fatal(err)
		}
		record := datum.(*goavro.Record)
		value := reflect.ValueOf(f)
		for i := 0; i < value.NumField(); i++ {
			fieldName := value.Type().Field(i).Name
			avroVal, err := record.Get(fieldName)
			if err != nil {
				t.Fatal(err)
			}
			avroMap := avroVal.(map[string]interface{})
			if len(avroMap) != value.Field(i).Len() {
				t.Fatalf("Got %v keys from goavro but expected %v", len(avroMap), value.Field(i).Len())
			}
			for _, k := range value.Field(i).MapKeys() {
				keyString := k.Interface().(string)
				avroMapVal := avroMap[keyString]
				structMapVal := value.Field(i).MapIndex(k).Interface()
				if !reflect.DeepEqual(avroMapVal, structMapVal) {
					t.Fatalf("Field %v key %v not equal: %v != %v", fieldName, k, avroMapVal, structMapVal)
				}
			}
		}
	}
}

func BenchmarkMapRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		record := MapTestRecord{map[string]int32{"value1": 1, "value2": 2, "value3": 3}, map[string]int64{"value1": 1, "value2": 2, "value3": 3}, map[string]float64{"value1": 1, "value2": 2, "value3": 3}, map[string]string{"value1": "12345", "value2": "67890", "value3": "abcdefg"}, map[string]float32{"value1": 1, "value2": 2, "value3": 3}, map[string]bool{"true": true, "false": false}, map[string][]byte{"value1": {1, 2, 3, 4}, "value2": {100, 200, 255}}}
		err := record.Serialize(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapGoavro(b *testing.B) {
	schemaJson, err := ioutil.ReadFile("maps.avsc")
	if err != nil {
		b.Fatal(err)
	}
	codec, err := goavro.NewCodec(string(schemaJson))
	if err != nil {
		b.Fatal(err)
	}
	someRecord, err := goavro.NewRecord(goavro.RecordSchema(string(schemaJson)))
	if err != nil {
		b.Fatal(err)
	}
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		someRecord.Set("IntField", map[string]interface{}{"value1": int32(1), "value2": int32(2), "value3": int32(3)})
		someRecord.Set("LongField", map[string]interface{}{"value1": int64(1), "value2": int32(2), "value3": int32(3)})
		someRecord.Set("FloatField", map[string]interface{}{"value1": float32(1), "value2": float32(2), "value3": float32(3)})
		someRecord.Set("DoubleField", map[string]interface{}{"value1": float32(1), "value2": float32(2), "value3": float32(3)})
		someRecord.Set("StringField", map[string]interface{}{"value1": "12345", "value2": "67890", "value3": "abcdefg"})
		someRecord.Set("BoolField", map[string]interface{}{"true": true, "false": false})
		someRecord.Set("BytesField", map[string]interface{}{"value1": []byte{1, 2, 3, 4}, "value2": []byte{100, 200, 255}})

		codec.Encode(buf, someRecord)
	}

}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]MapTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		if err != nil {
			t.Fatal(err)
		}
		datum, err := DeserializeMapTestRecord(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *datum, f)
	}
}
