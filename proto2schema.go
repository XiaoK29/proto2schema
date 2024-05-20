package proto2schema

import (
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/emicklei/proto"
)

// 字段类型映射表

var fieldTypeMap = map[string]string{
	"uint32":              "int",
	"uint64":              "int",
	"int32":               "int",
	"int64":               "int",
	"sint32":              "int",
	"sint64":              "int",
	"string":              "str",
	"google.protobuf.Any": "any",
	"bool":                "bool",
	"float":               "float",
	"double":              "float",
}

// Proto2schema 将.proto文件转换为schema字符串
func Proto2schema(path string) string {
	// 读取.proto文件
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	defer f.Close()

	lineBreak := "\n"
	if runtime.GOOS == "windows" {
		lineBreak = "\r\n"
	}

	// 解析.proto文件
	parser := proto.NewParser(f)
	definitions, err := parser.Parse()
	if err != nil {
		log.Fatal("Error parsing proto file:", err)
	}

	// 构建schema字符串
	var builder strings.Builder
	for _, definition := range definitions.Elements {
		message, ok := definition.(*proto.Message)
		if !ok {
			continue
		}

		builder.WriteString("schema ")
		builder.WriteString(message.Name)
		builder.WriteString(":")
		builder.WriteString(lineBreak)

		for _, element := range message.Elements {
			switch field := element.(type) {
			case *proto.NormalField:
				builder.WriteString("    ")
				builder.WriteString(field.Name)
				if field.Optional {
					builder.WriteString("?")
				}
				builder.WriteString(": ")

				if field.Repeated {
					builder.WriteString("[")
				}
				builder.WriteString(getFieldType(field.Type))
				if field.Repeated {
					builder.WriteString("]")
				}
				builder.WriteString(lineBreak)

			case *proto.MapField:
				builder.WriteString("    ")
				builder.WriteString(field.Name)
				builder.WriteString(": {")
				builder.WriteString(getFieldType(field.KeyType))
				builder.WriteString(":")
				builder.WriteString(getFieldType(field.Type))
				builder.WriteString("}")
				builder.WriteString(lineBreak)
			}
		}

		builder.WriteString(lineBreak)
	}

	return builder.String()
}

// getFieldType 获取字段类型，如果是默认类型则返回映射后的类型，否则原样返回
func getFieldType(fieldType string) string {
	if mappedType, ok := fieldTypeMap[fieldType]; ok {
		return mappedType
	}
	return fieldType
}
