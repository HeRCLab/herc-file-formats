package mlpx

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ToJSON converts an existing MLPX object to a JSON string and returns it.
func (mlp *MLPX) ToJSON() ([]byte, error) {
	b, err := json.MarshalIndent(mlp, "", "\t")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// WriteJSON calls ToJSON() and then overwrites the specified path with it's
// return.
func (mlp *MLPX) WriteJSON(path string) error {
	b, err := mlp.ToJSON()
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

// FromJSON reads an in-memory JSON string and generates an MLPX object. It
// does not validate the data which is read.
func FromJSON(data []byte) (*MLPX, error) {
	mlp := &MLPX{}
	err := json.Unmarshal(data, mlp)
	if err != nil {
		return nil, err
	}
	return mlp, err
}

// ReadJSON is a utility function which reads a file from disk, then calls
// FromJSON() on it. It does not validate the MLPX file.
func ReadJSON(path string) (*MLPX, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	mlp, err := FromJSON(data)
	if err != nil {
		return nil, err
	}

	return mlp, nil
}
