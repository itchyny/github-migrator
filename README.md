# github-migrator
This tool migrates a GitHub repository to another.
This is especially useful to move a repository from GitHub Enterprise to github.com.

## Basic usage
```bash
export GITHUB_MIGRATOR_SOURCE_API_TOKEN=xxx
export GITHUB_MIGRATOR_SOURCE_API_ENDPOINT=http://localhost/api/v3 # This might be the endpoint of GitHub Enterprise
export GITHUB_MIGRATOR_TARGET_API_TOKEN=yyy
go run . [owner]/[source] [owner]/[target]
```

## Disclaimer
This tool is stil under construction.
I assume no responsibility according to what happens using this tool.

## Bug Tracker
Report bug at [Issues・itchyny/github-migrator - GitHub](https://github.com/itchyny/github-migrator/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
