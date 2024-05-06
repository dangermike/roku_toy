# Roku Toy

I want to play with my roku, at least until I need glasses.

## Commands

* **device**
  * **list**: shows all devices by USN and URL. If set, alias is also shown
  * **alias**: Creates an alias for the given USN. These are stored in `~/.config/roku_toy/aliases`. Note that reusing a USN or name will overwrite previous aliases.
  * **unalias**: deletes a previously set alias by USN or name.
* **channel**
  * **list**: Show all channels on the provided device
  * **get**: Show the currently active channel on the provided device

## Links

* [external control API documentation](https://sdkdocs-archive.roku.com/External-Control-API_1611563.html)
* [SSDP spec](https://datatracker.ietf.org/doc/html/draft-cai-ssdp-v1-03)
* [fuzzy library](https://github.com/lithammer/fuzzysearch)
