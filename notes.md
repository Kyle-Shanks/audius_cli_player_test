# Notes and things

## Useful Things
- This is a super useful article about how to use commands
    - https://charm.sh/blog/commands-in-bubbletea/

## TODO
- Think of a clean key map for all of the things
- Set up a logging system that works better than printing to the console
    - Something like this: https://github.com/charmbracelet/bubbletea/blob/6b77c8fc10d43195ab29e6e09f93272623ce4e9c/logging.go
    - or this
    ```go
    if len(os.Getenv("DEBUG")) > 0 {
        f, err := tea.LogToFile("debug.log", "debug")
        if err != nil {
            fmt.Println("fatal:", err)
            os.Exit(1)
        }
        defer f.Close()
    }
    ```
- Update the help text to display things for search and stuff
- Add queue methods in player
    - Just need next and prev for now
- Update the main view names to page
- Update the keymap system to have the player and the pages handle all of the keymapping except for quit and toggle help
- Update the tracks table to save the column titles to dynamically get the info from the tracks
- Write a system to add the track Id to the filename when downloading and then only deleting all of the temp files on quit
    - Also check to see if we have the track file before playing to skip downloading again

## Goals
- User Favorite Tracks
- Way for user to save their userId for favorites
    - Also update favorites page to display 'No User Id' if none provided
- Global controls and Help text
    - Need to think of a control scheme that will be clean and intuitive
    - Top level help text for sure
    - Help text in specific places as needed

## Future Goals
- More mouse interactions when the UI is nailed down
- Work on sending messages out for tracking plays and stuff
- Collection View
- User Favorited Playlists
- Collection Search
- Queue View and more robust and managable queue system
- Shuffle
- App resizing
- User Login & Writes
    - Favoriting Tracks
- User Profile View
- User Search

