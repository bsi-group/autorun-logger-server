# Compiling

This document assumes that the **golang** tool set has been installed.

Golang applications always have a **src** directory which contains the applications source code, along with any associated projects that are referenced by the primary application. The following shows where the applications source code resides:

    /autorun-logger-server/source/src/info-assure

## gb

The project uses [gb](https://getgb.io) for building the project. **gb** allows for reproducible builds and vendoring so that all dependencies are kept with the project source.

To install **gb**, create a temporary directory and set the GOPATH environment variable to the new temporary directory.
```
$ export GOPATH=/home/bsmith/tempgb
```
Then download the source code for **gb**
```
go get github.com/constabulary/gb/...
```
Navigate to the **gb** sub-directory:
```
cd  /home/bsmith/tempgb/src/github.com/constabulary/gb
```
Build the project
```
go build
```
Copy the binaries to the local path
```
cp ../../../bin/* /usr/local/bin
```
The **gb** command maybe aliased with git, so check with:
```
alias gb
```
If the alias exists then you can unaliase by:
```
unalias gb
```
## Compile with gb

To compile the application use the following commands (assuming the same directory structure):
```
$ cd /autorun-logger-server/source/
$ gb build all
```
