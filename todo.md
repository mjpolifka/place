# Todo List

## Feature: Locate

- Validate args don't do anything dangerous
    - command seems fine, we only ever do a comparison operator
        - remove nested length checks and enforce correct number of arguments for each command
    - process name not sure
        - strip out endings of `.exe`
    - x, y, height, width should be fine b/c `strconv`?
        - set a max int size so it doesn't wrap around
        - ensure they are non-negative
- Enable use of "instance" arg
    - Can getHWND return a list of matches along with their HWND's?
    - I'm hoping this list is already sorted in a sane way, so I can use the index as the "instance" number, which translates to "smallest number is the window opened the longest ago" (same order as when you hover over the taskbar and it shows previews)