# Alien Invasion

## Overview
This program attempts to solve the **Alien Invasion** coding puzzle. More
details regarding the problem will be added to this README soon.

## Requirements
At present, this code is only being tested on MacOS, but should compile just as
well on most Linux distros that support Golang.

* Tested with Golang v1.11.2

## Building
To build, once you've cloned the repo, simply do the following:

```bash
> make
```

This will build an executable `alien-invasion` in the root of the repo. If you
do not have [dep](https://golang.github.io/dep/) installed, the `Makefile`
should take care of installing that into your `GOPATH` for you.

## Running
Once you've built the `alien-invasion` executable, to run the simulator with the
example world and a given number of aliens, simply just run the following:

```bash
# Run with 3 aliens
> ./alien-invasion -N 3 --use-example-map
```

For more help, simply run:

```bash
> ./alien-invasion --help
Alien invasion simulator! See https://github.com/thanethomson/alien-invasion for more details.

Usage:
  alien-invasion [flags]

Flags:
  -N, --alien-count int    the number of aliens to simulate (default 2)
  -h, --help               help for alien-invasion
      --use-example-map    use the example world map instead of loading one
  -m, --world-map string   the file from which to load the world map (default "world-map.txt")
```

## World Map
World maps are stored in text files such that each line gives information about
a particular city. For example:

```
Foo north=Bar west=Baz south=Qu-ux
Bar south=Foo west=Bee
```

The first line indicates that there is a city called `Foo`. To the north of
`Foo` is a city called `Bar`, to the west of `Foo` is a city called `Baz`, and
to the south of `Foo` is a city called `Qu-ux`.

See the [maps](./maps/) folder for some example maps.

## Assumptions
The following assumptions have been made when looking at the problem definition:

* A minimum of 2 aliens is necessary for a meaningful simulation.
* There will only be a single contiguous map of cities (i.e. there won't be
  patches of unconnected cities littered throughout the map).
* The map is visualised as a grid, where cities occupy single squares within
  that grid.
* Necessarily, it will be possible for there to be empty squares within the map
  where there will be no cities. This is to be considered **lava**, where the
  aliens will always tend towards self-preservation and not venture into these
  empty squares.
* The city map cannot contain contradictory data, otherwise the program must
  indicate this as an input data error. For example, if city A is west of city
  B, then city C cannot simultaneously be west of city B.
* Aliens cannot travel diagonally on the map (e.g. North-West, or South-East).
  Thus they will only move in one of the four primary directions: North, East,
  South and West.
* If an alien is not yet trapped (i.e. it can move in at least one of the four
  primary directions), the random direction it will choose will be one of the
  **possible** directions. In other words, an alien will never even consider
  moving in a direction in which it will be blocked, unless it has absolutely
  nowhere to go.
* When aliens are randomly placed across the map, they are placed on cities and
  not on empty spots on the map.
* The stop criterion for the program of 10,000 moves per alien should rather be
  considered as 10,000 **potential moves**, otherwise trapped aliens could
  result in the program running forever.

