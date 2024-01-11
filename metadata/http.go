package metadata

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/walterwanderley/sqlc-grpc/converter"
)

func (s *Service) HttpMethod() string {
	query := trimHeaderComments(strings.ReplaceAll(s.Sql, "`", ""))
	query = strings.ToUpper(query)
	if strings.HasPrefix(query, "SELECT") && s.HasSimpleParams() {
		return "get"
	}
	if strings.HasPrefix(query, "DELETE") && s.HasSimpleParams() {
		return "delete"
	}
	if strings.HasPrefix(query, "UPDATE") {
		return "put"
	}

	return "post"
}

func trimHeaderComments(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "--") && !strings.HasPrefix(s, "/*") {
		return s
	}
	i := strings.Index(s, "\n")
	if i != -1 {
		return trimHeaderComments(s[i+1:])
	}
	return strings.TrimSpace(s)
}

func (s *Service) HttpPath() string {
	path := "/" + converter.ToKebabCase(removePrefix(s.Name))
	method := s.HttpMethod()

	if (method == "get" || method == "delete") &&
		len(s.InputNames) == 1 && !s.HasCustomParams() && !s.HasArrayParams() {
		path = fmt.Sprintf("%s/{%s}", path, converter.ToSnakeCase(converter.CanonicalName(s.InputNames[0])))
	}
	return path
}

func (s *Service) HttpBody() string {
	switch s.HttpMethod() {
	case "get", "delete":
		return ""
	default:
		if s.HasArrayParams() {
			return s.InputNames[0]
		}
		return "*"
	}
}

func (s *Service) HttpResponseBody() string {
	if s.HasArrayOutput() {
		return "list"
	} else if s.HasCustomOutput() {
		return converter.ToSnakeCase(converter.CanonicalName(s.Output))
	}
	return ""
}

func (s *Service) HttpOptions() []string {
	if len(s.CustomProtoOptions) > 0 {
		return s.CustomProtoOptions
	}
	res := make([]string, 0)
	res = append(res, "option (google.api.http) = {")
	res = append(res, fmt.Sprintf("    %s: \"%s\"", s.HttpMethod(), s.HttpPath()))
	body := s.HttpBody()
	if body != "" {
		res = append(res, fmt.Sprintf("    body: \"%s\"", body))
	}
	responseBody := s.HttpResponseBody()
	if responseBody != "" {
		res = append(res, fmt.Sprintf("    response_body: \"%s\"", responseBody))
	}
	res = append(res, "};")
	return res
}

func removePrefix(s string) string {
	p := prefix(s)
	if p == s {
		return s
	}

	p = strings.ToLower(p)
	for _, reserved := range []string{
		"create", "add", "insert", "list", "get", "read", "update", "modify", "delete", "remove",
	} {
		if strings.HasPrefix(p, reserved) {
			return s[len(p):]
		}
	}

	return s
}

func prefix(s string) string {
	var res = make([]rune, 0, len(s))
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			return string(res)
		} else {
			res = append(res, r)
		}
	}
	return string(res)
}
