#include "test_util.h"
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include "mlpx.h"

int main(int argc, char** argv) {
	int handle, snapc, layerc, layeridx;
	char* layerid;
	should_equal(MLPXOpen("test1.mlpx", &handle), 0);

	should_equal(MLPXGetNumSnapshots(handle, &snapc), 0);
	should_equal(snapc, 1);

	should_equal(MLPXGetNumLayers(handle, 0, &layerc), 0);
	should_equal(layerc, 3);

	should_equal(MLPXGetLayerIDByIndex(handle, 0, 1, &layerid), 0);
	str_should_equal(layerid, "hidden0");

	should_equal(MLPXClose(handle), 0);
}
