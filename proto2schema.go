package proto2schema

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/emicklei/proto"
)

// 字段类型映射表

var defaultFieldTypeMap = map[string]string{
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

	fieldTypeMap := GenFieldTypeMap(definitions, false)
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

				fieldType, err := getFieldType(fieldTypeMap, field.Type)
				if err != nil {
					log.Fatal("NormalField", err)
				}
				builder.WriteString(fieldType)
				if field.Repeated {
					builder.WriteString("]")
				}
				builder.WriteString(lineBreak)

			case *proto.MapField:
				builder.WriteString("    ")
				builder.WriteString(field.Name)
				builder.WriteString(": {")
				keyType, err := getFieldType(fieldTypeMap, field.KeyType)
				if err != nil {
					log.Fatal("MapField KeyType", err)
				}
				builder.WriteString(keyType)
				builder.WriteString(":")
				fieldType, err := getFieldType(fieldTypeMap, field.Type)
				if err != nil {
					log.Fatal("MapField fieldType", err)
				}
				builder.WriteString(fieldType)
				builder.WriteString("}")
				builder.WriteString(lineBreak)
			}
		}

		builder.WriteString(lineBreak)
	}

	return builder.String()
}

// GenFieldTypeMap
func GenFieldTypeMap(definitions *proto.Proto, useIntegersForNumbers bool) map[string]string {
	fieldTypeMap := make(map[string]string)
	for key, value := range defaultFieldTypeMap {
		fieldTypeMap[key] = value
	}

	for _, definition := range definitions.Elements {
		switch visitee := definition.(type) {
		case *proto.Message:
			fieldTypeMap[visitee.Name] = visitee.Name
		case *proto.Enum:
			var builder strings.Builder
			elementsLen := len(visitee.Elements) - 1
			for i, e := range visitee.Elements {
				v, ok := e.(*proto.EnumField)
				if !ok {
					continue
				}

				value := v.Name
				if useIntegersForNumbers {
					value = strconv.Itoa(v.Integer)
				}

				builder.WriteString(fmt.Sprintf(`"%v"`, value))
				if elementsLen > i {
					builder.WriteString(` | `)
				}

				fieldTypeMap[v.Name] = v.Name
			}
			fieldTypeMap[visitee.Name] = builder.String()
		}
	}

	return fieldTypeMap
}

// getFieldType 获取字段类型，如果是默认类型则返回映射后的类型，否则原样返回
func getFieldType(fieldTypeMap map[string]string, fieldType string) (string, error) {
	value, ok := fieldTypeMap[fieldType]
	if !ok {
		return "", fmt.Errorf(`this "%v" is not currently supported`, fieldType)
	}

	return value, nil
}
