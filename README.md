# Zebra
<img src="./logo/zebra.svg" alt="zebra-logo" width="300px"/>
Zebra is a CLI tool that allows to read, parse and analyse log files to retrieve metrics.

## Behavior
The behavior of the tool is the following  
Zebra will read one or several log files, given by user. The tool will extract only supported format logs and aggregate them.  
The aggregations are the following:  
Logs are aggregated by levels, showing number of lines for each level  
Logs are aggregated by service, showing for each service:  
- Number of lines  
- Average duration of processes  
Zebra also shows number of errors encountered during parsing (non existing files, format errors, etc) but does not natively give details on those errors.  
If the user is requesting slowest logs using --top flag, the tool will also give full log lines for slowest logs, depending on the number requested  

## Logs Format
Currently, only one format of logs is supported :

### DATE LEVEL service=SERVICE message=MESSAGE duration=DURATION [PROPS]

- **DATE** is date and time, RFC3339 format (YYYY\-MM\-DDTHH\:MM\:SSZ)
- **LEVEL** is a string in [DEBUG, INFO, WARNING, ERROR, FATAL]
- **SERVICE** is a string describing service logged
- **MESSAGE** is a message
- **DURATION** is the duration of the process, in milliseconds, in the format numberDURATION where DURATION can be h(our)/m(inute)/s(second)/ms(millisecond)
- **PROPS** is one or multiple key-value pairs, in the format key=value, separated by spaces. This is optionnal

### JSON format

JSON format is supported and the following rules are to follow:

- The following fields are mandatory:
    - **date** field needs to be followed by a RFC3339 formatted date string
    - **service** field is followed by a service compatible string
    - **message** field is followed by the string message
    - **duration** field is the duration of the process
    - any **other** field is going to be an extra field
    - format and rules applied to default format are the same for json (date, duration, etc)


## Flags
When using Zebra, it is possible to add flags to modify behavior. Some flags are necessary, some are optionals.
### Necessary flags
- **files**: log files to analyse
### Optional flags
- **startDate**: log date to start from
- **endDate**: log date to end to
- **service**: filter logs by service
- **level**: filter logs by level
- **top**: number of slowest logs to show
