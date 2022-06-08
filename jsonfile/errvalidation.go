package jsonfile

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"github.com/jamestunnell/topdown/sliceutil"
)

type ErrValidation struct {
	Errors []gojsonschema.ResultError
}

func (err *ErrValidation) Error() string {
	errStrings := sliceutil.Map(err.Errors, resultErrToString)
	errsStr := strings.Join(errStrings, "\n")

	return fmt.Sprintf("JSON validation result contains errors: \n%s", errsStr)
}

func resultErrToString(r gojsonschema.ResultError) string {
	return r.String()
}
