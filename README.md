# herc-file-formats

This repository contains specifications and implementations for bespoke file
formats created for the HeRC lab, or for file formats we use in our work that
don't have good quality existing public implementations.

## Structure

Each directory is a separate file format. Directories should be named after
their file format, and should contain a `doc/` directory documenting the
specification for that format, as well as any other relevant documentation.

## File Formats

| Format | Current Version | Status | Purpose |
|-|-|-|-|
| [TNX (Trivial Network eXchange)](./tnx) | unreleased | on-hold | Format for representing compute graphs, with a focus on neurual networks. |
| [MLPX (MLP eXchange)](./mlpx) | unreleased | in-progress | Portable format for checkpointing MLP networks. |
