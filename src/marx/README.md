## Marx

> Father of the union movement

![Marx Logo](https://notbad.software/img/marx_logo.jpg "Stylised portrait of Karl Marx")

Marx takes a set of CSV files in a directory and creates a union (geddit?) of their contents.

By default it will slurp up all csv files in the current working directory and emit the output to stdout. You can overwrite these defaults with flags if you like:

Note that if you try and pipe to a file in the input directory Marx will try and union the output file, fail, and exit immediately.

```
./marx --input . > ./output/output.csv
```
