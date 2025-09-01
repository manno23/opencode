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

function run_client_server --description 'Run dev server + TUI; clean up on exit'
    # Nuke any stale cleanup from previous runs
    functions -q __opencode_cleanup_exit; and functions -e __opencode_cleanup_exit

    # Start the dev server in background; detach from our tty for safety
    bun run dev serve --hostname 0.0.0.0 --port 4096 --log-level INFO --print-logs \
        </dev/null >>.server-logs 2>&1 &

    # Track the process group (fish gives each job its own PGID = leader PID)
    set -g __opencode_server_pgid $last_pid

    # Shell-exit safety net (does NOT hijack your TUI)
    function __opencode_cleanup_exit --on-event fish_exit
        if set -q __opencode_server_pgid
            # Kill entire group so helpers/watchers die too
            kill -TERM -- -$__opencode_server_pgid 2>/dev/null
            sleep 0.2
            kill -KILL -- -$__opencode_server_pgid 2>/dev/null
            set -e __opencode_server_pgid
        end
        functions -e __opencode_cleanup_exit
    end

    # ---- Your TUI ----
    set -l WORKINGDIR "$HOME/data/code/opencode"
    cd $WORKINGDIR

    set -l TUI_CLIENT "packages/tui/opencode"
    set -gx OPENCODE_SERVER "http://0.0.0.0:4096"
    set -gx OPENCODE_APP_INFO (printf \
        '{"hostname":"localhost","git":false,"path":{"home":"%s","config":"%s/.config/opencode","data":"%s/.local/share/opencode","root":"%s","cwd":"%s","state":"%s/.local/state/opencode"},"time":{}}' \
        $HOME $HOME $HOME $WORKINGDIR $WORKINGDIR $HOME)

    if test -x "$TUI_CLIENT"
        # Run directly so it fully owns the TTY (Ctrl-C goes to the TUI)
        "$TUI_CLIENT"
        set -l rc $status
    else
        echo "$TUI_CLIENT doesn't exist."
        set -l rc 127
    end

    # ---- Normal-path cleanup (fires after the TUI exits) ----
    if set -q __opencode_server_pgid
        kill -TERM -- -$__opencode_server_pgid 2>/dev/null
        wait $__opencode_server_pgid 2>/dev/null
        set -e __opencode_server_pgid
    end
    functions -e __opencode_cleanup_exit

    return $rc
end
