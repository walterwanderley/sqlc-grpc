package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

type sqlcConfigV2 struct {
	SQL []struct {
		Engine string `json:"engine,omitempty" yaml:"engine"`
		Gen    struct {
			Go *sqlGoConfig `json:"go,omitempty" yaml:"go"`
		} `json:"gen" yaml:"gen"`
	} `json:"sql" yaml:"sql"`
}

func (v2 sqlcConfigV2) toV1() (v1 sqlcConfig) {
	v1.Packages = make([]PackageConfig, 0)
	for _, sql := range v2.SQL {
		if sql.Gen.Go == nil {
			continue
		}
		v1.Packages = append(v1.Packages, PackageConfig{
			Name:                      sql.Gen.Go.Package,
			Path:                      sql.Gen.Go.Out,
			Engine:                    sql.Engine,
			EmitInterface:             sql.Gen.Go.EmitInterface,
			EmitResultStructPointers:  sql.Gen.Go.EmitResultStructPointers,
			EmitParamsStructPointers:  sql.Gen.Go.EmitParamsStructPointers,
			EmitMethodsWithDBArgument: sql.Gen.Go.EmitMethodsWithDBArgument,
			SqlPackage:                sql.Gen.Go.SqlPackage,
		})
	}
	return
}

type sqlGoConfig struct {
	EmitInterface             bool   `json:"emit_interface" yaml:"emit_interface"`
	EmitResultStructPointers  bool   `json:"emit_result_struct_pointers" yaml:"emit_result_struct_pointers"`
	EmitParamsStructPointers  bool   `json:"emit_params_struct_pointers" yaml:"emit_params_struct_pointers"`
	EmitMethodsWithDBArgument bool   `json:"emit_methods_with_db_argument,omitempty" yaml:"emit_methods_with_db_argument"`
	Package                   string `json:"package" yaml:"package"`
	Out                       string `json:"out" yaml:"out"`
	SqlPackage                string `json:"sql_package" yaml:"sql_package"`
}

func readConfigV2(name string) (sqlcConfig, error) {
	var cfgV2 sqlcConfigV2
	f, err := os.Open(name)
	if err != nil {
		return sqlcConfig{}, err
	}
	defer f.Close()

	switch name {
	case jsonConfig:
		err = json.NewDecoder(f).Decode(&cfgV2)
	default:
		err = yaml.NewDecoder(f).Decode(&cfgV2)
	}
	if err != nil {
		return sqlcConfig{}, err
	}
	return cfgV2.toV1(), nil
}
