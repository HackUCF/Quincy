# api/routes/competition

HTTP handlers for operator control over the running competition. These routes are the interface an organizer uses to halt and resume the competition mid-event, and to inspect its current state. They are not consumed by agents.

The package distinguishes two independent kinds of pause. A scoring pause leaves agents working — checks still run and the scoreboard still reflects the most recent result for every service — but results stop contributing to team pass and total counters, so no team loses uptime for the duration. A check pause goes further and stops the work itself: agents are handed a no-op instead of a real assignment and idle until the pause lifts. Each kind is held and toggled separately, so an operator can stop scoring without stopping the checks, or stop both.

Toggling is handled by a single parameterized handler shared across all four pause and unpause routes, with the desired state and pause kind bound at registration time rather than read from the request, so there is no request body to validate and no way for a caller to name a state the router does not already expose. A toggle that would be a no-op is reported rather than silently accepted: requesting a state the competition is already in aborts with 418 I'm a Teapot and leaves the stored state untouched. A real change returns 200 with a confirmation, and a database failure returns 500.

A separate read-only handler reports the full picture in one response — both pause kinds and, for each, the moment it entered its current state. Because all pause state is stored in the database rather than held in memory, every handler in this package requires the database sink and is replaced at startup with a 501 Not Implemented handler when none is configured.
