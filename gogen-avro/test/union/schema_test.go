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
{"UnionField":{"Int":1, "UnionType":0}},
{"UnionField":{"Long":2, "UnionType":1}},
{"UnionField":{"Float":3.4, "UnionType":2}},
{"UnionField":{"Double":5.6, "UnionType":3}},
{"UnionField":{"String":"testString", "UnionType":4}},
{"UnionField":{"Bool":true, "UnionType":5}},
{"UnionField":{"Bytes":"VGhpcyBpcyBhIHRlc3Qgc3RyaW5n", "UnionType":6}},
{"UnionField":{"UnionType":7}}
]
`

func TestPrimitiveUnionFixture(t *testing.T) {
	fixtures := make([]PrimitiveUnionTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	schemaJson, err := ioutil.ReadFile("union.avsc")
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
			fieldUnionIndex := int(value.Field(i).FieldByName("UnionType").Int())
			structVal := value.Field(i).Field(fieldUnionIndex).Interface()
			avroVal, err := record.Get(fieldName)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(structVal, avroVal) {
				t.Fatalf("Field %v not equal: %v != %v", fieldName, structVal, avroVal)
			}
		}
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]PrimitiveUnionTestRecord, 0)
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
		datum, err := DeserializePrimitiveUnionTestRecord(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *datum, f)
	}
}

func TestDefault(t *testing.T) {
	record := NewPrimitiveUnionTestRecord()
	assert.Equal(t, record.UnionField.Int, int32(1234))
}
