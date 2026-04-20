# Todo List

## Feature: Locate

- [ ] Validate args don't do anything dangerous
    - [x] command seems fine, we only ever do a comparison operator
        - [x] remove nested length checks and enforce correct number of arguments for each command
    - [x] process name not sure
        - [x] strip out endings of `.exe`
        - [x] also strip "path characters" and "control characters"
    - [ ] x, y, height, width should be fine b/c `strconv`?
        - [x] set a max int size so it doesn't wrap around
        - [x] ensure they are non-negative
        - [ ] set maxima based on display size; how do we get display size?
- [ ] Enable use of "instance" arg
    - Can getHWND return a list of matches along with their HWND's?
    - I'm hoping this list is already sorted in a sane way, so I can use the index as the "instance" number, which translates to "smallest number is the window opened the longest ago" (same order as when you hover over the taskbar and it shows previews)
- [ ] How to unminimize it if it's minimized?
    - [ ] Must do this before moving/resizing, nothing happens if it's minimized