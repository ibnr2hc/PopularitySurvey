# PopularitySurvey

## Table of contents

- [Overview](#overview)
- [System Requirements](#system-requirements)
- [Installation](#installation)
- [Uninstallation](#uninstallation)
- [User Guid](#user-guid)
  - [Basics](#basics)

## Overview

Ranks the followers of the specified user in order of popularity.
Retrieve the number of followers of a follower and display them in order of the number of followers.


## System Requirements

| System | Version |
| ------ | ------- |
| go | 1.16.15 darwin/arm64 |


## Installation
1. Install command
   ```bash
   $ make install
   ```
2. Setup environment variable
   - TWITTER_BEARER_TOKEN(`export TWITTER_BEARER_TOKEN=XXX`)
     - Twitter API Bearer Token.


## Uninstallation
1. Uninstall command
   ```bash
   $ make uninstall
   ```


## User Guid

### Basics
```bash
$ popularity_survey survey -s {{ Survey target screen name (e.g. ibnr2hc) }}
```
![f24b4dd7-677b-4c8d-98c7-32d6ab594733](https://user-images.githubusercontent.com/43519598/213014643-1abface5-685f-4313-b27f-8ab7f1979be2.jpeg)

