# Sesh

A nice TUI for managing your TMUX sessions


# Install

1. Install sesh
```bash
go install github.com/xemotrix/sesh
```
2. Set up key bindings in your `tmux.conf`

For example, if all your repos are in `~/repos`
```
bind-key i run-shell "tmux neww 'sesh switch -p ~/repos'" 
bind-key u run-shell "tmux neww 'sesh create -p ~/repos'" 
bind-key k run-shell "tmux neww 'sesh kill'" 
```

# Available commands

## sesh switch
`sesh switch -p ~/repos` will output a list of the currently active sessions + all directories in `~/repos`. 

Start typing to fuzzy-search them. If you want you can also navigate them using `ctrl+j` and `ctrl+k`. 

Once the one you want is highligted, press `enter` to:
- if it was an existing session switch to it.
- if it was a directory create a new session in that directory and switch to it.

## sesh create
`sesh create -p ~/repos` will promt a text input. There you can type the name of a new project. Feedback is provided on update about the validity of the name as a tmux session (all valid session names match the regex `^[A-Za-z](\w|-)*$`). If the name is valid, pressing enter will create a new directory in `~/repos`, a new session based on that directory, and finally switch to that session.

## sesh kill
`sesh kill` will output a list of the currently active sessions. There you can mark the ones you want to kill with `space` and confirm with `enter` to kill them.

# TODO commands
## sesh clone
Same as `sesh create` but given a git url, clone the repo, create the session there and switch to it.


