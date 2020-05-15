# snmpsim-cli-manager

[![Go Report Card](https://goreportcard.com/badge/github.com/inexio/snmpsim-cli-manager)](https://goreportcard.com/report/github.com/inexio/snmpsim-cli-manager)
[![GitHub license](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/inexio/snmpsim-cli-manager/LICENSE)
[![GitHub code style](https://img.shields.io/badge/code%20style-uber--go-brightgreen)](https://github.com/uber-go/guide/blob/master/style.md)

## Description

A tool for using the [snmpsim](https://github.com/etingof/snmpsim)  [REST API](https://github.com/etingof/snmpsim-control-plane) via the command-line written in go using the cobra CLI framework.

The tool uses our [snmpsim-restapi-go-client](https://github.com/inexio/snmpsim-restapi-go-client).

## Code Style

This project was written according to the **[uber-go](https://github.com/uber-go/guide/blob/master/style.md)** coding style.

## Features

### Full SNMPSIM Feature Support

This client allows you to create/delete:

- Labs
- Agents
- Engines
- Endpoints
- Users
- Tags
- All objects tagged with a certain tag

and also enables you to:

- Get data for all available components
- Activate/Deactivate Labs
- Tag/Untag all available objects
- Upload/Delete record files
- Read various metrics

## Requirements

Requires a running instance of [snmpsim](https://github.com/etingof/snmpsim) and its
[REST API server](https://github.com/etingof/snmpsim-control-plane).

To check if your setup works, follow the steps provided in the **'Tests'** section of this document.

## Installation

```
go get github.com/inexio/snmpsim-cli-manager
```

or

```
git clone https://github.com/inexio/snmpsim-cli-manager.git
```

## Setup

After downloading you have to decide wether you use the tool via source or if you want to use its binary.

For using the tool via source you can run:

```
go run snmpsim-cli-manager/snmpsim/main.go  
```

or you can compile it to a binary:

```
cd snmpsim-cli-manager/snmpsim
go install
cd $GOBIN
snmpsim  
```

After installing the tool you have to either declare a path to your config file or set certain environment variables for the tool to work.

These can be set as follows:

#### Config File

##### Using source

In **cmd/root.go**, in the **initConfig()** function you can see the following lines of code:

```
viper.AddConfigPath("config/")
viper.SetConfigType("yaml")
viper.SetConfigName("snmpsim-cli-manager-config")
```

ConfigPath is relative to the package location.
ConfigType and name can also be changed to match your desired type of config.

##### Using the binary

If you are using the binary you have to add a folder named **config** in the same directory the binary is located at.

In this **config** folder you have to add a file named **snmpsim-cli-manager-config.yaml**. This file has to contain all data required in the snmpsim-cli config.

##### Using the '--config' flag

If you want to use an alternative config file for only one command, you can use the --config flag with the command. Here's an example of this:

```
snmpsim get labs --config ~/go/src/snmpsim-cli-manager/snmpsim/config/configtwo.yaml
```

Note, that this only works for one command and doesn't permanently change the config file. 

##### Using the env-var SNMPSIM_CLI_CONFIG

You can also set an environment variable to read in the config file given in that variable.

To set this variable use:

```
export <YOUR_ENV_PREFIX>_CONFIG=/absolut/path/to/your/config
```

This way the config will always be read in instead of how '--config' has to be set every time you want to execute a command.


#### Env-Config

If you want to use the **setup-env** feature with the **--env-config** flag you'll need a valid config file located on your machine. 

An example of such a config can be found [here](https://github.com/inexio/snmpsim-cli-manager/blob/master/snmpsim/env-config/env-config.yaml)

#### Environment Variables

Also in the **root.go** file, in the **initConfig()** function you will find:

```
//Set env var prefix to only match certain vars
viper.SetEnvPrefix("SNMPSIM_CLI")
```

`SetEnvPrefix` can be changed to whatever prefix you prefer to have in your environment vars.

The needed environment variables can then be added as follows:

For the management endpoint:

```
export SNMPSIM_CLI_MGMT_HTTP_BASEURL="<your mgmt baseUrl>"
export SNMPSIM_CLI_MGMT_HTTP_AUTHUSERNAME="<your username>"
export SNMPSIM_CLI_MGMT_HTTP_AUTHPASSWORD="<your password>"
```

and for the metrics endpoint:

```
export SNMPSIM_CLI_METRICS_HTTP_BASEURL="<your metrics baseUrl>"
export SNMPSIM_CLI_METRICS_HTTP_AUTHUSERNAME="<your username>"
export SNMPSIM_CLI_METRICS_HTTP_AUTHPASSWORD="<your password>"
```

## Usage

When doing ```snmpsim help```, you can get information about all commands:
![snmpsim help](https://imgur.com/HdmuDmU.gif)

The following section will show you how to create a lab and do various operations with it.

```go
snmpsim create lab --name TestLab
```

Should return:

```
Lab has been created successfully.
Id: 1
```

You can get information about a component via the get command:

```
snmpsim get lab 1
```

Which should return:

```
Lab
  Id: 1
  Name: TestLab
  Power: off

  Agents(0)
    /

  Tags(0)
    /
```

To power the lab on use:

```
snmpsim power 1 --on
```

Which should return:

```
Power of Lab 1 has been successfully set to on
```

And after a second get operation the labs power value should be 'on':

```
Lab
  Id: 1
  Name: TestLab
  Power: on

  Agents(0)
    /

  Tags(0)
    /
```
![snmpsim lab](https://imgur.com/cWXrpus.gif)

To create an agent use the following command:

```
snmpsim create agent --name TestAgent --dataDir "/opt/snmpsim/data/agent1-test" 
```

which returns:

```
Successfully created agent.
Id: 1
```

To add a sub-component to its corresponding main-component use the add command:

```
snmpsim add agent-to-lab --agent 1 --lab 1

Agent 1 has been added to lab 1
```

Alternativly, you can add a sub-component to its main-component while creating it:

```
snmpsim create agent --name TestAgent --dataDir "/opt/snmpsim/data/agent1-test" --lab 1
```

And finally to remove a sub-component from its main-component use the remove command:

```
remove agent-from-lab --agent 1 --lab 1

Agent 1 has been removed from lab 1
```
![snmpsim agent](https://imgur.com/23xSzgs.gif)
## Tests

In order to test if your setup is operational you can use the following command:

```
go run main.go get labs
```

or

```
snmpsim get labs
```

The output should look something like this:

```
Labs(0)
  /
```

or if you have any labs configured it should look like this:

```
Labs(2) 

  Lab
    Id: 1
    Name: TestLab1
    Power: off

    Agents(1)


    Tags(1)


  Lab
    Id: 2
    Name: TestLab2
    Power: off

    Agents(1)


    Tags(1)
```

## Getting Help

If there are any problems or something does not work as intended, open an [issue](https://github.com/inexio/snmpsim-cli-manager/issues/new/choose) on GitHub.

If you have problems with the usage of the application itself you can use the built in **--help** flag to get some useful information about every command.

## Contribution

Contributions to the project are welcome.

We are looking forward to your bug reports, suggestions and fixes.

If you want to make any contributions make sure your go report does match up with our projects score of **A+**.

When you contribute make sure your code is conform to the **uber-go** coding style.

Happy Coding!
