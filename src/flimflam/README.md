## FlimFlam

> Kinda like a schemer but a bit flimsier

```
      .~~~~'\~~\
     ;       ~~ \
     |           ;
 ,--------,______|---.
/          \-----'    \
'.__________'-_______-'
```

FlimFlam is designed to generate a Google BigQuery compatible schema from the header line of a CSV file. BigQuery has the option to auto-detect schemas but it attempts to be clever about integers, dates, etc. When you just want everything to be stringly typed, you've got a FlimFlam.

`input.csv`
```
foo,bar,baz
1,2,3
```

```
flimflam < input.csv
```

Will produce:

```
foo:STRING,bar:STRING,baz:STRING
```
