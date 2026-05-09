# Place

**Put your Windows windows in their place**

Partly a Go experiment, partly trying to solve my docking-station woes  
Letting Codex do the DLL calls, learning that is currently out-of-scope  


## The Idea

Have a place for every window, and be able to put them in their place


### First, a basic way to move windows

- `place firefox is` prints the current x, y, width, height of firefox :ballot_box_with_check:
- `place firefox 200 10 1200 1300` moves firefox to position x=200, y=10, width=1200, height=1300 :ballot_box_with_check:


### Then later, a robust way to keep things where they belong

- `place create desktop` creates a new location that can have its own unique window positions :ballot_box_with_check:
- `place select desktop` sets the current location :ballot_box_with_check:
- `place save firefox` saves the current position of firefox to the currently selected location, for later placing
- `place firefox` puts firefox into its saved position for the currently selected location
- `place all` puts all windows which have a saved position within the currently selected location into their place
- `place list locations` list all saved locations
- `place list positions` list all saved window positions for the current location

### And finally, a way to place multiple instances of one executable's windows

any command which takes an executable name can optionally take an instance number as a next parameter
- `place firefox 2 200 10 1200 1300` locates the second instance of firefox
- `place save firefox 2` saves the current position of the second instance of firefox
- `place firefox 2` places the second instance of firefox

also, be able to copy existing locations when creating new ones
- `place create desktop copy basic` creates a new location, copying all window positions from a location named "basic"

## Installation

- Run `go build`
- Add the exe to a folder on your PATH
    - Recommend `C:\Users\[username]\AppData\Local\Place` then add that to PATH