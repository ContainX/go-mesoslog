# Go-Mesoslog

[![Build Status](https://travis-ci.org/gondor/go-mesoslog.svg)](https://travis-ci.org/gondor/go-mesoslog)

## Overview

Go-Mesoslog makes it easy to grab StdOut and StdErr logs from any task.  It allows you to find logs using the name of the task without the UUID portion.  Go-Mesoslog locates the logs via the master and then downloads the desired log from the allocated slave.

## Usage

```
go-mesoslog [global_flags] [action] <action_flags>
Actions
  list                      - List current application id's and task count (instances running)
  print [appId]             - Outputs the log for the given [appId] to StdOut.  Each running instance log will be outputed
  tail [appId]              - Tails the log from the given [appId] to StdOut.  All instances are part of the tail
  file [appId] [targetDir]  - Outputs the log for the given [appId] to a files in the [targetDir] prefixed with the instance TaskID
  help                      - Help about any cmmand

Global Flags
  -m, --master :            - Mesos Master host:port (eg. 192.168.2.1:5050)
```

### Binary Installation 

You can download the binaries (ver 0.1)

 * Architecture i386 [ [linux](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_linux_386.tar.gz?direct) / [windows](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_windows_386.zip?direct) / [darwin](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_darwin_386.zip?direct) / [freebsd](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_freebsd_386.zip?direct) / [openbsd](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_openbsd_386.zip?direct) ]
 * Architecture amd64 [ [linux](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_linux_amd64.tar.gz?direct) / [windows](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_windows_amd64.zip?direct) / [darwin](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_darwin_amd64.zip?direct) / [freebsd](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_freebsd_amd64.zip?direct) / [openbsd](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_openbsd_amd64.zip?direct) ]

Or by installing via deb packages (ver 0.1)

 * [ [amd64](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_amd64.deb?direct) / [armhf](https://dl.bintray.com//content/pacesys/utils/go-mesoslog_0.1_armhf.deb?direct) ]

### Build and Install the Binaries from Source

Add Go-Mesoslog and its package dependencies to your go `src` directory

    go get -v github.com/gondor/go-mesoslog

Once the `get` has completed, you should find your new `go-mesoslog` (or `go-mesoslog.exe`) executable sitting inside the `$GOPATH/bin/`

To update Go-Mesoslog's dependencies, use `go get` with the `-u` option.

    go get -u -v github.com/gondor/go-mesoslog

## Why Go-Mesoslog

Ideally in your larger clusters you should be log forwarding via Syslog, GELF, FluentD, etc but for smaller deployments where this isn't setup such as DEV or TEST environments this offers a quick view without exposing all of your cluster characteristics to the teams.

I wrote this quickly in one night to solve a few use cases.  Feel free to fork and contribute!

## Looking for Mesos/Marathon, Kubernetes or ECS deployment strategies?

Checkout my other project [go-mesoslog](https://github.com/gondor/go-mesoslog)

## License

This software is licensed under the Apache 2 license, quoted below.

Copyright 2015 Jeremy Unruh

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations under
the License.