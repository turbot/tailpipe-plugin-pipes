package source

import (
	"testing"

	"github.com/turbot/tailpipe-plugin-sdk/source"
)

func TestConformance(t *testing.T) {
	source.RunConformanceTests(t, &PipesAuditLogSource{})
}
