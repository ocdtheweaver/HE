# HE ↔ WolfHead OS — Integration Guide
*How HE programs interact with the WolfHead operating system*
*Status: Stubs implemented. Real compositor bindings: Pass 8 target.*

---

## Overview

WolfHead OS is a Linux-based, spatial, context-driven operating system where HE is the first-class application language. Every WolfHead app is a `.he` file. The OS itself is configured in HE. The shell is written in HE.

This guide covers:
- How HE programs access WolfHead features today (stub layer)
- What the real binding architecture will look like
- How to write WolfHead apps in HE now

---

## The `wolfhead` Module

```
summon "wolfhead" as wh
~ Also available as: ~
summon "os" as os
```

---

## Workspaces

WolfHead uses a spatial workspace model. Workspace 0 is **Genesis** — the launcher and home. Numbered workspaces hold apps.

### Current workspace

```
set ws to wh.workspace()
say "I am in: {ws.name}"
say "Workspace ID: {ws.id}"
say "Active: {ws.active}"
```

### Switch workspace

```
~ Switch to workspace 3 ~
wh.workspace(3)

~ Switch to Genesis ~
wh.workspace(0)
```

### Workspace lifecycle (planned)

```
~ Future API ~
set ws to wh.workspace(2)
tell ws to pause      ~ suspend background processing ~
tell ws to resume     ~ restore full processing ~
tell ws to close      ~ terminate workspace ~
```

### Z-axis navigation (planned)

WolfHead uses a three-dimensional workspace model: left/right for workspaces, in/out for zoom levels.

```
~ Future API ~
wh.zoom(1)    ~ zoom in (focus) ~
wh.zoom(-1)   ~ zoom out (overview) ~
```

---

## Contexts

Contexts are named profiles that change how the OS behaves: which apps are visible, what notifications are allowed, how CPU is allocated.

Built-in contexts: `Work`, `Social`, `Chill`. Custom contexts can be defined.

### Get available contexts

```
set contexts to wh.context()
for each ctx in contexts [
  say "Context: {ctx}"
]
```

### Switch context

```
wh.context("Work")
wh.context("Chill")
wh.context("Gaming")  ~ custom ~
```

### Context-aware app behaviour

```
create MyApp [
  has: [
    context is "Work"
  ]
  on context_change [
    set context to wh.context()
    if context is "Chill" then [
      say "Relaxing mode — lowering activity"
    ]
  ]
]
```

---

## Notifications

```
~ Simple notification ~
wh.notify("Meeting in 5 minutes")

~ Notification with title ~
wh.notify("Calendar", "Meeting in 5 minutes")

~ Planned: priority, actions ~
~ wh.notify("Alert", "Disk full", priority is "high") ~
```

---

## Launching Apps

```
wh.launch("browser")
wh.launch("terminal")
wh.launch("myapp.he")    ~ launch a HE app by file ~
```

---

## Gestures

WolfHead supports gesture input — drawing numbers with your finger or stylus switches workspaces. Apps can register gesture handlers.

```
~ Trigger a gesture programmatically ~
wh.gesture("swipe-right")
wh.gesture(3)              ~ draw the number 3 → switch to workspace 3 ~

~ Planned: register a reaction to a gesture ~
create MyApp [
  on gesture with "pinch" [
    say "Pinch detected!"
  ]
]
```

---

## Platform Detection

```
set plat to wh.platform()
say "Running on: {plat}"    ~ "WolfHead/Linux" ~

~ Cross-platform guard ~
if plat is "WolfHead/Linux" then [
  wh.notify("Hello from WolfHead!")
]
```

---

## Writing a WolfHead App

A complete WolfHead app in HE:

```
~ myapp.he — A simple counter app for WolfHead ~

summon "wolfhead" as wh

create CounterApp [
  has: [
    title is "Counter"
    count starts as 0
    context is "any"
  ]

  can: [
    start [
      wh.notify(title, "App started")
      say "{title} is running"
    ]
    increment [
      grow count by 1
      say "Count: {count}"
    ]
    reset [
      set count to 0
      wh.notify(title, "Counter reset")
    ]
    status [
      say "{title} — Count: {count} — Context: {context}"
    ]
  ]

  on click with "increment" [
    tell CounterApp to increment
  ]

  on click with "reset" [
    tell CounterApp to reset
  ]

  on context_change [
    set context to wh.context()
    tell CounterApp to status
  ]
]

tell CounterApp to start
```

---

## UI Integration

The `ui` module generates HTML output. In WolfHead, this will render to a Wayland surface.

```
summon "ui" as ui
summon "wolfhead" as wh

tell ui to window with [
  title is "My WolfHead App",
  children: [
    ui.text("Welcome to WolfHead"),
    ui.button("Open workspace 1", "open-workspace-1"),
    ui.button("Switch to Chill", "chill-mode")
  ]
]
```

Reaction handlers for UI events:

```
create App [
  on click with "open-workspace-1" [
    wh.workspace(1)
  ]
  on click with "chill-mode" [
    wh.context("Chill")
    wh.notify("App", "Switched to Chill mode")
  ]
]
```

---

## App Scoping (Planned)

WolfHead apps have three scopes:

| Scope    | Meaning                              | Declaration |
|----------|--------------------------------------|-------------|
| `system` | Runs in all contexts, all workspaces | `create App scope system [...]` |
| `global` | Runs in all contexts                 | `create App scope global [...]` |
| `user`   | Runs only in current context         | `create App [...]` (default) |

```
~ Future syntax ~
create StatusBar scope system [
  has: [
    time is ""
  ]
  on tick [
    summon "clock" as clock
    set time to clock.format(clock.now(), "HH:mm")
    say "Status: {time}"
  ]
]
```

---

## Physics Integration

```
summon "physics" as phys

create GameScene [
  can: [
    setup [
      phys.gravity(-9.8)
    ]
  ]
  on collision [
    say "Something collided!"
  ]
]
```

---


## WolfHead Apps with Method Chains and Closures

Pass 6 and 7 features make WolfHead app code cleaner:

```
summon "wolfhead" as wh
summon "clock" as clock

create StatusBar [
  has: [
    context is "Work"
    workspace is 0
    time is ""
  ]
  can: [
    refresh [
      set now to clock.now()
      time becomes "{now.hour}:{now.minute}"
      say "Status: {context} | WS {workspace} | {time}"
    ]
    switchContext(name) [
      context becomes name
      wh.context(name)
      wh.notify("Context", "Switched to {name}")
    ]
  ]
  on tick [
    tell StatusBar to refresh
  ]
]

~ Register callbacks using closures ~
set onWorkspaceChange to ability(wsId) [
  tell StatusBar to refresh
  set StatusBar.workspace to wsId
]

tell StatusBar to refresh
```

## Context-Aware App Logic

```
summon "wolfhead" as wh
summon "list" as lst

set workContexts to ["Work", "Social", "Chill"]

if wh.context() is one of workContexts then [
  say "Running in a known context"
]

for each ctx in workContexts [
  say "Available: {ctx}"
]

~ Switch context with validation ~
set requested to "Gaming"
if not (requested is one of workContexts) then [
  wh.notify("Warning", "{requested} is not a built-in context")
]
wh.context(requested)
```

## Workspace Field Iteration

```
summon "wolfhead" as wh
set ws to wh.workspace()

say "Current workspace properties:"
for each key, val in ws [
  say "  {key}: {val}"
]
```

## Background Architecture (When Fully Built)

### Current state (stubs)

All `wolfhead` module functions currently `fmt.Printf` their intent. No real OS calls are made. This is intentional — the language is being built first.

### Target architecture (Pass 8)

```
HE program
    │
    ▼ wh.workspace(3)
lang/eval/stdlib.go
    │
    ▼ IPC over Unix socket / D-Bus / Wayland protocol
WolfHead Compositor Daemon
    │
    ├── Wayland compositor (wlroots-based)
    ├── Workspace manager
    ├── Context engine
    ├── Gesture recogniser
    └── Notification daemon
```

### Wayland surface binding

When `tell ui to window` is called in a WolfHead context, the output will be piped to a Wayland client surface rather than stdout. The surface is managed by the compositor, which handles compositing, input routing, and lifecycle.

### HE as shell language

Long-term, WolfHead's configuration, startup scripts, and shell interactions are all HE:

```
~ ~/.wolfhead/startup.he ~

summon "wolfhead" as wh

~ Set up workspaces ~
wh.workspace(1)
wh.context("Work")

~ Launch startup apps ~
wh.launch("terminal")
wh.launch("browser")

~ Configure gestures ~
summon "io" as io
set config to io.read("~/.wolfhead/gestures.he")
```

---

## Quick Reference

```
summon "wolfhead" as wh

wh.workspace()         → current workspace object (.id, .name, .active)
wh.workspace(n)        → switch to workspace n
wh.context()           → list all context names
wh.context("Work")     → switch context
wh.notify("msg")       → OS notification
wh.notify("title","msg") → notification with title
wh.launch("app")       → launch app
wh.gesture("name")     → trigger gesture
wh.platform()          → "WolfHead/Linux"
```

---

*WolfHead OS is under active development alongside HE.*
*The binding layer (Pass 8) will replace these stubs with real compositor calls.*
