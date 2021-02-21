# Extra Types

The current package contains types that I require in projects where the type never
changes, and used globally.


## Current Type Support

 * Duration - Ability to store `time.Duration` over JSON and database.
 * Numeric values - Ability to store and load `int` and `uint` family even when they are string for example.
 * Bool - Ability to take boolean value as int, string and boolean and convert to `bool` type, with `nil` support.


# TODO
  - [x] Add Tests for nil duration
  - [ ] Add more test covers for nil duration
  - [ ] Add More int and uint type support
