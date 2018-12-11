# audigo
3D Led CubeのPCレス(Raspberry pi)音響サービス(SDL版)  

<!-- toc -->  
* [💊 Requirements](#-requirements)
* [📌 Installing](#-installing)
* [🎧 Usage](#-usage)
* [🌏 REST Api](#️-rest-api)
* [🎃 Notes](#-notes)
* [🎤 Third party](#-third-party)
<!-- tocstop -->  

# Getting Started
## 💊 Requirements

**ALL**  
* git  
* dep 
    ```sh
    $ go get -u github.com/golang/dep/cmd/dep
    ```
* Go 1.11 or later  

see https://github.com/veandco/go-sdl2  
  
**Linux**  
```sh
$ sudo apt install libsdl2-dev libsdl2-mixer-dev
```

**Windows**  
Setup SDL2 for mingw  
  
  
  
## 📌 Installing

1. Goto GOPATH  
    **WIndows**
    ```sh
    $ cd %GOPATH%
    ```

    **Others**
    ```sh
    $ cd $HOME/go
     or
    $ cd $GOPATH
    ```

2. Get src
    ```sh
    $ git clone https://github.com/code560/audigo-sdl.git ./src/github.com/code560/audigo-sdl
    $ cd ./src/github.com/code560/audigo-sdl
    $ dep ensure
    ```

3. Build
    ```sh
    $ go build -a
    ```

4. Create log directory
    ```sh
    $ mkdir log
    ```

5. Set startup  
    Create service file.    
    ```sh
    Copy .service file  
    $ sudo cp install/audigo.service /etc/systemd/system/.

    Copy startup shell file  
    $ sudo mkdir -p /opt/audigo/bin
    $ sudo chmod 755 /opt/audigo/bin
    $ sudo cp install/audigo.sh /opt/audigo/bin/.
    $ sudo chown root:root /opt/audigo/bin/audigo.sh
    $ sudo chmod 755 /opt/audigo/bin/audigo.sh
    ```
  
    Enable startup service  
    ```sh
    $ sudo systemctl enable audigo

    or Update .service
    $ sudo systemctl daemon_update

    Self start service
    $ sudo systemctl start audigo

    Self stop service
    $ sudo systemctl stop audigo

    Check service state
    $ sudo systemctl status audigo
    $ sudo systemctl list-dependencies
    ```


# 🎧 Usage
Start audio service  
```sh
$ ./audigo-sdl
```

## 🔨 Options

### port
add port number (default port 80)  

```sh
Listening port 5701
$ ./audigo-adl 5701

Listening port 8080
$ ./audigo-adl 8080

Listening port 80
$ ./audigo-adl
```

## 📖 help

```sh  
$ ./audigo-sdl -h

NAME:
   audigo-sdl - Audio service by LED CUBU

USAGE:
   audigo-sdl.exe [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     server, s  Instant server mode. (default)
     client, c  Instant client REPL mode.
     repl, r    Instant local REPL mode.
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cd value, -c value  change current directory by repl
   --help, -h            show help
   --version, -v         print the version
```

Client mode help
```sh
$ ./audigo-sdl c -h

NAME:
   audigo-sdl.exe client - Instant client REPL mode.

USAGE:
   audigo-sdl.exe client [command options] [arguments...]

OPTIONS:
   --cd value, -c value      change current directory by repl
   --domain value, -d value  set request domain url by client (default: "http://audigo.local")
```  


## 📂 Directory layout

Add sound file in audigo-adl/asset/audio

```sh
audigo
 |-- audigo
 |-- asset
      |-- audio
           |-- bgm_wave.wav
           |-- hogehoge.wav
           |-- foobar.wav
           |-- ...

```
  
Support sound file format  
* wave
* sample rate: 44100
  
  
# 🌏️ REST Api
| REST | URI                             | note                          | arguments (json)     |
|------|---------------------------------|-------------------------------|----------------------|
| GET  | /audio/v1/ping                  | I Can Fly !                   | none                 |
| POST | /audio/v1/init/\<content id>    | init players in memory        | none                 |
| POST | /audio/v1/play/\<content id>    | play sound                    | src: "bgm_wave.wav"<br> (file name in ./asset/audio/) <br><br>loop: true or false<br> (loop play or single play) <br>        |
| POST | /audio/v1/stop/\<content id>    | stop content player sound     | none                 |
| POST | /audio/v1/volume/\<content id>  | change volume                 | vol: 0.5 (0 - 1, 0 is silent)          |


# 🎃 Notes

| Platform / Architecture        | x86 | x86_64 |
|--------------------------------|-----|--------|
| Windows (7, 10 or Later)       | -   | ✓     |
| Rasbian (STRETCH or Later)     | ✓  | -      |
| OSX (10.14 or Later)           | -   | -     |


# 🎤 THIRD PARTY

Use libs
* [gin-gonic/gin](https://github.com/gin-gonic/gin)
* [urfave/cli](https://github.com/urfave/cli)
* [go.uber.org/zap](https://github.com/uber-go/zap)
* [veandco/go-sdl2](https://github.com/veandco/go-sdl2)


音声ファイルを使用させていただいております。
* [効果音ラボ](https://soundeffect-lab.info)  
* [あみたろの声素材工房](http://www14.big.or.jp/~amiami/happy/)



以上  