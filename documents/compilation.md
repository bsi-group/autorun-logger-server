# Compiling

This document assumes that the **golang** tool set has been installed.

Golang applications always have a **src** directory which contains the applications source code, along with any associated projects that are referenced by the primary application. The following shows where the applications source code resides:

    /autorun-logger-server/source/src/info-assure

To compile golang applications, the **go build** command requires the **GOPATH** environment variable to be set. Assuming the source code is using the directory structure detailed above, then set the environment variable like so:

   $ export GOPATH=/autorun-logger-server/source

The environment variable is always set one directory level above the **src** folder.

To compile the application use the following commands (assuming the same directory structure):

```
$ cd /autorun-logger-server/source/src/info-assure
$ go build -o arl
```
