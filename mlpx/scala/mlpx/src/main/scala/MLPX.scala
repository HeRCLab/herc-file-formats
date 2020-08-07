package mlpx

import play.api.libs.json._

// Case classes are similar to C structs: just a data definition.
// These ones are carefully constructed to play nice with play-json.
case class Layer(predecessor: String,
                 successor: String,
                 neurons: Int,
                 weights: Option[List[Double]],
                 outputs: Option[List[Double]],
                 activations: Option[List[Double]],
                 deltas: Option[List[Double]],
                 biases: Option[List[Double]],
                 activation_function: Option[String])

case class Snapshot(layers: Map[String, Layer],
                    alpha: Option[Double])

// These companion object definitions allow automatic deserialization routines
// to be generated at compile time.
object Layer {
  implicit val layerReads = Json.reads[Layer]
}

object Snapshot {
  implicit val snapshotReads = Json.reads[Snapshot]
}

// Note: Some hackery required here due to weird schema key format.
class MLPXDocument(jsonSource: String) {
  private val _raw_json = Json.parse(jsonSource)
  var version = (_raw_json \ "schema").get(1).as[Double]
  var snapshots = (_raw_json \ "snapshots").as[Map[String, Snapshot]]

  // Returns [] if no snapshots present.
  // Returns ["initializer", "1", "2", ...] if IDs present.
  def sortedSnapshotIDs(): List[String] = {
    var initializer = List[String]()
    var numeric_ids = List[String]()
    var other_ids   = List[String]()
    for (k <- snapshots.keys) {
      if (k == "initializer") {
        initializer = List("initializer")
      // Cite: https://stackoverflow.com/a/9946141
      } else if (k.forall(_.isDigit)) {
        numeric_ids = numeric_ids :+ k // Append k to numeric ID list
      } else {
        other_ids = other_ids :+ k
      }
    }
    numeric_ids = numeric_ids.sortWith(_.toInt < _.toInt)
    initializer ++ numeric_ids ++ other_ids
  }

  def initializer(): Option[Snapshot] = {
    snapshots get "initializer"
  }  
}
