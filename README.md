# herc-file-formats

This repository contains specifications and implementations for bespoke file
formats created for the HeRC lab, or for file formats we use in our work that
don't have good quality existing public implementations.

## Structure

Each directory is a separate file format. Directories should be named after
their file format, and should contain a `doc/` directory documenting the
specification for that format, as well as any other relevant documentation.

Implementations should be in a sub-directory named after the implementation
language. So for example, a C implementation of TNX might reside in `./tnx/c/`.

## File Formats

| Format | Current Version | Status | Provenance | Purpose |
|-|-|-|-|-|
| [TNX (Trivial Network eXchange)](./tnx) | unreleased | on-hold | Bespoke | Format for representing compute graphs, with a focus on neurual networks. |
| [MLPX (MLP eXchange)](./mlpx) | unreleased | in-progress | Bespoke | Portable format for checkpointing MLP networks. |


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

When a release of a file format/implementation is complete, the relevant commit
should be tagged with `formatname-versionnumber`. A valid example might be
`tnx-0.0.3`.
