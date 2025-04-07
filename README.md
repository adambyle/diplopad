# Diplopad

**WORK IN PROGRESS**

Diplopad (**Diplo**macy Note**pad**) is an engine for *Diplomacy*, a strategy wargame by Allan B. Calhamer. It uses a smart, simple, and powerful adjudication system that is perfect for board analysis and simulation.

Some features that set Diplopad apart from other Diplomacy libraries:

* Diplopad has a simple and accessible API, useful both for basic scripts and for complex analysis and simulation tools.
* `Game` objects represent snapshot-like game *states*, instead of ever-transforming whole games. This allows applications to track a game's history or examine the outcome of multiple possible order sets on one game state.
* The `Arena` type allows for incremental, watchable adjudication. Orders can be added one-by-one, and the outcomes of hypothetical orders can be queried without applying them to the game. Arenas allow multiple possible order sets to be examined at once on the same underlying game state.
* Diplopad has built-in support for many useful functions, such as parsing orders from text and creating custom maps* to play the game on (so long as no extra mechanics are added with them).

In the future, Diplopad will support a higher-level game object that facilitates players' order submissions and press, like common web implementations of *Diplomacy*.

## *Custom map restrictions

Currently, Diplopad has some restrictions on creating maps whose provinces have complex coastal interactions. Namely:
1. If two adjacent provinces have identified coasts, one set of one province's coasts must connect to one set of the other province's coasts. But you cannot, for example, have ProvinceA SC - ProvinceB SC and ProvinceA NC - ProvinceB NC be distinct.

It is possible that these restrictions will be relaxed in the future.
