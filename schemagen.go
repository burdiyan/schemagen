package schemagen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Landoop/schema-registry"
	"github.com/asaskevich/govalidator"
	"gopkg.in/alanctgardner/gogen-avro.v5/generator"
	"gopkg.in/alanctgardner/gogen-avro.v5/types"
)

const (
	kindAvro = "Avro"
)

// SchemaConfig describes schemas to be downloaded.
type SchemaConfig struct {
	Subject string `yaml:"subject" valid:"required"`
	Version string `yaml:"version" valid:"required"`
	Package string `yaml:"package" valid:"required"`
}

// Config is the global application config.
type Config struct {
	Kind      string         `yaml:"kind" valid:"required"`
	Registry  string         `yaml:"registry" valid:"required"`
	Compile   bool           `yaml:"compile" valid:"required"`
	OutputDir string         `yaml:"outputDir" valid:"required"`
	Schemas   []SchemaConfig `yaml:"schemas" valid:"required"`
}

// Run uses the Config to download schemas and to compile them.
func Run(ctx context.Context, cfg Config) error {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		return err
	}

	switch cfg.Kind {
	case kindAvro:
		return generateAvro(ctx, cfg)
	default:
		return fmt.Errorf("kind %q is not supported", cfg.Kind)
	}
}

func generateAvro(ctx context.Context, cfg Config) error {
	client, err := schemaregistry.NewClient(cfg.Registry)
	if err != nil {
		return err
	}

	if err := os.Mkdir(cfg.OutputDir, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	for _, s := range cfg.Schemas {
		var schema string

		if s.Version == "latest" {
			sch, err := client.GetLatestSchema(s.Subject)
			if err != nil {
				return err
			}
			schema = sch.Schema
		} else {
			v, err := strconv.Atoi(s.Version)
			if err != nil {
				return fmt.Errorf("version %q is not valid: %v", s.Version, err)
			}

			sch, err := client.GetSchemaBySubject(s.Subject, v)
			if err != nil {
				return err
			}

			schema = sch.Schema
		}

		if err := compileAvroSchema(schema, s, cfg.OutputDir); err != nil {
			return err
		}
	}

	return nil
}

func compileAvroSchema(schema string, cfg SchemaConfig, out string) error {
	pkg := generator.NewPackage(cfg.Package)
	namespace := types.NewNamespace()

	_, err := namespace.TypeForSchema([]byte(schema))
	if err != nil {
		return err
	}

	for _, v := range namespace.Definitions {
		rec, ok := v.(*types.RecordDefinition)
		if !ok {
			continue
		}

		filename := generator.ToSnake(rec.Name()) + ".go"

		generateGoka(filename, rec, pkg)
	}

	if err := namespace.AddToPackage(pkg, codegenComment([]string{cfg.Package + ".avsc"}), false); err != nil {
		return err
	}

	var b bytes.Buffer
	if err := json.Indent(&b, []byte(schema), "", "    "); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(out, cfg.Package+".avsc"), b.Bytes(), 0755); err != nil {
		return err
	}

	target := path.Join(out, cfg.Package)

	if err := os.Mkdir(target, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	return pkg.WriteFiles(path.Join(out, cfg.Package))
}

// codegenComment generates a comment informing readers they are looking at
// generated code and lists the source avro files used to generate the code
//
// invariant: sources > 0
func codegenComment(sources []string) string {
	const fileComment = `// Code generated by github.com/burdiyan/schemagen. DO NOT EDIT.
/*
 * %s
 */`
	var sourceBlock []string
	if len(sources) == 1 {
		sourceBlock = append(sourceBlock, "SOURCE:")
	} else {
		sourceBlock = append(sourceBlock, "SOURCES:")
	}

	for _, source := range sources {
		_, fName := filepath.Split(source)
		sourceBlock = append(sourceBlock, fmt.Sprintf(" *     %s", fName))
	}

	return fmt.Sprintf(fileComment, strings.Join(sourceBlock, "\n"))
}
