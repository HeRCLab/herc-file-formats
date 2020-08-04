package wavegen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// WaveFile represents all of the data that could be encoded in a WaveGen JSON
// file.
type WaveFile struct {
	// Version should be the version specifier, currently 0.
	Version int

	// Parameters represents the WaveParameters.
	Parameters *WaveParameters `json:",omitempty"`

	// Signal represents the wave data.
	Signal *Signal `json:",omitempty"`
}

// ToJSON converts a WaveFile to an in-memory JSON representation and returns
// it.
func (wf *WaveFile) ToJSON() ([]byte, error) {
	if wf.Version != 0 {
		return nil, fmt.Errorf("Don't know how to write a wave file with version %d", wf.Version)
	}

	b, err := json.MarshalIndent(wf, "", "\t")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// WriteJSON generates a JSON representation using ToJSON(), then writes
// it to disk.
func (wf *WaveFile) WriteJSON(path string) error {
	b, err := wf.ToJSON()
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// FromJSON loads the data stored in a JSON representation of a wavegen file.
func FromJSON(data []byte) (*WaveFile, error) {
	wf := &WaveFile{}
	err := json.Unmarshal(data, wf)
	if err != nil {
		return nil, err
	}
	return wf, nil
}

// ReadJSON loads the data stored in a JSON file on disk to a wavegen file
// object.
func ReadJSON(path string) (*WaveFile, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	wf, err := FromJSON(data)
	if err != nil {
		return nil, err
	}

	return wf, nil
}
