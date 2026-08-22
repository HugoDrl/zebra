# Zebra
<img src="./logo/zebra.svg" alt="zebra-logo" width="300px"/>
Zebra is a CLI tool that allows to read, parse and analyse log files to retrieve metrics.

## Behavior
The behavior of the tool is the following  
Zebra will read one or several log files, given by user. The tool will extract only supported format logs and aggregate them.  
Zebra also shows number of errors encountered during parsing (non existing files, format errors, etc) but does not natively give details on those errors.  
If the user is requesting slowest logs using --top flag, the tool will also give full log lines for slowest logs, depending on the number requested  

## Logs Values

- **date** field is date and time, RFC3339 format (YYYY\-MM\-DDTHH\:MM\:SSZ)
- **level** field is one of [DEBUG, INFO, WARNING, ERROR, FATAL]
- **service** is a string
- **message** is a string
- **duration** is a duration written with a decimal value, positive or negative, rounded or floatting, followed by one of 'ns', 'us', 'ms', 's', 'm', 'h'
- **other** fields are optionnal and treated as extras. There can be zero to many other fields, treated as key-value pairs data

### Zebra format
```c
DATE LEVEL service=SERVICE message=MESSAGE duration=DURATION other=[EXTRAS]
```

### JSON format
```json
{"date": "DATE", "level": "LEVEL", "message": "MESSAGE", "duration": "DURATION", ["other": "EXTRAS"]}
```

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
- **json**: enable json format parsing

## Output values
Zebra will output a json string containing following informations:
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
  "slowest_logs": [
    {
      "date": "log_date",
      "level": "log_level",
      "duration": "duration_in_ns",
      "message": "log_message",
      "service": "log_service"
    }
  ]
}
```
An example for a scanning log output using `--top 1` could be
```json
{
  "number_of_lines": {
    "warning": 2
  },
  "service_performance": {
    "database": {
      "name": "database",
      "number_of_lines": 2,
      "average_duration": 50000000
    }
  },
  "file_errors": null,
  "parse_errors_count": 0,
  "slowest_logs": [
    {
      "date": "2026-01-01T00:00:00Z",
      "level": "warning",
      "duration": 50000000,
      "message": "hello from db",
      "service": "database"
    }
  ]
}
```
