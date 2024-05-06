# Roku Toy

I want to play with my roku, at least until I need glasses.

## Commands

* `device`
  * `list`: shows all devices by USN and URL. If set, alias is also shown
  * `alias`: Creates an alias for the given USN. These are stored in `~/.config/roku_toy/aliases`. Note that reusing a USN or name will overwrite previous aliases.
  * `unalias`: deletes a previously set alias by USN or name.
* `channel`
  * `list`: Show all channels on the provided device
  * `get`: Show the currently active channel on the provided device
  * `set`: Change to the provided channel by ID or name (fuzzy match)

### Usage notes

All of the `channel` commands have to target a single Roku. The target device can be specified using `--device` (`-d`) by alias or USN. You can also use `--first` (`-1`) to use the first device found on the network. The `--first` argument should not be used if you have more than one Roku on your network as there reporting order is not consistent. The commands will work but will be slower than if you provide `--device` or `--first` as the application has to wait for any straggler devices to report.

The "Home" application is not reported when listing applications in the Roku API. While it can be set if you know the ID, this application assumes that the ID is not known and will send the home key if `channel set home` or `channel set 0` is called.

Channel set by number does not query the applications from the device, but rather just passes the number. Setting the channel by name does fetch the channels, but has the advantage of fuzzy-matching the names. Calling `channel set apple` will launch "Apple TV" or `channel set plex` will launch "Plex - Free Movies &amp; TV".

## Examples

For these examples the actual USN of my Roku has been replaced with `ABCDEFGHIJKL`. I can specify that device with `-d ABCDEFGHIJKL` or `-d living_room` (after setting the alias). Also, since there is only one Roku on my network, I could use `-1`.

```bash
$ roku_toy device list                            # list devices on the network
ABCDEFGHIJKL http://192.168.1.176:8060/

$ roku_toy device alias ABCDEFGHIJKL living_room  # set alias

$ roku_toy device list                            # list devices on the network
ABCDEFGHIJKL http://192.168.1.176:8060/ living_room

$ roku_toy channel get -d living_room             # get current channel
YouTube (837)

$ roku_toy channel set -d living_room netflix     # set chanel to Netflix

$ roku_toy channel get -d living_room             # get current channel
Netflix (12)
```

## Links

* [external control API documentation](https://sdkdocs-archive.roku.com/External-Control-API_1611563.html)
* [SSDP spec](https://datatracker.ietf.org/doc/html/draft-cai-ssdp-v1-03)
* [fuzzy library](https://github.com/lithammer/fuzzysearch)
