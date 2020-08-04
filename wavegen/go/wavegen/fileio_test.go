package wavegen

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestJSONEncoding(t *testing.T) {
	w := &WaveParameters{
		SampleRate:           10,
		Offset:               0,
		Duration:             1,
		Frequencies:          []float64{1},
		Phases:               []float64{0},
		Amplitudes:           []float64{1},
		Noises:               []string{"none"},
		NoiseMagnitudes:      []float64{1.0},
		GlobalNoise:          "none",
		GlobalNoiseMagnitude: 0.0,
	}

	sig, err := w.GenerateSyntheticData()
	if err != nil {
		t.Error(err)
	}

	wf := &WaveFile{
		Version:    0,
		Parameters: w,
		Signal:     sig,
	}

	data, err := wf.ToJSON()
	if err != nil {
		t.Error(err)
	}

	decoded, err := FromJSON(data)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(wf, decoded) {
		t.Logf("wf = %v", wf)
		t.Logf("decoded = %v", decoded)
		t.Errorf("JSON decoding is not idempotent")
	}

	text := `
{
  "Version": 0,
  "Parameters": {
    "SampleRate": 10,
    "Offset": 0,
    "Duration": 1,
    "Frequencies": [
      1
    ],
    "Phases": [
      0
    ],
    "Amplitudes": [
      1
    ],
    "Noises": [
      "none"
    ],
    "NoiseMagnitudes": [
      1
    ],
    "GlobalNoise": "none",
    "GlobalNoiseMagnitude": 0
  }
}
`
	decoded, err = FromJSON([]byte(text))
	if err != nil {
		t.Error(err)
	}
	wf.Signal = nil
	if !cmp.Equal(wf, decoded) {
		t.Logf("wf = %v", wf)
		t.Logf("decoded = %v", decoded)
		t.Errorf("JSON decoding does not handle nil signal correctly")
	}

}
