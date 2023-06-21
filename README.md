# arkectl

## Install

Open your terminal and install the CLI with GoLang:

```
go install github.com/arkemishub/arkectl@arkectl
```

or with Homebrew:

```
brew tap arkemishub/homebrew-arkectl
brew install arkectl
```

## Setup

Create a new environment variable within your `.zshrc` or `.bashrc` file.

```bash
export ARKEPATH=$HOME/arke
```

## Usage

### Install arkectl

```bash
arkectl install
```

### Create a new project

```bash
arkectl create-app <project-name>
```

Flags:

- `--interactive` - Starts the interactive mode of the command
- `--local` - Clones arke repositories locally instead of using Docker.
- `--frontend` - Creates only the frontend app. Works only with --local flag
- `--console` - Clones only the console repository. Works only with --local flag
- `--backend` - Creates only the backend app. Works only with --local flag. Under development 🚧

### Start development server

```bash
arkectl start
```

### Update docker images

```bash
arkectl update
```
