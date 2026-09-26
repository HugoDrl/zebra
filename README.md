# Zebra
<img src="./logo/zebra.svg" alt="zebra-logo" width="300px"/>
Zebra is an HTTP tool that allows to read, parse and analyse log files to retrieve metrics.

## Setup
Once you've cloned this repo, to build the project you can use build.sh,
and start it using dist/zebra executable

```bash
bash ./build.sh
./dist/zebra
```

## Behavior
The behavior of the tool is the following  
Zebra will read one or several log files, given by requests. The tool will extract only supported format logs and aggregate them.  
Zebra also shows number of errors encountered during parsing (non existing files, format errors, etc) but does not natively give details on those errors.  

## Logs Values

- **date** field is date and time, RFC3339 format (YYYY\-MM\-DDTHH\:MM\:SSZ)
- **level** field is one of [DEBUG, INFO, WARNING, ERROR, FATAL]
- **service** is a string
- **message** is a string
- **duration** is a duration written with a decimal value, positive or negative, rounded or floatting, followed by one of 'ns', 'us', 'ms', 's', 'm', 'h'
- **other** fields are optionnal and treated as extras. There can be zero to many other fields, treated as key-value pairs data

### Zebra format
```c
DATE LEVEL service=SERVICE message=MESSAGE duration=DURATION [other=EXTRAS]
```

### JSON format
```json
{"date": "DATE", "level": "LEVEL", "message": "MESSAGE", "duration": "DURATION", ["other": "EXTRAS"]}
```

## Server

When started, Zebra starts an HTTP server, waiting for requests
Currently available endpoints are:  

- `/parse`  
starts to parse designated file, using `file` query parameter.  
`json` query parameter is used to enable log parsing. truthy values are **true** and **1**. no value or falsy value will result in default format parsing.  
this endpoint will not return anything other than HTTP code indicating success or failure, but will start parsing job.  
currently, parsed logs are stored in RAM, and can be retrieved using below endpoints.

- `/logs`  
returns parsed logs.  
several query parameters are available to filter logs:  
`start-date` (inclusive) - expected RFC3339 format  
`end-date` (inclusive) - expected RFC3339 format  
`level` - allowed values are `FATAL - ERROR - WARNING - INFO - DEBUG` case insensitive  
`service`  
currently, no sorting parameters are possible. Logs will be returned using the following format:

```json
[
    {
      "date": "log_date",
      "level": "log_level",
      "duration": "duration_in_ns",
      "message": "log_message",
      "service": "log_service"
    }
]
```

- `/metrics`  
serves basic logs metrics. Analyses all logs with no filter available, and outputs the following format:  

```json
{
  "number_of_lines": {
    "service": "number_of_lines_for_this_service"
  },
  "service_performance": {
    "service_name": {
      "name": "service_name",
      "number_of_lines": "number_of_lines_for_this_service",
      "average_duration": "average_log_duration_for_this_service_is_ns",
    }
  },
  "file_errors": "number_of_files_errors",
  "parse_errors_count": "number_of_parse_errors_encountered",
}
```

- `/slowest_logs`  
retrieves n slowest logs, given by integer query parameter `number_of_logs`  
each log retrieved here use the same format as `/logs` logs response
