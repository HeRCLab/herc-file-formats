Scala MLPX Library
------------------

This project provides a Scala package named `mlpx`, which implements basic support for reading MLPX files.


## Classes

 - `MLPXDocument` :: Top-level object that handles deserialization of the MLPX file. Contains a `Map[String, Snapshot]` containing all the snapshots.
 - `Snapshot` :: Simple object representing a complete snapshot of an MLP network. Contains a `Map[String, Layer]` containing all the layers in that snapshot.
 - `Layer` :: Simple object representing a layer in an MLP network.


## Usage examples

Loading an MLPX file quickly:

    val mlpx_file = mlpx.Utils.fileRead("example.mlpx")
    val doc = new mlpx.MLPXDocument(mlpx_file)

Print the snapshot IDs in an MLPX file in sorted order:

    for (id <- doc.sortedSnapshotIDs) println(id)

Extract a particular snapshot:

    doc.snapshots("1")

Extract the initializer snapshot:

    // Either of these will work:
    val Some(init) = doc.initializer  // Returns an Option[Snapshot] type.
    
    val init = doc.snapshots("initializer")

Get the schema version for the MLPX document:

    doc.version

Get a particular layer of a particular snapshot:

    doc.snapshots("2").layers("input")
