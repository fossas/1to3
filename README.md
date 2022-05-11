# 1to3

Utilities to help migrate from [FOSSA CLI](https://github.com/fossas/fossa-cli) v1 to v3.

## Installation

For now, you will need Go 1.17 installed on your machine.

```shell
git clone git@github.com:fossas/1to3.git
cd 1to3
go run main.go
```

## Usage

```shell
# Example directory
$ tree
.
├── my-test-directory
│   ├── package.json
│   └── requirements.txt
└── src
    └── Gemfile

2 directories, 3 files

# Targets detected by CLI v3
$ fossa list-targets
[ INFO] Found project: bundler@src/
[ INFO] Found target: bundler@src/
[ INFO] Found project: npm@my-test-directory/
[ INFO] Found target: npm@my-test-directory/
[ INFO] Found project: setuptools@my-test-directory/
[ INFO] Found target: setuptools@my-test-directory/

# Targets detected by CLI v1
$ fossa1.1.6 init
WARNING Filtering out suspicious module: my-test-directory (my-test-directory)
WARNING Filtering out suspicious module: my-test-directory (my-test-directory)

# Generated CLI v3 configuration block to exclude targets implicitly excluded by v1
$ fossa-1to3 targets
targets:
  exclude:
  - type: npm
    path: my-test-directory
  - type: pip
    path: my-test-directory
```
