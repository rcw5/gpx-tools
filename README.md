# GPX Simplifier

## What does it do?

This small CLI performs operations on GPX tracks:

- Splitting a track into multiple files
- Simplifying a track by reducing the number of trackpoints

This is useful for owners of Garmin GPX devices which have been known to truncate tracks that are too long or too complex.

## How do I use it?

```
NAME:
  gpx-tools - Simple suite of tools to manipulate GPX files

USAGE:
  gpx-tools [global options] command [command options] [arguments...]

VERSION:
  0.0.3

COMMANDS:
    simplify, sim  Simplify a GPX File
    split, spl     Split a GPX file into a number of smaller files
    help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
  --help, -h     show help
  --version, -v  print the version
```

## How does it work?

Simplifying a track involves calculating the cross-track-error (XTE) of each point in the route:

> It's pretty simple, really: for each triplet of vertices A-B-C, we compute how much cross-track error we'd introduce by going straight from A to C (the maximum cross-track error for that segment is the height of the triangle ABC, measured between vertex B and edge AC.)  If we need to remove 40 points, we just sort the points by that metric and remove the 40 smallest ones.

> It's actually a little more complicated than that, because removing a
point changes the result for its two nearest neighbors.  When we remove one, we recompute the neighbors and then sort them back into the list at their new locations.

(Mostly taken from https://github.com/gpsbabel/gpsbabel/blob/8e968d504e001b2df844aea4c8b6b420cde18652/smplrout.cc)
