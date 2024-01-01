# Sesh

A nice TUI for managing your TMUX sessions


# Install
1. Clone this repo
2. Build it
```bash
cd path/to/sesh 
go build -o ./bin/sesh .
```
3. Set up key bindings in your `tmux.conf`

For example, if all your repos are in `~/repos`
```
bind-key i run-shell "tmux neww 'path/to/sesh/bin/sesh switch ~/repos'" 
bind-key u run-shell "tmux neww 'path/to/sesh/bin/sesh create ~/repos'" 
```

# Available commands

## sesh switch
`sesh switch ~/repos` will output a list of the currently active sessions + all directories in `~/repos`. 

Start typing to fuzzy-search them. If you want you can also navigate them using `ctrl+j` and `ctrl+k`. 

Once the one you want is highligted, press `enter` to:
- if it was an existing session switch to it.
- if it was a directory create a new session in that directory and switch to it.

## sesh create
`sesh create ~/repos` will promt a text input. There you can type the name of a new project. Feedback is provided on update about the validity of the name as a tmux session (all valid session names match the regex `^[A-Za-z](\w|-)*$`). If the name is valid, pressing enter will create a new directory in `~/repos`, a new session based on that directory, and finally switch to that session.

# TODO commands
## sesh clone
Same as `sesh create` but given a git url, clone the repo, create the session there and switch to it.

## sesh kill
Kill one or many active sessions

