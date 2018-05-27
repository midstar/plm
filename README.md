![GitHub Logo](/images/logo.png)

# plm - Process Load Monitor

[![Go Report Card](https://goreportcard.com/badge/github.com/midstar/plm)](https://goreportcard.com/report/github.com/midstar/plm)
[![AppVeyor](https://ci.appveyor.com/api/projects/status/github/midstar/plm?svg=true)](https://ci.appveyor.com/project/midstar/plm)
[![Coverage Status](https://coveralls.io/repos/github/midstar/plm/badge.svg?branch=master)](https://coveralls.io/github/midstar/plm?branch=master)

Process Load Monitor is a service and application used for monitoring processes. It is primary intended for test purposes, such as monitoring the highest memory consumed by a process and the overal memory allocation trend for processes.

Compared to other tools with similar purposes PLM requires zero configuration, i.e. you don't need to know the PID or make any configuration of the monitor before the measurement. Instead, this tool monitors ALL processes with high resolution (many measured points) in near time and less resolution (fewer measured points) in long time. 

This is a perfect tool to use in your CI Tool (such as Jenkins/Hudson, AppVeyor, CircleCI, Bamboo etc.) for testing that your application does not have any memory leek or is consuming too much memory.

## Features

* Create plots of process memory allocation over time. See [screenshot](images/screenshot_plot.png).
* Get the maximum or minimum memory allocation during a specific time priod or over the life time of the process
* Zero configuration prior measurement, all processes are measured all the time
* No need to know the PID (Process IDentity), only the name of the processes and optionally its command line arguments.
* Add failures if memory exceeds a predefined limit for a process (to be used in the CI tool / your tests)
* Graphical user interface to display all processes (including processes that has died) and also to plot them. See [screenshot](images/screenshot_overview.png).

## Example

Lets say you have an application called myapp.exe that you are integration testing in a Jenkins job. Your requirements is that the application shall not use more than 500 MB (=512 000 KB) memory.

To secure that the requirement is fullfilled add following build step (execute windows batch command) in your Jenkins job, before the actual integration test starts:

    plmc tagset START_TEST

And after your integration test has finished add following build step (execute windows batch command)

    plmc -from START_TEST -m myapp.exe -f 512000 maxmem

Thas all. If myapp.exe exceeds 500 MB the Jenkins job will be marked as fail.

You might want to plot the memory trend over the test. Add following as a post build action:

    plmc -from START_TEST -m myapp.exe plot myapp_plot.html

## Example using java, python or similar

If the application that you are interested in is a "script" such as java or python, the process name will always be java.exe or python.exe. This is not a problem since PLM is able to also match processes based on its command line arguments. For example if your java application is named myapp.jar you simply write (based on previous example):

    plmc -from START_TEST -m "java.exe myapp.jar" -f 512000 maxmem 

## Installation

Get the latest official Windows installer [here on GitHub](https://github.com/midstar/plm/releases).

Experimental / daily builds are found [on AppVeyor](https://ci.appveyor.com/project/midstar/plm/build/artifacts).

## Usage

PLM consist to following components:

* **PLM service (plm deamon)** - This is a background process that measures the processes and stores the measured values in RAM. It also consist of a web server with a user interface (normally at http://localhost:12124/) 
* **PLM Client (plmc)** - This is a command line application for listing processes, checking the maximum memory etc. It is intended to be called from your test tool / CI. In Windows this tool is normally located in C:\Program Files\plm\client\plmc.exe.

Write following to list the available commands in the PLM Client:

    plmc -h

## Configuration

The default configuration should suit most people. See the plm.config file for the available configuration parameters.

## Features to be added in future

* Measure the process CPU usage
* Add support for Linux
* Add support for Mac
* Store measured values on disk instead of RAM (currently all measurements are lost if computer is restarted)

## Author and license

This application is written by Joel Midstjärna and is licensed under the MIT License.