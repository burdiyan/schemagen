package schemagen

import (
	"bytes"
	"text/template"

	"gopkg.in/alanctgardner/gogen-avro.v5/generator"
	"gopkg.in/alanctgardner/gogen-avro.v5/types"
)

var (
	gokaCodecTpl   = template.Must(template.New("codec").Parse(`type {{.CodecType}} struct{}`))
	gokaEncoderTpl = template.Must(template.New("encoder").Parse(`
func ({{.CodecGoType}}) Encode(value interface{}) ([]byte, error) {
	v := value.({{.BaseGoType}})
	
	var b bytes.Buffer
	
	if err := v.Serialize(&b); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
`))

	gokaDecoderTpl = template.Must(template.New("decoder").Parse(`
func ({{.CodecGoType}}) Decode(data []byte) (interface{}, error) {
	return Deserialize{{.BaseType}}(bytes.NewReader(data))
}	
`))
)

func generateGoka(filename string, r *types.RecordDefinition, pkg *generator.Package) {
	params := struct {
		BaseType    string
		BaseGoType  string
		CodecType   string
		CodecGoType string
	}{
		BaseType:    r.Name(),
		BaseGoType:  "*" + r.Name(),
		CodecType:   r.Name() + "Codec",
		CodecGoType: "*" + r.Name() + "Codec",
	}

	pkg.AddImport(filename, "bytes")
	pkg.AddStruct(filename, params.CodecType, renderTemplate(gokaCodecTpl, params))
	pkg.AddFunction(filename, params.CodecGoType, "Encode", renderTemplate(gokaEncoderTpl, params))
	pkg.AddFunction(filename, params.CodecGoType, "Decode", renderTemplate(gokaDecoderTpl, params))
}

func renderTemplate(t *template.Template, params interface{}) string {
	var b bytes.Buffer
	if err := t.Execute(&b, params); err != nil {
		// Templates should never fail.
		panic(err)
	}

	return b.String()
}
