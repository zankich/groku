groku
=====

CLI remote for controlling your [Roku](https://www.roku.com/)

### Installation
You may download a precompiled binary from the latest [release](https://github.com/zankich/groku/releases)

Or install from source  
`go get github.com/zankich/groku`

### Usage
```
$ groku
usage: groku [--version] [--help] <command> [<args>]

CLI remote for your Roku

Commands:
  home            Return to the home screen
  rev             Reverse
  fwd             Fast Forward
  select          Select
  left            Left
  right           Right
  up              Up
  down            Down
  back            Back
  info            Info
  backspace       Backspace
  enter           Enter
  search          Search
  replay          Replay
  play            Play
  pause           Pause
  discover        Discover a roku on your local network
  text            Send text to the Roku
  apps            List installed apps on your Roku
  app             Launch specified app
```
Play or pause current video:  
`$ groku play`

Rewind video:  
`$ groku rev`

Fast Forward video:  
`$ groku fwd`

Enter text into search box:  
`$ groku text "Breaking Bad"`

Query installed apps:  
`$ groku apps`

Launch installed app:  
`$ groku app "Netflix"`
