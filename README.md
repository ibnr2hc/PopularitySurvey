# PopularitySurvey

## Overview

Ranks the followers of the specified user in order of popularity.
Retrieve the number of followers of a follower and display them in order of the number of followers.


## Building
```bash
$ cd src
$ go build
```

## Setup
### Environment variable
- TWITTER_BEARER_TOKEN
  - Twitter API Bearer Token (`export TWITTER_BEARER_TOKEN=XXX`)

## User Guid

### Table of contents

* [Basics](#basics)

### Basics
```bash
$ ./PopularitySurvey survey -s {{ Survey target screen name (e.g. ibnr2hc) }}
```

