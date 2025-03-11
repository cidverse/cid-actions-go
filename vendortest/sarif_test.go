package vendortest

import (
	_ "embed"
	"testing"

	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
	"github.com/stretchr/testify/assert"
)

//go:embed files/report.json
var reportJson string

func TestSarifParser(t *testing.T) {
	result, err := sarif.FromBytes([]byte(reportJson))
	assert.NoError(t, err)
	assert.Len(t, result.Runs, 1)
	assert.Len(t, result.Runs[0].Tool.Driver.Rules, 1)
	assert.Len(t, result.Runs[0].Results, 1)
}
