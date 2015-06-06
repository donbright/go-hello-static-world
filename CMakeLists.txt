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

# This Cmake file intends to help linking Go code statically with C code
# Please see the README.md file for a fuller explanation of how this process 
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
message(STATUS "Output file   :${OUTFILE}")
message(STATUS "new CFLAGS : ${CFLAGS}")
message(STATUS "new LDFLAGS: ${LDFLAGS}")

configure_file( ${INFILE} ${OUTFILE} )


# Call 'go build' to generate executable file
# GOPATH / GOBIN
# see https://golang.org/doc/articles/go_command.html
# note: order matters... binary must come first so bridge.go can be found.
set(GOPATH ${PROJECT_BINARY_DIR}:${PROJECT_SOURCE_DIR})
set(GOBIN "${PROJECT_BINARY_DIR}/bin")

message(STATUS "GOPATH=${GOPATH}")
message(STATUS "GOBIN=${GOBIN}")
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

# To dump alot more debug info with go build '-x', set the
# environment variable VERBOSE=1. (for example, VERBOSE=1 cmake)

if(DEFINED ENV{VERBOSE})
  set(GOFLAGS "-v" "-x" "--ldflags" "'-extldflags" "\"-static\"'")
else()
  set(GOFLAGS "--ldflags" "'-extldflags" "\"-static\"'")
endif()

#set(CUSTCMD strace -f go install -v hello)
#set(CUSTCMD go build hello)

set(GOCMD go install ${GOFLAGS} ${CMAKE_PROJECT_NAME})
add_custom_command(
  TARGET ${C_LIBRARY_NAME}
  POST_BUILD
  COMMAND GOBIN=${GOBIN} GOPATH=${GOPATH} ${GOCMD}
  COMMENT "Running ${GOCMD}"
)

# full static linking requires typing a lot of flags. this will
# create a 'go.sh' file that will include all that detail

string(REPLACE ";" " " GOCMD "${GOCMD}")
file(WRITE go.sh "${GOCMD}\n") 
file(COPY go.sh
  DESTINATION ${PROJECT_BINARY_DIR}
  FILE_PERMISSIONS OWNER_READ OWNER_WRITE OWNER_EXECUTE)
file(REMOVE tmp)
