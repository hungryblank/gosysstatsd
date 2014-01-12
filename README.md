# goSysStatsd

## Simple utility to import system metrics in statsd

goSysStatsd is a utility to add system metrics like memory usage and
disk usage to your statsd system.

Being written in golang goSysStatsd is a zero dependency tool that can
just be deployed by copying one file on the destination system.

## Installation

1. copy the executable on your system

## Usage

every time the command is executed it will update the specified statsd
server with disk and memory usage

```sh
# with default options statsd listening on localhost port 8125
gosystatsd

# custom host and port
gosystatsd -h myhost -p 9999
```
