groku
=====

CLI utility for controlling your [Roku](https://www.roku.com/)

###Installation
You may download a precompiled binary

#####Windows
--
- [32bit](https://github.com/zankich/groku/releases/download/0.3/groku_0.3_windows_386.zip)
- [64bit](https://github.com/zankich/groku/releases/download/0.3/groku_0.3_windows_amd64.zip)  

#####Linux
--
- [32bit](https://github.com/zankich/groku/releases/download/0.3/groku_0.3_linux_386.tar.gz)
- [64bit](https://github.com/zankich/groku/releases/download/0.3/groku_0.3_linux_amd64.tar.gz)
- [armhf](https://github.com/zankich/groku/releases/download/0.3/groku_0.3_linux_arm.tar.gz)

#####OSX
--
- [32bit](https://github.com/zankich/groku/releases/download/0.3/groku_0.3_darwin_386.zip)
- [64bit](https://github.com/zankich/groku/releases/download/0.3/groku_0.3_darwin_amd64.zip)

Or install from source  
`go get github.com/zankich/groku`

###Usage
```
$ groku
NAME:
   groku - roku CLI remote

USAGE:
   groku [global options] command [command options] [arguments...]

VERSION:
   0.3

COMMANDS:
   home         home
   rev          rev
   fwd          fwd
   select       select
   left         left
   right        right
   down         down
   up           up
   back         back
   info         info
   backspace    backspace
   enter        enter
   search       search
   replay       replay
   play         play/pause
   discover     discover roku on your local network
   text         send text to the roku
   apps         list installed apps on roku
   app          launch specified app
   help, h      Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --version, -v        print the version
   --help, -h           show help

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
