package metadata

import (
	"strings"

	"github.com/walterwanderley/sqlc-grpc/converter"
)

type Service struct {
	Name                string
	InputNames          []string
	InputTypes          []string
	Output              string
	Sql                 string
	Messages            map[string]*Message
	CustomProtoComments []string
	CustomProtoOptions  []string
}

func (s *Service) ParamsCallDatabase() string {
	if s.EmptyInput() {
		return ""
	}
	return ", " + strings.Join(s.InputNames, ", ")
}

func (s *Service) HasCustomParams() bool {
	if s.EmptyInput() {
		return false
	}

	return customType(s.InputTypes[0], s.Messages)
}

func (s *Service) HasSimpleParams() bool {
	if s.HasArrayParams() {
		return false
	}

	if !s.HasCustomParams() || s.EmptyInput() {
		return true
	}

	if msg, ok := s.Messages[converter.CanonicalName(s.InputTypes[0])]; ok {
		return !msg.HasComplexAttribute(s.Messages)
	}

	return false
}

func (s *Service) HasArrayParams() bool {
	if s.EmptyInput() {
		return false
	}

	return strings.HasPrefix(s.InputTypes[0], "[]") && s.InputTypes[0] != "[]byte"
}

func (s *Service) HasCustomOutput() bool {
	if s.EmptyOutput() {
		return false
	}

	return customType(s.Output, s.Messages)
}

func (s *Service) HasArrayOutput() bool {
	if s.EmptyOutput() {
		return false
	}
	return strings.HasPrefix(s.Output, "[]") && s.Output != "[]byte"
}

func (s *Service) EmptyInput() bool {
	return len(s.InputTypes) == 0
}

func (s *Service) EmptyOutput() bool {
	return s.Output == ""
}
