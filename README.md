# Datalab

![Datalab Logo](https://notbad.software/img/datalab_logo.jpg "Picture of a scientist holding a beaker")

Datalab is a collection of UNIX-y tools designed for working with CSV data. Read the individual READMEs to find out what they each do, but they are all designed along the following princples:

* Speed. We want to be able to deal with gigabytes of data quickly.
* Constant memory usage. The RAM used by the suite should not grow unbounded. We're happy to pay a CPU penalty for this, but we want to be able to run on commodity hardware.
* UNIX-y. Take data on STDIN and emit data on STDOUT where possible.
* Parallel. Where possible given the constraints of the operation being performed, use all available cores on the given machine.

You can use these tools individually or stitch them together into an ETL pipeline with some simple bash scripting. eg:

```
#!/bin/bash

set -e

export DATA_PATH=./big_data

echo "foo,bar,lalala" > $DATA_PATH/colcat_targets.csv

# Union related datasets
marx --quiet --input $DATA_PATH/masterfiles \

  # Add missing values to START_DATE_TIME, END_DATE_TIME, SHIFTEND_DATETIME, SHIFTSTART_DATETIME - 01/JAN/1900 00:00:00
  | gumption --quiet --add-missing "01/JAN/1900 00:00:00" --columns START_DATE_TIME,END_DATE_TIME,SHIFTEND_DATETIME,SHIFTSTART_DATETIME \

  # Create WorkDayID from ACTUAL_DATE
  | gumption --quiet --split " " --columns ACTUAL_DATE \

  | gumption --quiet --rename WorkDayID --columns ACTUAL_DATE \

  # Inner Join in a payperiod ID to based on the workday date
  | stanley --quiet --left $DATA_PATH/pp_id_mapping.csv --join-key WorkDayID \

  # Bring ID fields to the beginning of the file
  | trogdor --quiet --columns "EMPLOYEE_NUMBER,PP ID" \

  # Concatenate PayPeriod_End_Date with PayPeriodID on '/'
  | colcat --quiet --target_file $DATA_PATH/colcat_targets.csv \

  # Replace '-' in COSTCENTER_NUMBER and COSTCENTER_NAME with "Not Provided"
  | gumption --quiet --replace-cell "-,Not Provided" --columns COSTCENTER_NUMBER,COSTCENTER_NAME \

  > $DATA_PATH/cleaned_masterfile.csv
```
