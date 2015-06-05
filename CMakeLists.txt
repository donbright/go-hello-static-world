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

#
# This cmake file is intended to help with the situation where you have
# a Go language program that needs to call C static libraries. Since
# Go and C have very different build systems, it is a bit tricky to set up
# a project where you will be working on both C and Go files simultaneously.
# This attempts to make that process easier by allowing you to type the
# familiar 'cmake && make' for the entire process... so that it will 
# build the C libraries and then start the traditional 'go build' process
# after those libraries have been built, all in a separate out-of-source
# binary tree. 

# It creates a .go "bridge" file automatically from a template 'go.in' 
# file, by modifying the paths in the template file.
#
# It also creates all binary code in a separate subdirectory, as is traditional
# with C programs, but also manages to do this with the Go program.

# C-library files are typically named '.a' files on un*x and '.lib' on Win*


cmake_minimum_required (VERSION 2.8)

project (hello)

####
#### C Library section
####

set(C_LIBRARY_NAME hello)

# Make sure C library files (.a/.lib files) are built under 'lib' subdir
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${PROJECT_BINARY_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${PROJECT_BINARY_DIR}/lib)
set(CMAKE_BINARY_OUTPUT_DIRECTORY ${PROJECT_BINARY_DIR}/bin)

include_directories(${PROJECT_SOURCE_DIR}/include)

set(CSOURCES src/hello.c)
add_library( ${C_LIBRARY_NAME} ${CSOURCES})




####
#### Go language section
####

# Create bridge.go bridge file using 'go.in' template

set(INFILE ${PROJECT_SOURCE_DIR}/src/bridge/bridge.go.in)
set(OUTFILE ${PROJECT_BINARY_DIR}/src/bridge/bridge.go)
set(C_INCLUDE_PATH ${PROJECT_SOURCE_DIR}/include)
set(C_LIBRARY_PATH ${PROJECT_BINARY_DIR}/lib)

message(STATUS "Creating .go bridge file")
message(STATUS "Input template:${INFILE}" )
message(STATUS "   Output file:${OUTFILE}")
message(STATUS "replacing C_INCLUDE_PATH with ${C_INCLUDE_PATH}")
message(STATUS "replacing C_LIBRARY_PATH with ${C_LIBRARY_PATH}")
message(STATUS "replacing C_LIBRARY_NAME with ${C_LIBRARY_NAME}")

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

#set(CUSTCMD strace -f go install -v hello)
#set(CUSTCMD go build hello)

set(GOCMD go install -v ${CMAKE_PROJECT_NAME})
add_custom_command(
  TARGET ${C_LIBRARY_NAME}
  POST_BUILD
  COMMAND GOBIN=${GOBIN} GOPATH=${GOPATH} ${GOCMD}
  COMMENT "Running ${GOCMD}"
)



