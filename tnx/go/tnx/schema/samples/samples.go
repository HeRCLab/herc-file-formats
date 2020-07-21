// Package samples is used to store sample TNX objects that can be used for
// testing or demonstration purposes.
//
// To avoid circular imports, samples are returned as []byte, but should
// be able to be loaded using the `FromJSON()` method.

package samples

func SampleMLP3Layer() []byte {
	return []byte(`
{
  "schema": [
    "tnx",
    "0"
  ],
  "topology": {
    "nodes": [
      {
        "id": "input",
        "operation": "output",
        "outputs": [ "input->output0" ]
      },
      {
        "id": "hidden1",
        "operation": "mlplayer",
        "inputs": [ "hidden1<-input0" ],
        "outputs": [ "hidden1->output0" ]
      },
      {
        "id": "activaton1",
        "operation": "relu",
        "inputs": [ "activation1<-input0" ],
        "outputs": [ "activation1->output0" ]
      },
      {
        "id": "hidden2",
        "operation": "mlplayer",
        "inputs": [ "hidden2<-input0" ],
        "outputs": [ "hidden2->output0" ]
      },
      {
        "id": "activaton2",
        "operation": "relu",
        "inputs": [ "activation2<-input0" ],
        "outputs": [ "activation2->output0" ]
      },
      {
        "id": "hidden3",
        "operation": "mlplayer",
        "inputs": [ "hidden3<-input0" ],
        "outputs": [ "hidden3->output0" ]
      },
      {
        "id": "activaton3",
        "operation": "relu",
        "input": [ "activation3<-input0" ],
        "outputs": [ "activation3->output0" ]
      },
      {
        "id": "output",
        "operation": "output",
        "inputs": [ "output<-input0" ]
      }
    ],
    "links": [
      {
        "source": "input->output0",
        "target": "hidden1<-input0"
      },
      {
        "source": "hidden1->output0",
        "target": "activation1<-input0"
      },
      {
        "source": "activation1->output0",
        "target": "hidden2<-input0"
      },
      {
        "source": "hidden2->output0",
        "target": "activation2<-input0"
      },
      {
        "source": "activation2->output0",
        "target": "hidden3<-input0"
      },
      {
        "source": "hidden3->output0",
        "target": "activation3<-input0"
      },
      {
        "source": "activation3->output0",
        "target": "output<-input0"
      }
    ]
  },
  "parameters": {
    "input": {
      "dimensions": [
        25
      ]
    },
    "hidden1": {
      "neurons": 25,
      "activation": "activation1"
    },
    "hidden2": {
      "neurons": 15,
      "activation": "activation2"
    },
    "hidden3": {
      "neurons": 15,
      "activation": "activation3"
    },
    "output": {
      "dimensions": [
        5
      ]
    }
  }
}
`)

}
