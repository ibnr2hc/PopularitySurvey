# PopularitySurvey

## Table of contents

- [Overview](#overview)
- [Installation](#installation)
- [Uninstallation](#uninstallation)
- [User Guid](#user-guid)
  - [Basics](#basics)

## Overview

Ranks the followers of the specified user in order of popularity.
Retrieve the number of followers of a follower and display them in order of the number of followers.


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
