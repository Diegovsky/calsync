# Calsync

Calsync is a simple file-oriented tool to manage your Google Calendar events using a single source-of-truth file. It allows you to **create, edit, and manage events** in Google Calendar while keeping your personal calendar safe.  

Defining your whole schedule is as simple as:

```text
[Work]
11/10 - Meeting with Greg
12/10 - Super important project
20/10 - Getting fired

[Compilers] # University class
01/10/2025 - Lexical Analysis Workshop
15/10/2025 - Syntax Tree Presentation
28/10/2025 - Final Project Submission
````

Creates:
![assets/example.png](Screenshot of Google Calendar Showing)

---

## Overview

Calsync uses a single file, typically named `events.cal`, as the single source of truth for your events. This file is a custom format to describe events, grouped by tags.

Currently, it doesn't support specific times, only dates.

Calsync ensures that:

- **No existing events in your Google Calendar are modified**.  
- All events created by Calsync are stored under a dedicated `"Calsync Events"` calendar, so you can just delete it if you don't use it anymore.

> **Note on Google Calendar calendars:** In Google Calendar, a "calendar" is essentially a separate container for events. You can have multiple calendars (e.g., personal, work, or project-specific), and each calendar has its own set of events. Calsync creates and manages a custom calendar named `"Calsync"` to safely store all events it manages.

---

## `events.cal` Format

Calsync uses a simple DSL to describe events:

- `[<Tag>]` — Groups events under a common category or project (optional, isn't used for anything yet).  
- `DD/MM/YYYY - <Event Title> @<id>` — Describes an event:  
  - `DD/MM[/YYYY]` is the event date (year is optional, will be inferred as the current year if not specified)
  - `<Event Title>` is the name of the event  
  - `@<id>` is the Google Calendar ID assigned to the event. This is automatically added when the event is first created, and updated in the file on edits. 
- Anything after a `#` is considered a comment is currently deleted.

The `id` is very important for tracking changes to an event's name and date. It is also used to detect when an event is deleted from the file, ensuring it is deleted from the cloud calendar.

