# TNX

*Trivial Neural network eXchange*

**On Hiatus**: this project is tabled for the time being, as it's scope has
become rather large. The decision reached on July 20, 2020 was to leave this
project be for the present time, and instead build more application-specific
file formats, such as [mlpx](https://github.com/HeRCLab/mlpx). Later with the
benefit of hindsight, we can re-open this project when it's intented objectives
become more directly necessary for our use cases.

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
