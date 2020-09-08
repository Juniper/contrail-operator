package goimports

import (
	"fmt"
	"go/token"

	goimportsAPI "github.com/golangci/gofmt/goimports"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/imports"
)

var Analyzer = &analysis.Analyzer{
	Name: "goimports",
	Doc:  "Goimports checks whether imports are correctly listed",
	Run:  run,
}

func init() {
	imports.LocalPrefix = "github.com/Juniper/contrail-operator"
}

func run(pass *analysis.Pass) (interface{}, error) {
	var fileNames []string
	for _, f := range pass.Files {
		pos := pass.Fset.PositionFor(f.Pos(), false)
		fileNames = append(fileNames, pos.Filename)
	}
	for _, f := range fileNames {
		diff, err := goimportsAPI.Run(f)
		if err != nil {
			return nil, err
		}
		if diff == nil {
			continue
		}

		//TODO improve error reporting
		text := fmt.Sprintf("file %v is not formated well please run: goimports -local github.com/Juniper/contrail-operator ", f)
		pass.Report(analysis.Diagnostic{Pos: token.Pos(1), Message: text})
	}
	return nil, nil
}
