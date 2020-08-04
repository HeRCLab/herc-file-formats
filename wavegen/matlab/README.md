# WaveGen Matlab Client

[back to README](../../README.md)

WaveGen's Matlab client is implemented using Matlab's native JSON decoder
method. The client is implemented in [`wavegen.m`](./wavegen.m), and an example
of it's use can be seen in [`example.m`](./example.m). It is recommended to
simply copy the client code into your project folder.

Note that the `wavegen` object's fields are immutable non-member methods to
avoid getting the object into an inconsistent state. If you need to modify the
loaded data, you should first copy it outside of the object.
