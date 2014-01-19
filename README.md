# goSysStatsd

## Simple utility to import system metrics in statsd

goSysStatsd is a utility to add system metrics like memory usage and
disk usage to your statsd system.

Being written in golang goSysStatsd is a zero dependency tool that can
just be deployed by copying one file on the destination system.

## Installation

1. Copy the executable on your system
```
wget http://downloads.gofog.org/gosystatsd-0.0.1-x86_64
```
2. Use it


## Usage

every time the command is executed it will update the specified statsd
server with disk and memory usage

```sh
# with default options statsd listening on localhost port 8125
gosystatsd

# custom host and port
gosystatsd -h myhost -p 9999
```

you should start seeing data in statsd similar to this
```
   { 'system.memory.total': 7710,
     'system.memory.used': 3253,
     'system.memory.free': 4457,
     'system.memory.shared': 0,
     'system.memory.buffers': 68,
     'system.memory.cached': 984,
     'system.memory.available': 5509,
     'system.memory.usagePct': 29,
     'disk_usage.blocks.total.-dev-sda6': 19998104,
     'disk_usage.blocks.used.-dev-sda6': 6941112,
     'disk_usage.blocks.usagePct.-dev-sda6': 34
   }
```
