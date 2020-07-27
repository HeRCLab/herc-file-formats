#include "test_util.h"
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include "mlpx.h"

int main(int argc, char** argv) {
	int handle, snapc, layerc, layeridx, neuronc, initializer, pred, succ;
	char* layerid;
	char* funct;
	double weight, output, activation, delta, bias, alpha;
	should_equal(MLPXOpen("test1.mlpx", &handle), 0);

	should_equal(MLPXGetNumSnapshots(handle, &snapc), 0);
	should_equal(snapc, 1);

	should_equal(MLPXSnapshotGetNumLayers(handle, 0, &layerc), 0);
	should_equal(layerc, 3);

	should_equal(MLPXLayerGetIDByIndex(handle, 0, 1, &layerid), 0);
	str_should_equal(layerid, "hidden0");

	should_equal(MLPXLayerGetNeurons(handle, 0, 1, &neuronc), 0);
	should_equal(neuronc, 2);

	should_equal(MLPXGetInitializerSnapshotIndex(handle, &initializer), 0);
	should_equal(initializer, 0);

	should_equal(MLPXLayerGetPredecessorIndex(handle, 0, 1, &pred), 0);
	should_equal(pred, 0);

	should_equal(MLPXLayerGetSuccessorIndex(handle, 0, 1, &succ), 0);
	should_equal(succ, 2);

	should_equal(MLPXLayerSetWeight(handle, 0, 1, 3, 2.7), 0);
	should_equal(MLPXLayerGetWeight(handle, 0, 1, 3, &weight), 0);
	should_equal_epsilon(weight, 2.7, 0.00001);

	should_equal(MLPXMakeIsomorphicSnapshot(handle, "1", 0), 0);
	should_equal(MLPXLayerSetWeight(handle, 1, 1, 3, 3.7), 0);
	should_equal(MLPXLayerGetWeight(handle, 1, 1, 3, &weight), 0);
	should_equal_epsilon(weight, 3.7, 0.00001);
	should_equal(MLPXLayerGetWeight(handle, 0, 1, 3, &weight), 0);
	should_equal_epsilon(weight, 2.7, 0.00001);

	should_equal(MLPXLayerSetOutput(handle, 0, 1, 1, 5.2), 0);
	should_equal(MLPXLayerGetOutput(handle, 0, 1, 1, &output), 0);
	should_equal_epsilon(output, 5.2, 0.00001);

	should_equal(MLPXLayerSetActivation(handle, 0, 1, 1, 6.2), 0);
	should_equal(MLPXLayerGetActivation(handle, 0, 1, 1, &activation), 0);
	should_equal_epsilon(activation, 6.2, 0.00001);

	should_equal(MLPXLayerSetDelta(handle, 0, 1, 1, 7.2), 0);
	should_equal(MLPXLayerGetDelta(handle, 0, 1, 1, &delta), 0);
	should_equal_epsilon(delta, 7.2, 0.00001);

	should_equal(MLPXLayerSetBias(handle, 0, 1, 1, 8.2), 0);
	should_equal(MLPXLayerGetBias(handle, 0, 1, 1, &bias), 0);
	should_equal_epsilon(bias, 8.2, 0.00001);

	should_equal(MLPXLayerGetActivationFunction(handle, 0, 1, &funct), 0);
	str_should_equal(funct, "foobar");
	should_equal(MLPXLayerSetActivationFunction(handle, 0, 1, "baz"), 0);
	should_equal(MLPXLayerGetActivationFunction(handle, 0, 1, &funct), 0);
	str_should_equal(funct, "baz");

	should_equal(MLPXSnapshotGetAlpha(handle, 0, &alpha), 0);
	should_equal_epsilon(alpha, 0.1, 0.0001);
	should_equal(MLPXSnapshotSetAlpha(handle, 0, 0.2), 0);
	should_equal(MLPXSnapshotGetAlpha(handle, 0, &alpha), 0);
	should_equal_epsilon(alpha, 0.2, 0.0001);

	should_equal(MLPXClose(handle), 0);
}
