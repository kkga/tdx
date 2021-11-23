# tdx

`tdx` is a todo manager for [iCalendar] files.

It is designed to work with [vdirsyncer] and expects to operate in its
[storage format][vdirstorage].

iCalendar-compatibility means it can be used as a CLI companion to any
CalDAV-enabled todo program, such as iOS Reminders. See more on how to set up
calendars and synchronization in [vdirsyncer documentation][vdirdocs].

[iCalendar]: https://en.wikipedia.org/wiki/ICalendar
[vdirsyncer]: https://github.com/pimutils/vdirsyncer
[vdirstorage]: https://vdirsyncer.pimutils.org/en/latest/vdir.html
[vdirdocs]: https://vdirsyncer.pimutils.org/en/stable/index.html

## Features

- adding todos
  - automatic date and priority parsing
- listing todos
  - sorting and filtering by fields
  - automatic hashtag parsing and output organized by tags
- completing todos
- editing todos in a `$VISUAL`/`$EDITOR` program
- deleting todos
- purging completed/cancelled todos

## Usage

### General usage

```
tdx -- todo manager for vdir (iCalendar) files.

Usage:
  tdx [flags]
  tdx [command]

Available Commands:
  add         Add todo
  list        List todos
  done        Complete todos
  edit        Edit todo
  show        Show todos
  delete      Delete todos
  purge       Delete done todos
  help        Help about any command
  completion  generate the autocompletion script for the specified shell

Flags:
  -h, --help          help for tdx
  -p, --path string   path to vdir folder
  -v, --version       version for tdx

Use "tdx [command] --help" for more information about a command.
```

### List command usage

```
List todos, optionally filtered by query.

Usage:
  tdx list [query] [flags]

Aliases:
  list, ls, l

Examples:
$ tdx list --sort prio --due 2

Flags:
  -l, --lists LISTS     filter by LISTS, comma-separated (e.g. 'tasks,other')
  -g, --group string    group listed todos, valid options: list, tag, none  (default "list")
  -a, --all             show todos from all lists (overrides -l)
  -d, --due N           filter by due date in next N days
  -S, --status STATUS   filter by STATUS: needs-action, completed, cancelled, any (default "needs-action")
  -t, --tag TAGS        filter todos by given TAGS
  -T, --no-tag TAGS     exclude todos with given TAGS
  -s, --sort FIELD      sort by FIELD: prio, due, status, created (default "prio")
      --description     show description in output
      --two-line        use 2-line output for dates and description
  -h, --help            help for list

Global Flags:
  -p, --path string   path to vdir folder
```

## Installation

### From release binaries

Download the compiled binary for your system from
[Releases](https://github.com/kkga/tdx/releases) page and put it somewhere in
`$PATH`.

### From source

Requires [Go](https://golang.org/) installed on your system.

Clone the repository and run `go build`, then copy the compiled binary somewhere
in `$PATH`.

If Go is [configured](https://golang.org/ref/mod#go-install) to install packages
in `$PATH`, it's also possible to install without cloning the repository:

```
go install github.com/kkga/tdx@latest
```

## Configuration

`tdx` is configured through environment variables.

| variable        | function                                                      |
| --------------- | ------------------------------------------------------------- |
| `TDX_PATH`      | Path to [vdir] directory[^fn1]                                |
| `TDX_LIST_OPTS` | Default options for `<list>` command, see `tdx list -h`[^fn2] |
| `TDX_ADD_OPTS`  | Default options for `<add>` command, see `tdx add -h`[^fn3]   |
| `NO_COLOR`      | Disable color in output                                       |

[^fn1]: Either root path containing multiple collections or path to specific
collection containing `*.ics` files.

[^fn2]: For example, to show todos due in the next 2 days, from 'myList',
grouped by tag: `TDX_LIST_OPTS='-d 2 -l myList -g tag'`

[^fn3]: For example, to use a default list for new todos:
`TDX_ADD_OPTS='-l myList'`

[vdir]: http://vdirsyncer.pimutils.org/en/stable/vdir.html
