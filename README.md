# lightctl

A Linux backlight / LEDs controller similar to
[brightnessctl](https://github.com/Hummer12007/brightnessctl),
written in Go and featuring smooth brightness changes.


## Usage

```plain
Usage: lightctl [--device DEVICE] [--debug] <command> [<args>]

Options:
  --device DEVICE, -d DEVICE
  --debug, -D
  --help, -h             display this help and exit

Commands:
  increase
  decrease
  set
  list
```

### List

This command returns a list of controllable outputs (backlights / LEDs).
By specifying `-d DEVICE_NAME` to the other commands, one can
perform an action on a specific device.

```plain
$ lightctl list
intel_backlight
input3::capslock
input3::numlock
input3::scrolllock
phy0-led
platform::micmute
platform::mute
tpacpi::kbd_backlight
tpacpi::power
tpacpi::standby
tpacpi::thinklight
tpacpi::thinkvantage
```

### Increase

Increases the brightness of the specified `-d DEVICE`
by the specified value / percentage.

```plain
Usage: lightctl increase [--smooth-ms SMOOTH-MS] TO

Positional arguments:
  TO

Options:
  --smooth-ms SMOOTH-MS, -s SMOOTH-MS [default: 0]

Global options:
  --device DEVICE, -d DEVICE
  --debug, -D
  --help, -h             display this help and exit

```

**Example:**

```plain
lightctl increase -s 100 10%
```

This command smoothly increases (in 100 ms) the default 
device's brightness by `10%`. 
As an example, if the previous brightness percentage was 15%, this command
will increase the brightness to 25% by linearly increasing it
and completing the action in 100ms.


### Decrease

Same as increase, but reduces the brightness

```plain
$ lightctl decrease -h
Usage: lightctl decrease [--smooth-ms SMOOTH-MS] TO

Positional arguments:
  TO

Options:
  --smooth-ms SMOOTH-MS, -s SMOOTH-MS [default: 0]

Global options:
  --device DEVICE, -d DEVICE
  --debug, -D
  --help, -h             display this help and exit
```

### Set

Sets the brightness to the specified value / percentage

```plain
$ lightctl set -h
Usage: lightctl set [--smooth-ms SMOOTH-MS] TO

Positional arguments:
  TO

Options:
  --smooth-ms SMOOTH-MS, -s SMOOTH-MS [default: 0]

Global options:
  --device DEVICE, -d DEVICE
  --debug, -D
  --help, -h             display this help and exit
```

**Example:**

```plain
$ lightctl set -s 200 50%
```

Sets the brightness to 50% linearly in 200ms.

## Implementation

Under the hood `lightctl` uses two methods to change the brightness,
depending on which permission the user running the software has.

1. Direct write to `/sys/class/CLASS/DEVICE/brightness`
2. DBUS [`org.freedesktop.login1 SetBrightness()` method](https://www.freedesktop.org/software/systemd/man/org.freedesktop.login1.html)

The advantage of using the DBUS method is that the user doesn't need
any particular privilege to run the command since this is handled by
DBUS.

## Compiling

```bash
go build -o ~/go/bin/lightctl ./
~/go/bin/lightctl
```