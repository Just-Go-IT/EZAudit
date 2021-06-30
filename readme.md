[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
![](https://github.com/Just-Go-IT/EZAudit/actions/workflows/RunTests.yml/badge.svg)

## _EZAudit_

EZAudit is an audit framework that checks system security settings. The first step is to check the execution environment - the running Operating System and if Admin Rights are available. A previously written configuration file will be parsed and the containing audits will be executed. One configuration contains multiple audits which can have multiple steps. Steps might have conditions for a second step to be executed. Each step saves the output of the command that has been executed via Bash, Terminal, PowerShell, or whatever is used by the OS. 
There are 2 given configuration files based on CIS Benchmarks published from the [Center of Information Security].
After auditing, you receive the output by default in a zip file containing a debug.log, ResultReport.json and an artifact folder with all the collected artifacts from the auditsteps.

The whole project is built in a modular way in order to be extended without deeper knowledge of the program. You are able to add new modules or support other Operating Systems within minutes. You can find more information about adding new modules in the [User Manual] (GER).

This project emerged from a semester at the Hochschule Mannheim during summer of 2021. We are a team of 6 people called Just Go IT.   

## Quick Start
There are 2 ways of getting started:

*Case 1 - No new modules needed:*\
Download the latest release from [Downloads](#downloads). After extracting the folder you will have the executable to start the program and you are ready to go. 

*Case 2 - Adding new modules:*\
If you want to add new modules you will need an IDE (e.g. GoLand) to edit the code. As mentioned before, a guide for adding modules is written in the user manual.
Then you will have to build a new executable file.

Depending on your target hardware and operating system you have to set the go environment variables:
```sh
set GOOS=your OS
set GOARCH=your CPU architecture
```

with the following command you receive the list of possible variables:
```sh
go tool dist list
```

then navigate to the EZAudit folder within your terminal and type the following command to build the executable:
```sh
go build
```

After building, you can start the executable by double-clicking. The program will be started with predetermined default values. By default, the program is looking for a "config.json" within the executing folder and will start with the verbosity level 4 (Information). 
If you want to modify the starting parameters you can make use of the flags provided:

|Short Flag|Long Flag|Description
|:-----:|:-:|-----
| <code>-h</code> | <code>--help</code> | show the list of possible flags 
| <code>-d</code> | <code>--dryRun</code> | the DryRun validates the config file. Checks JSON syntax, <br> if the OS is matching and if Modules have valid parameters
| <code>-f</code> | <code>--force</code> | forces the program to execute the config file ignoring the operating system
| <code>-c=</code> | <code>--config=</code> | pass the relative path of the config file. For example <br> <code>--config=/Desktop/ExampleConfig.json</code>
| <code>-v=</code> | <code>--verbosity=</code> | sets the verbosity level for logging (1-5)
| <code>-nz</code> | <code>--noZip</code> | don't zip the result folder

## Downloads

Here is the link to the latest [GitHub] release.

## GoDoc

Here is the link to the [GoDoc] documentation.

## Manual

Here is the link to the [User Manual] (GER).

## Technical Documentation

Here is the link to the [Technical Documentation] (GER).

## Build Status

| OS | Status
| --- | ---
| Windows | ![](https://github.com/Just-Go-IT/EZAudit/actions/workflows/windows.yml/badge.svg)
| Ubuntu | ![](https://github.com/Just-Go-IT/EZAudit/actions/workflows/ubuntu.yml/badge.svg)

## License
[MIT License]

[Center of Information Security]:<https://downloads.cisecurity.org/#/>
[GitHub]:<https://github.com/Just-Go-IT/EZAudit/releases>
[MIT License]:<https://github.com/Just-Go-IT/EZAudit/blob/main/LICENSE.txt>
[GoDoc]:<https://just-go-it.github.io/EZAudit/>
[User Manual]:<https://github.com/Just-Go-IT/EZAudit/blob/main/Manual_GER.pdf>
[Technical Documentation]:<https://github.com/Just-Go-IT/EZAudit/blob/main/Technical_Documentation_GER.pdf>
