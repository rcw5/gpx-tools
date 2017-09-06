To simplify a route: https://github.com/gpsbabel/gpsbabel/blob/8e968d504e001b2df844aea4c8b6b420cde18652/smplrout.cc
GPSBAbel ^^^

* It's pretty simple, really: for each triplet of vertices A-B-C, we compute
 * how much cross-track error we'd introduce by going straight from A to C
 * (the maximum cross-track error for that segment is the height of the
 * triangle ABC, measured between vertex B and edge AC.)  If we need to remove
 * 40 points, we just sort the points by that metric and remove the 40
 * smallest ones.
 *
 * It's actually a little more complicated than that, because removing a
 * point changes the result for its two nearest neighbors.  When we remove
 * one, we recompute the neighbors and then sort them back into the list
 * at their new locations.




 