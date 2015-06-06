# Hello World using Go and Statically-linked C library

This program prints "Hello World". A go language program prints the word 
'Hello' and then it calls a statically linked C library to print 
'World'.

This is a fork of shadowmint's go-hello-static-world project

## Usage

To build and run, copy/paste these commands into a terminal:
   
    mkdir build && cd build && cmake ..
    VERBOSE=1 make
    ./bin/hello

## Status

This is an hello world program, a basic learning tool. It is not intended
for any professional use. See the license. 

## Design

The source tree is layed out as follows:

    src/hello/hello.go       # main go package
    src/bridge/bridge.go.in  # template, will generate 'glue' .go package
    src/c/world.c            # c code all lives under 'c' path
    src/c/world.h            
    CMakeLists.txt           # Cmake build file

As you can see, this is a sort of 'hybrid' go project layout. It 
has subdirs under 'src' for go packages. But it also has c code and
a Cmake build file.

During the cmake && make process, Cmake & go's builder will go through 
the following steps:

    Create bridge .go file from the bridge.go.in template
    Create c library (.a/.lib) file from c source code
    Call 'go install' to create main executable file, statically linked

The resulting binary tree will hopefully be as follows:

    build/bin/hello          # executable file generated by Go's builder
    build/lib/libworld.a     # statically linkable C library
    build/src/bridge.go      # "bridge" generated by cmake process
    build/cmake*             # usual Cmake generated files (cache, etc)
    build/Makefile           # cmake-generated, builds C lib & runs 'go install'
    build/pkg/machine/bridge.a # bridge library generated by 'go install'

Running ./bin/hello should produce output like this:

    Hello (Invoking c statically linked library...)
    World
    (Done)

## GOBIN, GOPATH, What do to do after 'make'

To rebuild the C code, just type 'make' from the command shell.

To rebuild the 'go' code, you need to first set up the GOBIN and GOPATH
environment variables. These are printed during the 'cmake' run. Cmake also
generates shell scripts to set them automatically. Run one depending on your shell:

    . setenv.sh          # for bash shell (typical shell on linux)
    source setenv.csh    # for csh shell (typical shell on BSD)
    go build hello       # now you can run go commands

After this you can type 'go build hello' or 'go install hello' and it should
build and link the 'bin/hello' file as needed. 

However if the bridge code and/or the C code changes, you will need to 
possibly rerun make and/or cmake to regenerate the C library file (.lib or .a)
and also regenerate the bridge.go glue code.
 
## On the static glue of bridge.go.in and bridge.go

The magic static glue of the bridge file works using 'cgo' as follows:

bridge.go.in has two lines like this:

    // #cgo CFLAGS: ${CFLAGS}
    // #cgo LDFLAGS: ${LDFLAGS}

bridge.go, generated by CmakeLists.txt, has two lines like this:

    // #cgo CFLAGS: -I/tmp/go-hello-static-world/src/c
    // #cgo LDFLAGS: -L/tmp/go-hello-static-world/build/lib -lworld

## Windows(TM)

Untested in this fork.

## Differences from shadowmint's original code

    1. remove usage of separate cmake BindConfig.txt file
    2. put all generated files under 'build' directory (bridge.go)
    3. move call of 'go build' into cmake process (go install)
    4. rearrange and simplify directory structure
    5. link everything statically in the binary, not just the C code


## See Also

<http://golang.org/doc/code.html> How to Write Go Code (golang.org)
<http://blog.golang.org/c-go-cgo> C-Go (especially for the bridge.go file)
<http://golang.org/misc/cgo/testso/cgoso.go> setting LDFLAGS on different OSes
<http://blog.hashbangbash.com/2014/04/linking-golang-statically/> Static Go-linking

## Why do this?

When transforming a program from C to Go, where the program is
50,000 lines of C, it may be easier to transform it in small pieces
and of course some C libraries may not be transformable ever. 

Static linking is a nice way to simplify the user experience. 
No dynamic libraries to worry about.

## How do you know if a file is really statically linked?

ldd and file can help you

    don@serebryanya[build]$ ldd ./bin/hello 
	not a dynamic executable
    don@serebryanya[build]$ file bin/hello 
        bin/hello: ELF 64-bit LSB  executable, x86-64, version 1 (SYSV), 
        statically linked, for GNU/Linux 2.6.24, 

Thanks for reading. Thanks shadowmint.


