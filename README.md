# github-migrator
This tool migrates a GitHub repository to another.
This is especially useful to move a repository from GitHub Enterprise to github.com.

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
- Pull requests
  - A pull request is converted to an issue
  - Comments and review comments are migrated as issue comments
  - Created dates, Labels
  - Pull request numbers (issue numbers) are same as the original repository
- Repository information
  - Description, Homepage (only when the target repository has blank description, homepage)
- Labels
  - Label name, description and colors
- All the other things will be lost
  - Issue and pull request reactions
  - Diffs (split) view of pull requests
  - Wiki
  - Projects, Milestones (will be implemented in the near future)
  - Default branch, Protection rules
  - Webhooks, Notifications, Integrations

## Disclaimer
This tool is still under construction.
I assume no responsibility according to what happens using this tool.

## Bug Tracker
Report bug at [Issues・itchyny/github-migrator - GitHub](https://github.com/itchyny/github-migrator/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
