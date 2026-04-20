# Place

**Put your windows in their place**

Partly a Go experiment, partly trying to solve my docking-station woes  
Letting Codex do the DLL calls, learning that is currently out-of-scope  


## The Idea

Have a place for every window, and be able to put them in their place


### First, a basic way to move windows

- `place locate firefox 1 200 10 1200 1300` moves the first instance of firefox to position x=200, y=10, width=1200, height=1300


### Then later, a robust way to keep things where they belong

- `place create desktop copy basic` creates a new location that can have unique window positions, copying all window positions from a location named "basic"
- `place save firefox 2` saves the current position of the second instance of firefox to the current location, for later placing
- `place select desktop` sets the current location
- `place firefox` puts the first instance of firefox into its saved position for the current location
- `place all` puts all windows which have a saved position within the "desktop" location into their place
- `place list locations` list all saved locations
- `place list positions` list all saved window positions for the current location