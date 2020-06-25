# TNX

*Trivial Neural network eXchange*

**WiP**: this repo does not contain any working code yet. The proposed format
can be found in [`./doc`](./doc).

This repository contains code to interface with data in the TNX format, created
to support several of [HeRC Lab](https://cse.sc.edu/~jbakos/group/)'s ongoing
projects. It also includes a tool to convert TNX to and from
[ONNX](http://onnx.ai/) format. Note that only a limited subset of ONNX is
supported.


## Motivation

The motivation of the TNX format is to create a means for storing typologies of
simple multilayer neural networks, as well as not only the needed
initialization values of the network elements, but also intermediary and output
values related to a specific execution of that network. A key design goal is to
enable multiple network implementations on different platforms/languages to be
validated against one another.
