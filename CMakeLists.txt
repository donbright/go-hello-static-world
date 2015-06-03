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
# This file will help to build and statically-link a C language library 
# file to the Go language. It creates a .go "bridge" file automatically 
# from a template 'go.in' file, by modifying the paths in the template 
# file.
#

#
# The way it works is a bit unusual compared to most Cmake situations
# Go has it's own build system built into the language, so cmake is not
# really needed. What is needed, however, is to separate the 'build' 
# files from the 'source' files. 
#
# In this situation, we are building C library files (.a on linux/bsd and
# .lib on win*)... but we are also building a bridge .go file. 
#
# Go also needs files to be layed out in a certain directory structure.
#
# The solution is to build the C library files as normal, generate the 
# .go file as a standard dependency automatically, then 
# copy any .go file that has changed since last build from the source tree
# to the binary tree. 

# C-library files are typically named '.a' files on un*x and '.lib' on Win*


cmake_minimum_required (VERSION 2.8)

project (Hello)

# Flags
#if(NOT WIN32)
#  set(CMAKE_C_FLAGS "-Wall -pedantic -Wimplicit-function-declaration -Wreturn-type -Wstrict-prototypes -Wextra -fPIC")
#  set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Werror")
#  set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -std=c99")
#endif()


####
#### C Library section
####

set(C_LIBRARY_NAME hello)

# Make sure C library files (.a files) are built under 'lib' subdir
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${PROJECT_BINARY_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${PROJECT_BINARY_DIR}/lib)

include_directories(${PROJECT_SOURCE_DIR}/include)

set(CSOURCES src/hello.c)
add_library( ${C_LIBRARY_NAME} ${CSOURCES})




####
#### Go language section
####

# Create bridge.go bridge file using 'go.in' template

set(INFILE ${PROJECT_SOURCE_DIR}/demo/src/bridge/bridge.go.in)
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
add_custom_command(
  OUTPUT demo.go
  COMMAND cp ../demo/demo.go .
  COMMENT "Copying .go files"
)
add_custom_target( gofiles DEPENDS demo.go )
add_dependencies( hello gofiles )

# Call 'go build' to generate executable file
#
# Note that we set the GOPATH environment variable here before running
# go build. This is necessary for go build to work. However it may not
# run correctly depending on the 'shell' program used to execute the
# 'go build' command. It should work on most Linux/BSD command shells. It may
# need to be modified for other platforms. 

add_custom_command(
  TARGET ${C_LIBRARY_NAME}
  POST_BUILD
  COMMAND GOPATH=${PROJECT_BINARY_DIR} go build
  COMMENT "Running go build"
)



