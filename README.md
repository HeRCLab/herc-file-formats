# herc-file-formats

[![HeRCLab](https://circleci.com/gh/HeRCLab/herc-file-formats.svg?style=svg)](https://app.circleci.com/pipelines/github/HeRCLab/herc-file-formats?branch=master) [![HeRCLab](https://goreportcard.com/badge/github.com/HeRCLab/herc-file-formats)](https://goreportcard.com/report/github.com/HeRCLab/herc-file-formats)

| MLPX Docs | Wavegen Docs |
|---|---|
| [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/herclab/herc-file-formats/mlpx) | [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/herclab/herc-file-formats/wavegen) |

This repository contains specifications and implementations for bespoke file
formats created for the HeRC lab, or for file formats we use in our work that
don't have good quality existing public implementations.

## How To Install

Because of the nature of this repository, all tools and libraries are usually
installed as a single bundle. However, check the documentation for each
sub-project, many can be installed independently of each other if you so desire

Usually, downloading the latest release as pre-compiled binaries from the
GitHub repository is preferred.

Otherwise, you can run `make build && sudo make install`.

The `./build_release.sh` script can be used to generate a `.deb` file.

Check each sub-project for dependencies.

## Structure

Each directory is a separate file format. Directories should be named after
their file format, and should contain a `doc/` directory documenting the
specification for that format, as well as any other relevant documentation.

Implementations should be in a sub-directory named after the implementation
language. So for example, a C implementation of TNX might reside in `./tnx/c/`.

Tools implemented in this repository should be placed in a sub-directory named
`cmd`. For example, a tool called `tnxutil` that operates on `tnx` files would
be long in `./tnx/cmd/tnxutil/`. To be clear, this applies to tools intended
for users to use and install on their system, not tools internal to this
repository used for building or CI.

### What belongs in This Repo

* Specifications of file formats used in the HeRC lab, either bespoke ones we
  have defined for ourselves, or third-party formats that we have had to work
  with.
* Libraries for inter-operating with such file formats.
* Tools whose express purpose is to interact with one of these file formats,
  or which are used to convert to or from such a file format.
	* Tools which have a broader purpose beyond say, validating or
	  converting a format belong in their own repos.
	* For example, a tool for converting MLPX files to ONNX could be placed
	  in `./mlpx/cmd/mlpx2onnx`, but a tool for executing MLPs which uses
	  MLPX as an input format should be placed in it's own repo.



## File Formats

| Format | Current Version | Status | Provenance | Purpose |
|-|-|-|-|-|
| [TNX (Trivial Network eXchange)](./tnx) | unreleased | on-hold | Bespoke | Format for representing compute graphs, with a focus on neurual networks. |
| [MLPX (MLP eXchange)](./mlpx) | unreleased | in-progress | Bespoke | Portable format for checkpointing MLP networks. |
| [Wavegen](./wavegen) | 0.0.3 | in-progress | Bespoke  | Wave generation tool & portable format for exchanging such waves. |


*Current Version* should be either "unreleased", for formats which have no
stable release yet, or the version number of the library contained in the
folder.  In general, all implementations of a format in different should have
the same version number.

*Status* should be one of "unreleased", "in-progress", "maintenance", or
"complete".
* **Unreleased formats have no stable public API, and may change or be removed
  at any time without notice.**

*Provenance* should be either "Bespoke" for in-house formats, or the origin of
the format otherwise.

## Versioning

Each individual format or tool may have it's own internal version. When a new
version needs to be released, a monotonically increasing version number will be
selected as the herc-file-formats release, and used for any generated packages.
This version number will attempt to follow semantic versioning to the extend
possible.

## Version History

### 0.0.1

Released formats:
* mlpx 0.0.1
* wavegen 0.0.4

### 0.0.2 (upcoming)


## Future Work

There are several existing bespoke formats, which should be "brought into the
fold" of this repo. These formats were created before this repo.

* Several bespoke image formats in `herc-imgtool` in the
  [herc-tools-public](https://github.com/HeRCLab/herc-tools-public) repo.
	* `de2rawh` AKA `hif24`
	* `de2raw`
	* `hif1`
	* `hif8`
* `herc-imgtool` itself


