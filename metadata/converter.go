package metadata

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
	"unicode"
)

func exprToStr(e ast.Expr) string {
	switch exp := e.(type) {
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", exprToStr(exp.X), exp.Sel.Name)
	case *ast.Ident:
		return exp.String()
	case *ast.StarExpr:
		return "*" + exprToStr(exp.X)
	case *ast.ArrayType:
		return "[]" + exprToStr(exp.Elt)
	default:
		panic(fmt.Sprintf("invalid type %v", exp))
	}
}

func toProtoType(typ string) string {
	if strings.HasPrefix(typ, "*") {
		return toProtoType(typ[1:])
	}
	if strings.HasPrefix(typ, "[]") && typ != "[]byte" {
		return "repeated " + toProtoType(typ[2:])
	}
	switch typ {
	case "json.RawMessage", "[]byte":
		return "bytes"
	case "sql.NullBool":
		return ".google.protobuf.BoolValue"
	case "sql.NullInt32":
		return ".google.protobuf.Int32Value"
	case "int":
		return "int64"
	case "int16":
		return "int32"
	case "sql.NullInt64":
		return ".google.protobuf.Int64Value"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "sql.NullFloat64":
		return ".google.protobuf.DoubleValue"
	case "sql.NullString":
		return ".google.protobuf.StringValue"
	case "sql.NullTime", "time.Time":
		return ".google.protobuf.Timestamp"
	case "uuid.UUID", "net.HardwareAddr", "net.IP":
		return "string"
	default:
		return typ
	}
}

func bindToProto(src, dst, attrName, attrType string) []string {
	res := make([]string, 0)
	switch attrType {
	case "sql.NullBool":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Bool(%s.%s.Bool) }", dst, attrName, src, attrName))
	case "sql.NullInt32":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Int32(%s.%s.Int32) }", dst, attrName, src, attrName))
	case "sql.NullInt64":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Int64(%s.%s.Int64) }", dst, attrName, src, attrName))
	case "sql.NullFloat64":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Float64(%s.%s.Float64) }", dst, attrName, src, attrName))
	case "sql.NullString":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.String(%s.%s.String) }", dst, attrName, src, attrName))
	case "sql.NullTime":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = timestamppb.New(%s.%s.Time) }", dst, attrName, src, attrName))
	case "time.Time":
		res = append(res, fmt.Sprintf("%s.%s = timestamppb.New(%s.%s)", dst, attrName, src, attrName))
	case "uuid.UUID", "net.HardwareAddr", "net.IP":
		res = append(res, fmt.Sprintf("%s.%s = %s.%s.String()", dst, attrName, src, attrName))
	case "int16":
		res = append(res, fmt.Sprintf("%s.%s = int32(%s.%s)", dst, attrName, src, attrName))
	default:
		res = append(res, fmt.Sprintf("%s.%s = %s.%s", dst, attrName, src, attrName))
	}
	return res
}

func bindToGo(src, dst, attrName, attrType string, newVar bool) []string {
	res := make([]string, 0)
	switch attrType {
	case "sql.NullBool":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("%s = sql.NullBool{Valid: true, Bool: v.Value}", dst))
		res = append(res, "}")
	case "sql.NullInt32":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("%s = sql.NullInt32{Valid: true, Int32: v.Value}", dst))
		res = append(res, "}")
	case "sql.NullInt64":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("%s = sql.NullInt64{Valid: true, Int64: v.Value}", dst))
		res = append(res, "}")
	case "sql.NullFloat64":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("%s = sql.NullFloat64{Valid: true, Float64: v.Value}", dst))
		res = append(res, "}")
	case "sql.NullString":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("%s = sql.NullString{Valid: true, String: v.Value}", dst))
		res = append(res, "}")
	case "sql.NullTime":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("if err = v.CheckValid(); err != nil { err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), validation.ErrUserInput)", attrName))
		res = append(res, "return }")
		res = append(res, "if t := v.AsTime(); !t.IsZero() {")
		res = append(res, fmt.Sprintf("%s.Valid = true", dst))
		res = append(res, fmt.Sprintf("%s.Time = t } }", dst))
	case "time.Time":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("if err = v.CheckValid(); err != nil { err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), validation.ErrUserInput)", attrName))
		res = append(res, "return }")
		res = append(res, fmt.Sprintf("%s = v.AsTime()", dst))
		res = append(res, fmt.Sprintf("} else { err = fmt.Errorf(\"%s is required%%w\", validation.ErrUserInput)", attrName))
		res = append(res, "return }")
	case "uuid.UUID":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if %s, err = uuid.Parse(%s.Get%s()); err != nil {", dst, src, attrName))
		res = append(res, fmt.Sprintf("err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), validation.ErrUserInput)", attrName))
		res = append(res, "return }")
	case "net.HardwareAddr":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if %s, err = net.ParseMAC(%s.Get%s()); err != nil {", dst, src, attrName))
		res = append(res, fmt.Sprintf("err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), validation.ErrUserInput)", attrName))
		res = append(res, "return }")
	case "net.IP":
		if newVar {
			res = append(res, fmt.Sprintf("%s := net.ParseIP(%s.Get%s())", dst, src, attrName))
		} else {
			res = append(res, fmt.Sprintf("%s = net.ParseIP(%s.Get%s())", dst, src, attrName))
		}
	case "int16":
		if newVar {
			res = append(res, fmt.Sprintf("%s := int16(%s.Get%s())", dst, src, attrName))
		} else {
			res = append(res, fmt.Sprintf("%s = int16(%s.Get%s())", dst, src, attrName))
		}
	default:
		if newVar {
			res = append(res, fmt.Sprintf("%s := %s.Get%s()", dst, src, attrName))
		} else {
			res = append(res, fmt.Sprintf("%s = %s.Get%s()", dst, src, attrName))
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

func toKebabCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}
