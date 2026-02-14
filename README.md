# Nothing gets past my bow!

- To get started running Quincy see [USAGE.md](/USAGE.md)
- For information about the source code or development see [DEVELOPMENT.md](/DEVELOPMENT.md)
- The API specification is available here: [api/API_SPEC.md](/api/API_SPEC.md) (this is nowhere near final, just what has been implemented so far)

## About

- Quincy is a cybersecurity competition scoring engine built to be fast, customizable, and scalable.
- This <ins><b>is not</b></ins> meant to be an all-in-one or one-size-fits-all solution. You should probably use Quotient or Scorestack.
- This exists to meet the specific requirements of [HPCC](https://plinko.horse), a beginner competition run for UCF students by UCF students:
  - It needs to be fast. Checks need to happen <i>very</i> often (&lt;10s) for a better competitor experience. Scores and service status need to be served quickly and efficiently from storage.
  - It needs to be distributed. The scorechecks cannot come from the same app or box that the website is served from. Scoring agents should be as independent as possible. Agents need to be droppable on random servers or containers.
  - It needs to support hyper-customizable service checks. While generalized checks are helpful for quickly standing up a functioning scoreboard, they are not ideal for creating the best possible competition.
  - It cannot be container native. It's annoying to use on our private cloud, which only has managed servers.
- This project has some limitations:
  - This repository consists of an API and its scoring agent. The frontend should be implemented using the API to suit the specific needs of your competition.
  - The scoreboard, password change form, injects, and all other user facing components need to be implemented to suit the specific needs of your competition.
  - The API is completely unauthenticated:
    - The agents and API operate under assumed trust.
    - Anything with access to the API can post a score for any service, or change a users password for any team.
    - Access should be granted using a combination of firewalls and/or proxies.
    - The planned architecture for plinko will completely offload authentication and API interaction to the frontend. Users will have no network access to the API.
    - This may change in the future.

## Slop Warning

- AI is used for the some documentation, testing automation, troubleshooting, and research.
- AI was not and should not be used for the writing of any application source code.