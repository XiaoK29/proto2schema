package proto2schema

import (
	"log"
	"os"
	"strings"

	"github.com/emicklei/proto"
)

var fileldType = map[string]string{
	"uint32":              "int",
	"uint64":              "int",
	"int32":               "int",
	"int64":               "int",
	"sint32":              "int",
	"sint64":              "int",
	"string":              "str",
	"google.protobuf.Any": "any",
}

func Proto2schema(path string) string {
	// 读取.proto文件
	f, err := os.Open("./test.proto")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	defer f.Close()

	// 解析.proto文件
	parser := proto.NewParser(f)
	definitions, err := parser.Parse()
	if err != nil {
		log.Fatal("Error parsing proto file:", err)
	}

	var builder strings.Builder
	for _, definition := range definitions.Elements {
		message, ok := definition.(*proto.Message)
		if !ok {
			continue
		}

		builder.WriteString(`schema `)
		builder.WriteString(message.Name)
		builder.WriteString(`:`)
		builder.WriteString("\n")
		for _, element := range message.Elements {
			if field, ok := element.(*proto.NormalField); ok {
				builder.WriteString(" ")
				builder.WriteString(" ")
				builder.WriteString(" ")
				builder.WriteString(" ")
				builder.WriteString(field.Name)
				if field.Optional {
					builder.WriteString(`?`)
				}

				builder.WriteString(`: `)
				if field.Repeated {
					builder.WriteString(`[`)
				}

				builder.WriteString(isDefultFileldType(field.Type))
				if field.Repeated {
					builder.WriteString(`]`)
				}

				builder.WriteString("\n")
				continue
			}

			if field, ok := element.(*proto.MapField); ok {
				builder.WriteString(" ")
				builder.WriteString(" ")
				builder.WriteString(" ")
				builder.WriteString(" ")
				builder.WriteString(field.Name)
				builder.WriteString(`: `)
				builder.WriteString(`{`)
				builder.WriteString(isDefultFileldType(field.KeyType))
				builder.WriteString(`:`)
				builder.WriteString(isDefultFileldType(field.Type))
				builder.WriteString(`}`)
				builder.WriteString("\n")
			}
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

func isDefultFileldType(str string) string {
	fieldType, ok := fileldType[str]
	if ok {
		return fieldType
	}

	return str
}
