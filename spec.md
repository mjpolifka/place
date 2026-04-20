# Spec

## General

- [x] enforce correct number of arguments for given command
- [ ] powershell install script runs `go build` then moves `place.exe` to $HOME\AppData\Local\Microsoft\WindowsApps so it's on $PATH
    - [ ] checks that go is installed before running `go build`
    - [ ] verifies `place.exe` exists before trying to move it
- [ ] locations are saved to their own `[name].json` file
    - [ ] positions are the elements of the json
- [ ] selected location is saved to `config.json` file
- [ ] files live in `./.Place`


## Command: Locate

- [ ] `place locate process-name instance x y width height` should move the given instance of the process-name.exe window to location x, y and resize it to width x height
- [ ] gets a list of windows of `process-name.exe`
    - [x] as a start, get the first instance
    - Can getHWND return a list of matches along with their HWND's?
    - I'm hoping this list is already sorted in a sane way, so I can use the index as the "instance" number, which translates to "smallest number is the window opened the longest ago" (same order as when you hover over the taskbar and it shows previews)
- [x] moves the correct window to `x, y`
- [x] resizes the correct window to `width x height`
- [ ] validates args don't do anything dangerous
    - [x] process name
        - [x] strips out whitespace (this seems unnecessary)
        - [x] enforces not blank (also seems unnecessary)
        - [x] enforces no "path characters"
        - [x] enforces no "control characters"
        - [x] strips out endings of `.exe`
    - [ ] x, y, height, width
        - [x] send through `stringconv.Atoi`, will error on non-int
        - [x] enforce non-negative
        - [x] enforce max value so it doesn't wrap around
        - [ ] enforce maxima based on display size
    - [ ] instance
        - [ ] send through `stringconv.Atoi`, will error on non-int
        - [ ] enforce non-negative
        - [ ] enforce max value of `10`
- [x] unminimizes the window first
    - [x] must do this before moving/resizing, nothing happens if it's minimized
    - [x] similar issue when maximized, as soon as you move it the old size returns
    - restoring if not maximized or minimized doesn't seem to have any ill effects, but not 100% sure

## Command: Create

- [ ] `place create name` should create a new json file `name.json`
    - [ ] validates `name` input isn't dangerous
        - [ ] enforces no "path characters"
        - [ ] enforces no "control characters"
        - [ ] enforces no special characters, especially `.`
    - [ ] contains '{"name":"name", "positions":[]}'
    - [ ] updates `config.json` so `name` is selected
        - [ ] checks existence of `config.json` first
        - [ ] creates `config.json` if it doesn't exist

### Command: Copy

- Not MVP


## Command: Select

- [ ] `place select name` should change the currently selected location to `name`
    - [ ] validates `name` input isn't dangerous
        - [ ] enforces no "path characters"
        - [ ] enforces no "control characters"
        - [ ] enforces no special characters, especially `.`
    - [ ] checks if `name.json` exists
        - [ ] asks to create it if it doesn't exist
    - [ ] updates currently selected location in-memory
    - [ ] saves the selection to `config.json`


## Command: Save

- [ ] `place save process-name instance` should save the position of the given instance of the process-name.exe window to the current location
    - [ ] gets or already has the current location
        - [ ] asks to select a location if none is selected
    - [ ] validates `process name` isn't dangerous
        - [x] enforces no "path characters"
        - [x] enforces no "control characters"
        - [x] strips out endings of `.exe`
    - [ ] validates `instance` isn't dangerous
        - [ ] send through `stringconv.Atoi`, will error on non-int
        - [ ] enforce non-negative
        - [ ] enforce max value of `10`
    - [ ] gets the current size and location of the window
    - [ ] checks for an existing save for this process/instance
    - [ ] updates the values in-memory
    - [ ] saves the values to the correct json file


## Command: Process-Name

- [ ] `place process-name instance` puts the given instance of the process-name.exe window into its saved position in the current location
    - [ ] if `instance` is omitted, `1` is used
    - [ ] validates `process name` isn't dangerous
        - [x] enforces no "path characters"
        - [x] enforces no "control characters"
        - [x] strips out endings of `.exe`
    - [ ] validates `instance` isn't dangerous
        - [ ] send through `stringconv.Atoi`, will error on non-int
        - [ ] enforce non-negative
        - [ ] enforce max value of `10`
    - [ ] gets or already has the current location
        - [ ] asks to select a location if none is selected
    - [ ] checks for an existing save for this process/instance before attempting move
        - [ ] errors if no location is saved
    - [ ] moves the correct window to the correct location
    - [ ] resizes the correct window to the correct size
    


## Command: All

- [ ] `place all` puts all windows into their correct positions in the current location
    - [ ] gets or already has the current location
        - [ ] asks to select a location if none is selected
    - [ ] gets list of positions from json file
    - [ ] includes name, instance, x, y, width, height for each item
    - [ ] moves and resizes each item to match those values


## Command: List

- [ ] `place list` errors
    - [ ] reminds user to include either `locations` or `positions` sub-command

### Command: Locations

- [ ] `place list locations` lists all saved locations
    - [ ] lists all files in `./.Place` that end with `.json`
    - [ ] strips `.json` from all filenames
    - [ ] removes `config` from the list
    - [ ] displays list to user

### Command: Positions

- [ ] `place list positions` lists all saved positions in the current location
    - [ ] gets or already has the current location
        - [ ] asks to select a location if none is selected
    - [ ] gets list of positions from json file
    - [ ] includes name, instance, x, y, width, height for each item
    - [ ] displays list to user