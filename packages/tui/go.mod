module github.com/sst/opencode

go 1.24.0

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/alecthomas/chroma/v2 v2.18.0
	github.com/charmbracelet/bubbles/v2 v2.0.0-beta.1
	github.com/charmbracelet/bubbletea/v2 v2.0.0-beta.4
	github.com/charmbracelet/glamour v0.10.0
	github.com/charmbracelet/lipgloss/v2 v2.0.0-beta.3
	github.com/charmbracelet/x/ansi v0.9.3
	github.com/fsnotify/fsnotify v1.8.0
	github.com/google/uuid v1.6.0
	github.com/lithammer/fuzzysearch v1.1.8
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/mattn/go-runewidth v0.0.16
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6
	github.com/muesli/reflow v0.3.0
	github.com/muesli/termenv v0.16.0
	github.com/rivo/uniseg v0.4.7
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3
	golang.org/x/image v0.28.0
	golang.org/x/text v0.26.0
	rsc.io/qr v0.2.0
)

replace (
	github.com/charmbracelet/x/input => ./input
	github.com/sst/opencode-sdk-go => ../sdk/go
)

require (
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/charmbracelet/colorprofile v0.3.1 // indirect
	github.com/charmbracelet/lipgloss v1.1.1-0.20250404203927-76690c660834 // indirect
	github.com/charmbracelet/x/cellbuf v0.0.14-0.20250505150409-97991a1f17d1 // indirect
	github.com/charmbracelet/x/exp/slice v0.0.0-20250327172914-2fdc97757edf // indirect
	github.com/charmbracelet/x/input v0.3.7 // indirect
	github.com/charmbracelet/x/term v0.2.1 // indirect
	github.com/charmbracelet/x/windows v0.2.1 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/yuin/goldmark v1.7.8 // indirect
	github.com/yuin/goldmark-emoji v1.0.5 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.31.0 // indirect
)
