# db/misc

Miscellaneous database utility functions that don't belong in a more specific subpackage.

## Files

- **duration.go** - `GetCompDuration()` returns how long the competition has been running by querying the first and last timestamps in the `scores` table.
