package gofmt

import (
	"fmt"
	"go/token"

	gofmtAPI "github.com/golangci/gofmt/gofmt"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "gofmt",
	Doc:  "Gofmt checks whether code was gofmt-ed and simplified   ",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	var fileNames []string
	for _, f := range pass.Files {
		pos := pass.Fset.PositionFor(f.Pos(), false)
		fileNames = append(fileNames, pos.Filename)
	}
	for _, f := range fileNames {
		diff, err := gofmtAPI.Run(f, true)
		if err != nil {
			return nil, err
		}
		if diff == nil {
			continue
		}

		//TODO improve error reporting
		text := fmt.Sprintf("file %v is not gofmt -s", f)
		pass.Report(analysis.Diagnostic{Pos: token.Pos(1), Message: text})
	}
	return nil, nil
}
