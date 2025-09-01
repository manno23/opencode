#!/usr/bin/fish
# #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  # #
#  Commands:
#    opencode [project]         start opencode tui                        [default]
#    opencode run [message..]   run opencode with a message
#    opencode auth              manage credentials
#    opencode agent             manage agents
#    opencode upgrade [target]  upgrade opencode to the latest or a specific
#                               version
#    opencode serve             starts a headless opencode server
#    opencode models            list all available models
#    opencode github            manage GitHub agent
#
#  Positionals:
#    project  path to start opencode in                                    [string]
#
#  Options:
#        --help        show help                                          [boolean]
#    -v, --version     show version number                                [boolean]
#        --print-logs  print logs to stderr                               [boolean]
#        --log-level   log level
#                              [string] [choices: "DEBUG", "INFO", "WARN", "ERROR"]
#    -m, --model       model to use in the format of provider/model        [string]
#    -c, --continue    continue the last session                          [boolean]
#    -s, --session     session id to continue                              [string]
#    -p, --prompt      prompt to use                                       [string]
#        --agent       agent to use                                        [string]
#        --port        port to listen on                      [number] [default: 0]
#    -h, --hostname    hostname to listen on        [string] [default: "127.0.0.1"]


# Start server in background and remember its PID (make it global so handlers can see it)
bun run dev serve --hostname 0.0.0.0 --port 4096 --log-level INFO --print-logs >.server-logs 2>&1 &
set -g __opencode_server_pid $last_pid

# Kill the server if the user hits Ctrl-C / closes the terminal
function __opencode_cleanup --on-signal INT --on-signal TERM --on-signal HUP --on-event fish_exit
    if set -q __opencode_server_pid
        echo "Exiting, killing process $__opencode_server_pid"
        kill -TERM $__opencode_server_pid 2>/dev/null
        sleep 0.2
        kill -KILL $__opencode_server_pid 2>/dev/null
        set -e __opencode_server_pid
    end
    functions -e __opencode_cleanup
end

# Your app logic
set -l WORKINGDIR "$HOME/data/code/opencode"
cd $WORKINGDIR

set -l TUI_CLIENT "packages/tui/opencode"
set -gx OPENCODE_SERVER "http://0.0.0.0:4096"
set -gx OPENCODE_APP_INFO (printf '{"hostname":"localhost","git":false,"path":{"home":"%s","config":"%s/.config/opencode","data":"%s/.local/share/opencode","root":"%s","cwd":"%s","state":"%s/.local/state/opencode"},"time":{}}' $HOME $HOME $HOME $WORKINGDIR $WORKINGDIR $HOME)

if test -x $TUI_CLIENT
    $TUI_CLIENT
    set -l tui_status $status
else
    echo "$TUI_CLIENT doesn't exist."
    set -l tui_status 127
end

# Normal-path cleanup (also happens if TUI exits normally)
if set -q __opencode_server_pid
    kill -TERM $__opencode_server_pid 2>/dev/null
    wait $__opencode_server_pid 2>/dev/null
    set -e __opencode_server_pid
end
functions -e __opencode_cleanup

return $tui_status
