package samples

import (
	"testing"

	"github.com/herclab/tnx/go/tnx/schema"
)

func TestSampleMLP3Layer(t *testing.T) {
	_, err := schema.FromJSON(SampleMLP3Layer())
	if err != nil {
		t.Error(err)
	}
}
