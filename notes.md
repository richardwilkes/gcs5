### Remaining work for a basic port to Go

#### Unison-specific work that needs to be done

- Printing support in unison
- Drag & drop support in the unison Table object
- Carefully comb over the interface and identify areas where things aren't working well on Windows and Linux, since I
  spend nearly all of my development time on macOS and may have not noticed deficiencies there

#### GCS-specific work that needs to be done

- Detail editors
  - Modifiers section
  - Melee Weapons section
  - Ranged Weapons section
- Settings editors
  - Attributes
  - Body Type
- Library configuration dialogs
- Prompt to save when closing a modified document
- General completion of menu item actions (many are currently placeholders)
  - File
    - Save
    - Save As...
    - Print...
    - Recent Files list
    - Export To list
  - Edit
    - Duplicate
    - Convert to Container
  - Item
    - New Advantage
    - New Advantage Container
    - New Advantage Modifier
    - New Advantage Modifier Container
    - Add Natural Attacks Advantage
    - New Skill
    - New Skill Container
    - New Technique
    - New Spell
    - New Spell Container
    - New Ritual Magic Spell
    - New Carried Equipment
    - New Carried Equipment Container
    - New Other Equipment
    - New Other Equipment Container
    - New Equipment Modifier
    - New Equipment Modifier Container
    - New Note
    - New Note Container
    - Copy to Character Sheet
    - Copy to Template
    - Apply Template to Character Sheet
    - Open Page Reference
    - Open Each Page Reference
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
- Make final decision on whether tables that have no content in the sheet should be hidden (complicates updates and
  makes it harder for users to discover they exist)
  - I'm currently thinking of just having a preference for making empty lists vanish when printing. That way they are
    discoverable in the UI, but don't clutter up the printed sheet if you don't want them to.
- When a sheet is scrolled, tooltips don't seem to be coming up in the correct spot for cells in the tables