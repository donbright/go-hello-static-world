# Copyright 2012 Douglas Linder
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Modified by Don Bright 2015


# This cmake file is intended to help with the situation where you have 
# a Go language program that needs to call C static libraries. Please 
# see the README.md file for a fuller explanation of how this process 
# works.

cmake_minimum_required (VERSION 2.8)

project (hello)

####
#### C Library section
####

set(C_INCLUDE_PATH ${PROJECT_SOURCE_DIR}/src/c)
set(C_LIBRARY_NAME world)
set(CSOURCES src/c/world.c)

# Make sure C library files (.a/.lib files) are built under 'lib' subdir
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${PROJECT_BINARY_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${PROJECT_BINARY_DIR}/lib)
set(CMAKE_BINARY_OUTPUT_DIRECTORY ${PROJECT_BINARY_DIR}/bin)

include_directories(${C_INCLUDE_PATH})

add_library( ${C_LIBRARY_NAME} ${CSOURCES} )




####
#### Go language section
####

# Create bridge.go bridge file using 'go.in' template

set(INFILE ${PROJECT_SOURCE_DIR}/src/bridge/bridge.go.in)
set(OUTFILE ${PROJECT_BINARY_DIR}/src/bridge/bridge.go)
set(C_LIBRARY_PATH ${PROJECT_BINARY_DIR}/lib)

set(CFLAGS "-I${C_INCLUDE_PATH}")
set(LDFLAGS "-L${C_LIBRARY_PATH} -l${C_LIBRARY_NAME}")

message(STATUS "Creating .go bridge file")
message(STATUS "Input template:${INFILE}" )
message(STATUS "   Output file:${OUTFILE}")
message(STATUS "replacing CFLAGS with ${CFLAGS}")
message(STATUS "replacing LDFLAGS with ${LDFLAGS}")

configure_file( ${INFILE} ${OUTFILE} )

# Populate the build tree with .go files
##.. foreach... (not implemented)
## copy if changed... (not implemented)
#add_custom_command(
#  OUTPUT demo.go
#  COMMAND cp ../demo/demo.go .
#  COMMENT "Copying .go files"
#)
#add_custom_target( gofiles DEPENDS demo.go )
#add_dependencies( hello gofiles )

# Call 'go build' to generate executable file
#
# Note that we set the GOPATH environment variable here before running
# go build. This is necessary for go build to work. However it may not
# run correctly depending on the 'shell' program used to execute the
# 'go build' command. It should work on most Linux/BSD command shells. It may
# need to be modified for other platforms. 

# GOPATH see https://golang.org/doc/articles/go_command.html
# note: order matters... binary must come first so bridge.go can be found.
set(GOPATH ${PROJECT_BINARY_DIR}:${PROJECT_SOURCE_DIR})
set(GOBIN "${PROJECT_BINARY_DIR}/bin")

message("GOPATH=${GOPATH}")
message("GOBIN=${GOBIN}")
message(STATUS "executable file: ${GOBIN}/${CMAKE_PROJECT_NAME}")

# create shell programs to help user set GOPATH and GOBIN

set(SHFILE ${PROJECT_BINARY_DIR}/setenv.sh)
file(WRITE  ${SHFILE} "")
file(APPEND ${SHFILE} "GOPATH=${GOPATH}\n")
file(APPEND ${SHFILE} "GOBIN=${GOBIN}\n")
file(APPEND ${SHFILE} "export GOPATH\n")
file(APPEND ${SHFILE} "export GOBIN\n")
file(APPEND ${SHFILE} "echo GOPATH=$GOPATH\n")
file(APPEND ${SHFILE} "echo GOBIN=$GOBIN\n")

set(CSHFILE ${PROJECT_BINARY_DIR}/setenv.csh)
file(WRITE  ${CSHFILE} "")
file(APPEND ${CSHFILE} "setenv GOPATH ${GOPATH}\n")
file(APPEND ${CSHFILE} "setenv GOBIN ${GOBIN}\n")
file(APPEND ${CSHFILE} "echo GOPATH=$GOPATH\n")
file(APPEND ${CSHFILE} "echo GOBIN=$GOBIN\n")


# setup 'go install' for Makefile

# set up static linking of Go code itself. Use a few cmake tricks with quotes
set(GOFLAGS "--ldflags" "'-extldflags" "\"-static\"'")

# use the -x option when debugging the build process. it will dump more info 
#set(GOFLAGS "-x" "--ldflags" "'-extldflags" "\"-static\"'")

#set(CUSTCMD strace -f go install -v hello)
#set(CUSTCMD go build hello)

set(GOCMD go install ${GOFLAGS} ${CMAKE_PROJECT_NAME})
add_custom_command(
  TARGET ${C_LIBRARY_NAME}
  POST_BUILD
  COMMAND GOBIN=${GOBIN} GOPATH=${GOPATH} ${GOCMD}
  COMMENT "Running ${GOCMD}"
)



