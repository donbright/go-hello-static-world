# Go-C static linked demo

Use Go language with statically linked C libraries, with CMake build system.

This is a fork of shadowmint's go-static-linking project

## Usage

    mkdir build
    cd build
    cmake ..
    make
    ./demo

## Basic design of Cmake file

   Step 1 - create static C library (libHello.a)
   Step 2 - create .go 'bridge' from .go.in template
   Step 3 - build executable binary by calling 'go build'

## Windows(TM)

Untested in this fork.

## Differences from shadowmint's original code

1. remove usage of separate cmake BindConfig.txt file
2. put all generated files under 'build' directory (bridge.go)
3. move call of 'go build' into cmake process

