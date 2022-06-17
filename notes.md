### Remaining work for a basic port to Go

#### GCS-specific work that needs to be done

- Add undo records for edit operations that don't already have them
- Implement prompting for substitution text when moving items onto a sheet
- Add monitoring of the library directories for file changes
  - Perhaps also add manual refresh option, for those platforms where disk monitoring is less than optimal
- Settings editors
  - Attributes
  - Body Type
- Library configuration dialogs
- Completion of menu item actions
  - Item
    - Copy to Character Sheet
    - Copy to Template
    - Apply Template to Character Sheet
  - Library
    - Update <library> to <version>
    - Change Library Locations
  - Settings
    - Attributes...
    - Default Attributes...
    - Body Type...
    - Default Body Type...
- Printing support for sheets (requires support in unison first)

#### Unison-specific work that needs to be done

- Printing support
- Carefully comb over the interface and identify areas where things aren't working well on Windows and Linux, since I
  spend nearly all of my development time on macOS and may not have noticed deficiencies on the other platforms
