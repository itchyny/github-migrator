# github-migrator
This tool migrates a GitHub repository to another.
This is especially useful to move a repository from GitHub Enterprise to github.com.

![](https://user-images.githubusercontent.com/375258/71414326-cda1b480-2699-11ea-9de9-411e954bdb70.jpg)

## Usage
```bash
export GITHUB_MIGRATOR_SOURCE_API_TOKEN=xxx
export GITHUB_MIGRATOR_SOURCE_API_ENDPOINT=http://localhost/api/v3 # This might be the endpoint of GitHub Enterprise
export GITHUB_MIGRATOR_TARGET_API_TOKEN=yyy
# export GITHUB_MIGRATOR_TARGET_API_ENDPOINT=https://api.github.com # No need to specify the endpoint of github.com
go run . [old-owner]/[source] [new-owner]/[target]
```
Be sure to use this tool before pushing the git tree to the new origin (otherwise the links in the merged commits are lost).

### Other options
Sometimes same user has different user id on GitHub and Enterprise.
```bash
export GITHUB_MIGRATOR_USER_MAPPING=user-before1:user-after1,user-before2:user-after2,user-before3:user-after3
```

## Requirements
- Go 1.13+
- API tokens to access the source and target repositories.

## Features
- Issues
  - Issue description with the link to the original repository
  - Issue comments with the user name and icon (within the comment)
  - Created dates, Labels
  - Issue numbers are same as the original repository
  - Various events (including title changes, issue locking, assignments, review requests and branch deletion in a pull request)
- Pull requests
  - A pull request is converted to an issue
  - Comments and review comments are migrated as issue comments
  - Created dates, Labels
  - Pull request numbers (issue numbers) are same as the original repository
  - Number of changed files, insertions and deletions
  - Entire diff (excluding large file diffs)
  - Commits list and link to the corresponding /compare/ page
- Repository information
  - Description, Homepage (only when the target repository has blank description, homepage)
- Labels
  - Label name, description and colors
  - Label changes in issue and pull request
- Projects
  - Projects, columns and cards
  - Note that column automations are not migrated (cannot be set via API)
- Milestones
  - Milestone titles, description and due date
  - Connect issues to milestones
- Webhooks
  - Webhook URL, content type and events the hooks is trigger for.
- All the other things will be lost
  - Images posted to issue and pull request comments.
  - Emoji reactions to issue and pull request comments
  - Diffs (split) view of pull requests
  - Wiki
  - Default branch, Protection rules
  - Notifications, Integrations

## Bug Tracker
Report bug at [Issues・itchyny/github-migrator - GitHub](https://github.com/itchyny/github-migrator/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.

## Previous works and references
- [fastlane/monorepo: Scripts to migrate to a monorepo](https://github.com/fastlane/monorepo)
  - This tool greatly influenced me, especially for investigating the usage of the import api.
- [aereal/migrate-gh-repo: migrate GitHub (incl. Enterprise) repositories with idempotent-like manner](https://github.com/aereal/migrate-gh-repo)
  - For the idea of keeping the issue and pull request numbers.
- [Complete issue import API walkthrough with curl](https://gist.github.com/jonmagic/5282384165e0f86ef105)
  - Comprehensive tutorial for using the import api (which is not listed in the official api document yet).
