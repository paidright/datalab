## Shoppadocket

![Shoppadocket Logo](https://notbad.software/img/shoppadocket_logo.jpg "Shop A Docket company logo")

Shoppadocket takes a directory of CSV files and produces some metadata based on their contents.

The intended use case is to provide some data on which to base an intuition about whether the data is likely complete. It can also be used to provide a receipt to a client about what data we have ingested from them.

It provides the following data points:
* Total rows
* Column headers
* Number of unique values per column

### Usage
```
shoppadocket < input.csv > receipt.txt
```
