# Alien Invasion

## Overview
This program attempts to solve the **Alien Invasion** coding puzzle. More
details regarding the problem will be added to this README soon.

## Assumptions
The following assumptions have been made when looking at the problem definition:

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

