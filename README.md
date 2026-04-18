# Nothing gets past my bow!

### [Read the docs!](/docs/README.md)

## Quick Start (are you kidding me?)

This scoreboard is not designed for a quick start, so this is about as fast as it gets.

```sh
# pull quincy
wget https://github.com/HackUCF/Quincy/releases/latest/download/quincy-linux-amd64 -O ./quincy
chmod +x ./quincy

# run setup
./quincy api dump-config
mkdir scripts
```

1. Edit `config.yaml` and define your competition ([guide](/docs/USAGE.md#configuration), [default config](/src/api/config/default-config.yaml))
2. Add corresponding scripts to your new `scripts/` directory ([guide](/docs/USAGE.md#check-scripts))

```sh
# in two different terminals
./quincy api start
PATH=$(pwd)/scripts:$PATH ./quincy agent start
```

For a longer explanation, see the [usage guide](/docs/USAGE.md).

## About

Quincy is a cybersecurity competition scoring engine purpose-built for [HPCC](https://plinko.horse), a beginner-level competition run by UCF students.

### Project Requirements
This <ins><b>is not</b></ins> meant to be an all-in-one or one-size-fits-all solution. You should probably use [Scoring Engine](https://github.com/scoringengine/scoringengine), [Quotient](https://github.com/dbaseqp/Quotient), or [Scorestack](https://github.com/scorestack/scorestack) (if you hate yourself).
- It needs to be:
  - Fast: Checks need to happen <i>very</i> often (&lt;10s) for a better competitor experience. Scores and service status need to be served quickly and efficiently from storage.
  - Distributed: The scorechecks cannot come from the same app or box that the website is served from. Scoring agents should be as independent as possible. Agents need to be droppable on random servers or containers.
  - Customizable: It needs to support hyper-customizable service checks. While generalized checks are helpful for quickly standing up a functioning scoreboard, they are not ideal for creating the best possible competition.
  - Simple: It needs to have very little overhead to get it up and running. The project needs to be maintainable into the future by a small team of student developers. Complex components or advanced architecture can't be involved.
  - Stable: Failure resistance needs to be taken seriously. The scoring engine needs to behave for the entire competition.

<!-- WIP: i can't figure out a good way to format/write this section -->
### Justifications for Strange Decisions (uhhh guys?)

| Decision | Reasoning | Solution |
|----------|-----------|----------|
| Quincy has no UI or UX | By using an off-the-shelf scoreboard, the competition ends up needing to be shaped around its behavior and functionality. Adding inject submission to the same application as service uptime monitoring seems like feature creep. | Only expose an HTTP REST API. The frontend (scoreboard, injects, password changes) needs to be implemented to explicity meet the needs of each competition. Want to run the competition entirely through a Discord bot? Go ahead. Quincy will handle the service uptime monitoring. |
| Quincy is unauthenticated | The API has no authentication. Every endpoint is exposed to anyone with network access. <br><br> This one might change in the future. | Limit network access to the API with reverse proxies, firewalls, or security groups. Assume trust between the agent and API. It sounds like this adds complexity, but it actually simplifies deployment. It does effectively require the frontend to perform server-side requests on behalf of the user. |
| Quincy has no usable service checks | 1. Realism: it helps that every check is as detailed and specific to the service usecase. <br> 2. Competitive integrity: simple HTTP GET requests don't test much and can be easily cheesed. <br> 3. Performance: custom score checks can be written with performance in mind. | *Only* support custom score checks, but be inclusive within that. Agents can run on all common OSes, and can use any executable as a scorecheck. Quick drafts can be written and use Python/Bash scripts. Detailed browser automation can be run with Selenium scripts. Heavily optimized score checks can be written in Rust. |


### Slop :)

1. AI is used for *some* but not all documentation, scripting, troubleshooting, and research.
2. AI was not (and should not be) used for the writing of any application source code.
