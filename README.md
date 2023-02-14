<div align="center">
  <img width="128" src="./logo.svg" alt="JSON logo" />
  <h1>json-exec</h1>
  <p>A utility for capturing stderr and stdout output from a command and displaying them as JSON messages</p>
  <hr />
  <br />
  <a href="https://goreportcard.com/report/github.com/josh-hogle/json-exec" target="_blank">
    <img src="https://goreportcard.com/badge/github.com/josh-hogle/json-exec?style=for-the-badge" />
  </a>
  <a href="#">
    <img src="https://img.shields.io/badge/stability-alpha-ff69b4?style=for-the-badge" />
  </a>
  <a href="https://en.wikipedia.org/wiki/MIT_License" target="_blank">
    <img src="https://img.shields.io/badge/license-MIT-maroon?style=for-the-badge" />
  </a>
  <a href="#">
    <img src="https://img.shields.io/badge/support-community-purple?style=for-the-badge" />
  </a>
  <a href="https://conventionalcommits.org">
    <img src="https://img.shields.io/badge/Conventional%20Commits-1.0.0-orange.svg?style=for-the-badge" />
  </a>
</div>
<br />
<hr />
<br />

<!-- omit in toc -->
## Table of Contents
- [👁️ Overview](#️-overview)
- [✅ Requirements](#-requirements)
- [✴️ Installation](#️-installation)
- [▶️ Execution](#️-execution)
  - [➡️ run Command](#️-run-command)
  - [➡️ version Command](#️-version-command)
  - [➡️ Sample output messages](#️-sample-output-messages)
- [⛏️ Building from Source](#️-building-from-source)
- [📃 License](#-license)
- [❓ Questions, Issues and Feature Requests](#-questions-issues-and-feature-requests)

## 👁️ Overview

`json-exec` is a simple tool to execute any arbitrary command with a set of arguments and capture its output from stdout and/or stderr and produce a JSON message containing the result.

This utility was created in order to be able to produce consumable logs in JSON format for various services and executables.

Upon execution, `json-exec` will produce 2 log messages:
- Both messages will contain `@level`, `@timestamp` and `@message` fields as well as any additional fields passed via command-line flags.
- The first log message will contain `command` and `args` fields which will be populated with the command and an array of arguments passed to the command, respectively.
- The second log message will contain an `exit_code` field indicating the command's exit code. Additionally the `stdout` and `stderr` fields will contain output from stdout and stderr, respectively, if they are enabled.
- If the command produces and error, an `error_message` field will also be included in the second message.

## ✅ Requirements

This software is supported on the following platforms:
- Windows 10 or later (64-bit)
- MacOS (64-bit))
- Linux (64-bit)

There are no additional requirements to use this utility.

## ✴️ Installation

To install the utility on your system, simply download and unpack the appropriate archive file for your OS and then copy the `json-exec` binary to a folder on your path.

## ▶️ Execution

The general usage of `json-exec` is as follows:

```
Usage:
  json-exec [command]

Available Commands:
  help        Help about any command
  run         Executes an arbitrary system command with optional flags
  version     Display application version information

Flags:
  -c, --config-file string       Path to the configuration settings file
  -f, --field stringToString     one or more additional fields to include in the output (default [])
  -h, --help                     help for json-exec
      --level-field string       alternate name for the level field (default "@level")
  -l, --log-level string         adjust output log level - must be one of: debug, info, warn, error, fatal, panic or none (default "info")
      --message-field string     alternate name for the message field (default "@message")
      --timestamp-field string   alternate name for the timestamp field (default "@timestamp")

Use "json-exec [command] --help" for more information about a command.
```

### ➡️ run Command

The `run` command allows you to run an arbitrary command with or without arguments.

```
Usage:
  json-exec run [flags] <command> [command args]

Flags:
  -h, --help            help for run
      --ignore-stderr   ignore stderr output from the command
      --ignore-stdout   ignore stdout output from the command

Global Flags:
  -c, --config-file string       Path to the configuration settings file
  -f, --field stringToString     one or more additional fields to include in the output (default [])
      --level-field string       alternate name for the level field (default "@level")
  -l, --log-level string         adjust output log level - must be one of: debug, info, warn, error, fatal, panic or none (default "info")
      --message-field string     alternate name for the message field (default "@message")
      --timestamp-field string   alternate name for the timestamp field (default "@timestamp")
```

To run a simple system command without any arguments just call it as you normally would adding `json-exec run` to the start:

```
json-exec run chown root:root /root
```

If you are calling a system command that requires arguments, supply the `--` flag before the actual command so that the CLI does not interpret the command arguments as flags to `json-exec`:

```
json-exec run -- chown -hR root:root /root
```

To add extra fields to the output, use the `--field` flag repeatedly passing `key=value` for each pair:

```
json-exec run --field "@level=info" --field "@module=entrypoint" -- chown -hR root:root /root
json-exec run --field "app=myapp" -- chown -hR root:root /root
```

To ignore output from stderr, use the `--ignore-stderr` flag. To ignore output from `stdout`, use the `--ignore-stdout` flag.

```
json-exec run --ignore-stderr cat /etc/hosts
json-exec run --ignore-stdout curl -v https://google.com
```

### ➡️ version Command

The `version` command displays version information.

```
Usage:
  json-exec version [flags]

Flags:
  -h, --help        help for version
  -p, --plaintext   output plaintext instead of JSON
  -v, --verbose     display full version information including build and release date

Global Flags:
  -c, --config-file string       Path to the configuration settings file
  -f, --field stringToString     one or more additional fields to include in the output (default [])
      --level-field string       alternate name for the level field (default "@level")
  -l, --log-level string         adjust output log level - must be one of: debug, info, warn, error, fatal, panic or none (default "info")
      --message-field string     alternate name for the message field (default "@message")
      --timestamp-field string   alternate name for the timestamp field (default "@timestamp")
```

To display the version information in JSON output:

```
json-exec version
```

In this case only a single message will be printed and will include the `build`, `release_date` and `version` fields.

To display the version information in plaintext:

```
json-exec version --plaintext
```

### ➡️ Sample output messages

Attempting to change ownership on `/root` as a normal user:

```
{"@level":"info","@message":"Executing command: chown root:root /root","@timestamp":"2021-04-30T01:47:33.546988Z","args":["root:root","/root"],"command":"chown"}

{"@level":"warn","@message":"Command completed with exit code 1","@timestamp":"2021-04-30T01:47:33.547896Z","error_message":"exit status 1","exitCode":1,"stderr":"chown: changing ownership of '/root': Operation not permitted\n","stdout":""}
```

Attempting to recursively change ownership on `/root` as `root`:

```
{"@level":"info","@message":"Executing command: chown -hR root:root /root","@timestamp":"2021-04-30T01:49:42.127222Z","args":["-hR","root:root","/root"],"command":"chown"}

{"@level":"info","@message":"Command completed with exit code 0","@timestamp":"2021-04-30T01:49:42.142997Z","exitCode":0,"stderr":"","stdout":""}
```

`cat`'ing the `/etc/group` file (notice the newlines are replaced automatically with `\n` characters):

```
{"@level":"info","@message":"Executing command: cat /etc/hosts","@timestamp":"2021-04-30T01:52:45.072169Z","args":["/etc/hosts"],"command":"cat"}

{"@level":"info","@message":"Command completed with exit code 0","@timestamp":"2021-04-30T01:52:45.072848Z","exitCode":0,"stderr":"","stdout":"# This file was automatically generated by WSL. To stop automatic generation of this file, add the following entry to /etc/wsl.conf:\n# [network]\n# generateHosts = false\n127.0.0.1\tlocalhost\n127.0.1.1\tJOSH-DESKTOP.localdomain\tJOSH-DESKTOP\n\n# The following lines are desirable for IPv6 capable hosts\n::1     ip6-localhost ip6-loopback\nfe00::0 ip6-localnet\nff00::0 ip6-mcastprefix\nff02::1 ip6-allnodes\nff02::2 ip6-allrouters\n"}
```

Executing the `version` sub-command:

```
{"@level":"info","@message":"json-exec version 0.1.0 build abcdef (Released 29 Apr 2021)","@timestamp":"2021-04-30T01:42:54.222477Z","build":"abcdef","release_date":"29 Apr 2021","version":"0.1.0"}
```

Executing the `version` sub-command with the `--plaintext` option:

```
json-exec version 0.1.0 build abcdef (Released 29 Apr 2021)
```

## ⛏️ Building from Source

In order to build project from source, you will need the following software installed on your system:

- go 1.16 or later (<https://golang.org/dl/>)
- Standard build tools (eg: `git`, `make`, `tar` and `zip`)
- If building on a Windows machine, use WSL2

Once these tools have been installed, change to the root directory of the repository and run the command:

```
make clean all
```

Upon completion, the `dist` folder will contain `dev` and `release` subfolders with folders for each of the supported platforms.

## 📃 License

This utility is distributed under the MIT license.

## ❓ Questions, Issues and Feature Requests

If you have questions about this project, find a bug or wish to submit a feature request, please [submit an issue](https://github.com/josh-hogle/json-exec/issues).
