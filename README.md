# dwmon
`dwmon` is an extensible status bar printer for [dwm](https://dwm.suckless.org/) written in [Go](https://golang.org/).

It has number of plugins and it invokes those plugins to print the status string in specified format and to send it down one of the available sinks.

Plugins can listen for various events or unix signals and trigger status string re-generation. Also, plugins can schedule jobs that will trigger `dwmon` to regenerate status string periodically.

In order to use it, you need to specify:
 - list of plugins to enable
 - output format
 
This:
 
`dwmon --plugins=clock --format="{{ clock }}"`

will periodically print current date and time into `stdout` using "Mon, 02 Jan 2006, 03:04 PM" format.

In order to print it to DWM's status bar change the sink:

`dwmon --plugins=clock --sink=xsetroot --format="{{ clock }}"`

The format is normal [Go template](https://golang.org/pkg/text/template/) string where each plugin is a function.

# How to use

Add it into your `.xinitrc` and specify plugins to enable and format of the status string. Here is mine:

```
dwmon --plugins=wakeup,clock,shell --sink=xsetroot \
  --clock-interval=60 \
  --format="$(cat ~/.config/dwmon/format)" &
```

The format (`~/.config/dwmon/format`) is this:

```
{{shell "/bin/sh" "stat_mpd"}} | {{shell "/bin/sh" "stat_backlight"}} | {{shell "/bin/sh" "stat_vol"}} |{{shell "/bin/sh" "stat_net"}} |{{shell "/bin/sh"
"stat_bat"}} | {{clock}}
```

It outputs current song playing in mpd, level of the screen brightness, volume, status of my wifi, battery level and current date and time.

`--clock-interval` sent to `clock` plugin makes it repaint the status string every 60 seconds updating the clock.

Additionally, `wakeup` plugin allows me to ask `dwmon` to regenerate status string from outside of the process by sending `USR1` to it. I use it in my `sxhkd` configuration like that:

```
# volume
XF86Audio{RaiseVolume,LowerVolume,Mute}
  pamixer {-i 5,-d 5,--toggle-mute}; killall -USR1 dwmon
```

This changes the volume by calling to `pamixer` and then sends `USR1` signal to `dwmon` causing it to re-generate status string so I can see the updated volume value right away.

# Plugins

There are a short number of plugins supported right now but they provide all the extensibility necessary for basic use

- `clock` - provides a way to print current date and time
- `wakeup` - doesn't print anything but allows to send `USR1` signal to the process so it would re-generate the status string
- `echo` - prints whatever is sent to it
- `shell` - executes a shell command and prints output
- `mpd` - prints currently playing song from MPD and also listens for changes in the daemon.

Before plugins can be used, they must be enabled. You enable plugins by listing them in `--plugins` flag, like that: `--plugins=clock,wakeup`.

Once enabed a plugin function can be used in `--format` string multiple times.

Some plugin functions accept parameters. If a plugin function accepts parameters it accepts a variadic number of them. 

## Clock

A plugin to print current date and time. Accepts two parameters:

- `--clock-format` format for printing date and time. Default "Mon, 02 Jan 2006, 03:04 PM". Use [Go's date time layout syntax](https://golang.org/pkg/time/#Time.Format)
- `--clock-interval` how often to cause re-generation. *Note:* regenerating clock regenerates the whole status string.

**Format**

Use `clock` function in format. It accepts no parameters.

Example:
```
{{clock}}
```
will print current date and time in format specified by `--clock-format` or using the default one.

## Wakeup

A plugin that when enabled allows to send `USR1` signal to `dwmon` process causing it to re-generate status string.

**Format**

Function `wakeup` is exported as plugin system requires but it doesn't print anything.

## Echo

A plugin that prints whatever is supplied to it into status string. Useful for testing

**Format**

Use `echo` function in format. You can supply as many arguments to it as you want and it'll print all of them, comma separated.

Example:
```
{{echo "Hello" "World"}}
```
will print "Hello,World" into format.

## Shell

A plugin that prints output of a shell command into status string.

**Format**

Use function `shell` in format and supply it with the command you want to run. You can provide multiple arguments. Each of the arguments will be treated as an argument for the command executed by the plugin so make sure that the first argument is name of a binary or a script.

Example:

```
{{shell "echo" "Hello, world!"}}
```
will print "Hello, world!" into status string.

## MPD

A plugin that prints currently playing song in MPD.

Plugins is able to print the following properties:

 - `Title`
 - `Album`
 - `Artist`
 - `AlbumArtist`
 
**Format**

Use function `{{ mpd }}` order to print current song in requested format. Format specified using function parameters.
All of the unsupported parameters will be printed out verbatim.

Example:

```
{{ mpd "Title" " - " "Artist" }}
```

Prints title of the currently playing song, then it prints " - ", then it prints artist of the currently playing song.

# Sinks
dwmon supports two sinks at the moment of writing:
 - `stdout`
 - `xsetroot`
 
By default it uses sink `stdout` so it would be easier to try the tool and see the results.
 
Once the configuration is tested, sink can be switch to `xsetroot` with `--sink` option to output status string into DWM's status bar.
