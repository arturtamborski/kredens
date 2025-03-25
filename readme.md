# kredens

Simple and fast credential manager for your environment variables.

## Install

```bash
brew install arturtamborski/tap/kredens
```

## Usage Demo

```bash
# Store some credentials
$ kredens set AWS_KEY AKIAXXXXXXX
$ kredens set AWS_SECRET xxxxxxxxxxx
$ kredens set GITHUB_TOKEN ghp_xxxxxxxx

# List all keys
$ kredens keys
AWS_KEY
AWS_SECRET
GITHUB_TOKEN

# Get a specific value
$ kredens get AWS_KEY
AKIAXXXXXXX
# Or shorter:
$ kredens AWS_KEY
AKIAXXXXXXX

# List all credentials
$ kredens list
AWS_KEY=AKIAXXXXXXX
AWS_SECRET=xxxxxxxxxxx
GITHUB_TOKEN=ghp_xxxxxxxx

# Source in your shell
$ eval $(kredens source)

# Remove a credential
$ kredens del AWS_KEY
```

## Features

- Pure Go implementation
- Stores data in SQLite database (`~/.kredens.db`, not configurable)
- Fast and simple interface
- Homebrew installation

## Usage

```
kredens [command] [args...]

Commands:
  list          List all credentials
  keys          List all keys
  vals          List all values
  get KEY       Show value for KEY
  set KEY VAL   Store KEY with value VAL
  del KEY       Delete KEY
  source        Output credentials as export statements
  help          Show this help message
```
