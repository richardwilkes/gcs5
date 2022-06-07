### Remaining work for a basic port to Go

#### GCS-specific work that needs to be done

- Add undo records for edit operations that don't already have them
- Settings editors
  - Attributes
  - Body Type
- Library configuration dialogs
- Completion of menu item actions
  - Edit
    - Duplicate
  - Item
    - Copy to Character Sheet
    - Copy to Template
    - Apply Template to Character Sheet
  - Library
    - Show <library> on Disk
    - Update <library> to <version>
    - Change Library Locations
  - Settings
    - Attributes...
    - Default Attributes...
    - Body Type...
    - Default Body Type...
  - Help
    - Check for GCS updates...
- Printing support for sheets (requires support in unison first)

#### Unison-specific work that needs to be done

- Printing support
- Add keyboard navigation to the Table object
- Drag & drop support in the Table object
- Carefully comb over the interface and identify areas where things aren't working well on Windows and Linux, since I
  spend nearly all of my development time on macOS and may have not noticed deficiencies there
