package converter

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
	"unicode"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

func ExprToStr(e ast.Expr) (string, error) {
	switch exp := e.(type) {
	case *ast.SelectorExpr:
		x, err := ExprToStr(exp.X)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s.%s", x, exp.Sel.Name), nil
	case *ast.Ident:
		return exp.String(), nil
	case *ast.StarExpr:
		x, err := ExprToStr(exp.X)
		if err != nil {
			return "", err
		}
		return "*" + x, nil
	case *ast.ArrayType:
		elt, err := ExprToStr(exp.Elt)
		if err != nil {
			return "", err
		}
		return "[]" + elt, nil
	default:
		return "", fmt.Errorf("invalid type %T - %v", exp, exp)
	}
}

func ToProtoType(typ string) string {
	if strings.HasPrefix(typ, "*") {
		return ToProtoType(typ[1:])
	}
	if strings.HasPrefix(typ, "[]") && typ != "[]byte" {
		return "repeated " + ToProtoType(typ[2:])
	}
	switch typ {
	case "json.RawMessage", "[]byte":
		return "bytes"
	case "sql.NullBool", "pgtype.Bool":
		return "google.protobuf.BoolValue"
	case "sql.NullInt32", "pgtype.Int4", "pgtype.Int2":
		return "google.protobuf.Int32Value"
	case "pgtype.Uint32":
		return "google.protobuf.UInt32Value"
	case "int":
		return "int64"
	case "int16":
		return "int32"
	case "uint16":
		return "uint32"
	case "sql.NullInt64", "pgtype.Int8":
		return "google.protobuf.Int64Value"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "pgtype.Float4":
		return "google.protobuf.FloatValue"
	case "sql.NullFloat64", "pgtype.Float8":
		return "google.protobuf.DoubleValue"
	case "sql.NullString", "pgtype.Text", "pgtype.UUID":
		return "google.protobuf.StringValue"
	case "sql.NullTime", "time.Time", "pgtype.Date", "pgtype.Timestamp", "pgtype.Timestampz":
		return "google.protobuf.Timestamp"
	case "uuid.UUID", "net.HardwareAddr", "net.IP":
		return "string"
	case "sql.Result":
		return "ExecResult"
	default:
		if _, elementType := originalAndElementType(typ); elementType != "" {
			return elementType
		}
		return typ
	}
}

func BindToProto(src, dst, attrName, attrType string) []string {
	res := make([]string, 0)
	switch attrType {
	case "sql.NullBool", "pgtype.Bool":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Bool(%s.%s.Bool) }", dst, CamelCaseProto(attrName), src, attrName))
	case "pgtype.Int2":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Int32(int32(%s.%s.Int16)) }", dst, CamelCaseProto(attrName), src, attrName))
	case "pgtype.Uint32":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.UInt32(%s.%s.Uint32) }", dst, CamelCaseProto(attrName), src, attrName))
	case "sql.NullInt32", "pgtype.Int4":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Int32(%s.%s.Int32) }", dst, CamelCaseProto(attrName), src, attrName))
	case "sql.NullInt64", "pgtype.Int8":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Int64(%s.%s.Int64) }", dst, CamelCaseProto(attrName), src, attrName))
	case "pgtype.Float4":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Float(%s.%s.Float32) }", dst, CamelCaseProto(attrName), src, attrName))
	case "sql.NullFloat64", "pgtype.Float8":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Double(%s.%s.Float64) }", dst, CamelCaseProto(attrName), src, attrName))
	case "sql.NullString", "pgtype.Text":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.String(%s.%s.String) }", dst, CamelCaseProto(attrName), src, attrName))
	case "sql.NullTime", "pgtype.Date", "pgtype.Timestamp", "pgtype.Timestampz":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = timestamppb.New(%s.%s.Time) }", dst, CamelCaseProto(attrName), src, attrName))
	case "time.Time":
		res = append(res, fmt.Sprintf("%s.%s = timestamppb.New(%s.%s)", dst, CamelCaseProto(attrName), src, attrName))
	case "uuid.UUID", "net.HardwareAddr", "net.IP":
		res = append(res, fmt.Sprintf("%s.%s = %s.%s.String()", dst, CamelCaseProto(attrName), src, attrName))
	case "pgtype.UUID":
		res = append(res, fmt.Sprintf("if v, err := json.Marshal(%s.%s); err == nil {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.String(string(v))", dst, CamelCaseProto(attrName)))
		res = append(res, "}")
	case "int16":
		res = append(res, fmt.Sprintf("%s.%s = int32(%s.%s)", dst, CamelCaseProto(attrName), src, attrName))
	default:
		_, elementType := originalAndElementType(attrType)
		if elementType != "" {
			res = append(res, fmt.Sprintf("%s.%s = %s(%s.%s)", dst, CamelCaseProto(attrName), elementType, src, attrName))
		} else {
			res = append(res, fmt.Sprintf("%s.%s = %s.%s", dst, CamelCaseProto(attrName), src, attrName))
		}
	}
	return res
}

func BindToGo(src, dst, attrName, attrType string, newVar bool) []string {
	res := make([]string, 0)
	switch attrType {
	case "sql.NullBool", "pgtype.Bool":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("%s = %s{Valid: true, Bool: v.Value}", dst, attrType))
		res = append(res, "}")
	case "pgtype.Int2":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("%s = %s{Valid: true, Int16: int16(v.Value)}", dst, attrType))
		res = append(res, "}")
	case "pgtype.Uint32":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("%s = %s{Valid: true, Uint32: v.Value}", dst, attrType))
		res = append(res, "}")
	case "sql.NullInt32", "pgtype.Int4":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("%s = %s{Valid: true, Int32: v.Value}", dst, attrType))
		res = append(res, "}")
	case "sql.NullInt64", "pgtype.Int8":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("%s = %s{Valid: true, Int64: v.Value}", dst, attrType))
		res = append(res, "}")
	case "pgtype.Float4":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("%s = %s{Valid: true, Float32: v.Value}", dst, attrType))
		res = append(res, "}")
	case "sql.NullFloat64", "pgtype.Float8":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("%s = %s{Valid: true, Float64: v.Value}", dst, attrType))
		res = append(res, "}")
	case "sql.NullString", "pgtype.Text":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("%s = %s{Valid: true, String: v.Value}", dst, attrType))
		res = append(res, "}")
	case "sql.NullTime", "pgtype.Date", "pgtype.Timestamp", "pgtype.Timestampz":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("if err := v.CheckValid(); err != nil { err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), validation.ErrUserInput)", attrName))
		res = append(res, "return nil, err }")
		res = append(res, "if t := v.AsTime(); !t.IsZero() {")
		res = append(res, fmt.Sprintf("%s.Valid = true", dst))
		res = append(res, fmt.Sprintf("%s.Time = t } }", dst))
	case "time.Time":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("if err := v.CheckValid(); err != nil { err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), validation.ErrUserInput)", attrName))
		res = append(res, "return nil, err }")
		res = append(res, fmt.Sprintf("%s = v.AsTime()", dst))
		res = append(res, fmt.Sprintf("} else { err := fmt.Errorf(\"field %s is required%%w\", validation.ErrUserInput)", attrName))
		res = append(res, "return nil, err }")
	case "uuid.UUID":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v, err := uuid.Parse(%s.Get%s()); err != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), validation.ErrUserInput)", attrName))
		res = append(res, fmt.Sprintf("return nil, err } else { %s = v }", dst))
	case "pgtype.UUID":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("if err := json.Unmarshal([]byte(v), &%s); err != nil {", dst))
		res = append(res, fmt.Sprintf("err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), validation.ErrUserInput)", attrName))
		res = append(res, "return nil, err }")
		res = append(res, "}")
	case "net.HardwareAddr":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v, err = net.ParseMAC(%s.Get%s()); err != nil {", src, CamelCaseProto(attrName)))
		res = append(res, fmt.Sprintf("err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), validation.ErrUserInput)", attrName))
		res = append(res, fmt.Sprintf("return nil, err } else { %s = v }", dst))
	case "net.IP":
		if newVar {
			res = append(res, fmt.Sprintf("%s := net.ParseIP(%s.Get%s())", dst, src, CamelCaseProto(attrName)))
		} else {
			res = append(res, fmt.Sprintf("%s = net.ParseIP(%s.Get%s())", dst, src, CamelCaseProto(attrName)))
		}
	case "int16":
		if newVar {
			res = append(res, fmt.Sprintf("%s := int16(%s.Get%s())", dst, src, CamelCaseProto(attrName)))
		} else {
			res = append(res, fmt.Sprintf("%s = int16(%s.Get%s())", dst, src, CamelCaseProto(attrName)))
		}
	case "int":
		if newVar {
			res = append(res, fmt.Sprintf("%s := int(%s.Get%s())", dst, src, CamelCaseProto(attrName)))
		} else {
			res = append(res, fmt.Sprintf("%s = int(%s.Get%s())", dst, src, CamelCaseProto(attrName)))
		}
	case "uint16":
		if newVar {
			res = append(res, fmt.Sprintf("%s := uint16(%s.Get%s())", dst, src, CamelCaseProto(attrName)))
		} else {
			res = append(res, fmt.Sprintf("%s = uint16(%s.Get%s())", dst, src, CamelCaseProto(attrName)))
		}
	default:
		originalType, elementType := originalAndElementType(attrType)
		if newVar {
			if elementType != "" {
				res = append(res, fmt.Sprintf("%s := %s(%s.Get%s())", dst, originalType, src, CamelCaseProto(attrName)))
			} else {
				res = append(res, fmt.Sprintf("%s := %s.Get%s()", dst, src, CamelCaseProto(attrName)))
			}
		} else {
			if elementType != "" {
				res = append(res, fmt.Sprintf("%s = %s(%s.Get%s())", dst, originalType, src, CamelCaseProto(attrName)))
			} else {
				res = append(res, fmt.Sprintf("%s = %s.Get%s()", dst, src, CamelCaseProto(attrName)))
			}
		}
	}
	return res
}

func UpperFirstCharacter(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return str
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToPascalCase(str string) string {
	return UpperFirstCharacter(generator.CamelCase(ToSnakeCase(str)))
}

func ToKebabCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}

func CamelCaseProto(str string) string {
	return generator.CamelCase(ToSnakeCase(str))
}

func originalAndElementType(typ string) (original, element string) {
	typ = strings.TrimPrefix(typ, "[]")
	t := strings.Split(typ, ".")
	return t[0], strings.Join(t[1:], ".")
}
