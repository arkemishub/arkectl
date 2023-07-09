# arkectl

## Install

Open your terminal and install the CLI with GoLang:

```
go install github.com/arkemishub/arkectl@latest
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

| Flag          | Shorthand | Description                                                               | Default       |
| ------------- | --------- | ------------------------------------------------------------------------- | ------------- |
| `interactive` | `i`       | Starts the interactive mode.                                              | false         |
| `local`       | `l`       | Clones arke repositories locally instead of using Docker.                 | false         |
| `frontend`    | `f`       | Creates only the frontend app. Works only with `--local` flag.            | false         |
| `backend`     | `b`       | Creates only the backend app. Works only with `--local` flag.             | false         |
| `console`     | `c`       | Clones only the console repository. Works only with `--local` flag.       | false         |
| `template`    | `t`       | Template to be used for the frontend app. Works only with `--local` flag. | `nextjs-base` |

### Start development server

```bash
arkectl start
```

### Update docker images

```bash
arkectl update
```
