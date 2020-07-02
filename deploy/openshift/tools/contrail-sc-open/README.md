# Open Contrail neccessary ports on AWS

This tool alows to automatically open ports neccessary for Contrail on AWS in every Security Group attached do cluster resources.

## Build

In order to build this tool use `go build .` command.
Afterwards, you should have binary *contrail-sc-open* which is compiled tool.

## Usage

In order to use it run:
```
./contrail-sc-open -cluster-name <name of your Openshift cluster> -region <AWS region wherhe cluster is located>
```

Tool will log all security groups found and status whether it successfuly added new rules for Contrail ports.
