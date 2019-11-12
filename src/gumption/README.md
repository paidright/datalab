## Gumption

> It cleans things up

![Gumption Logo](https://notbad.software/img/gumption_logo.png "Picture of a tub of Gumption brand cleaning product")

Gumption contains a grab-bag of functions to perform general purpose cleaning tasks on CSV data.

It streams rows and processes lines in parallel. This keeps the memory usage low and throughput high. It's designed to handle large files.

The columns flag is a comma delimited list of column names to operate on. If this is left blank, Gumption will attempt to act on all columns if the operator permits it.

```
gumption --columns one,two --input /data --output /result --quiet
```

`--strip-leading-zeroes`
```
one,two
0000123,00abc
```
Becomes:
```
one,two
123,abc
```

`--unquote`
```
one,two
'123',abc
```
Becomes:
```
one,two
123,abc
```

`--commas-to-points`
```
one,two
"123,456",abc
```
Becomes:
```
one,two
123.456,abc
```

`--add-missing 999`
```
one,two
,abc
```
Becomes:
```
one,two
999,abc
```

`--replace-cell 123,xyz,456,qwe`
```
one,two
456,poi
123,abc
```
Becomes:
```
one,two
qwe,poi
xyz,abc
```

`--replace-cell-lookup 123,two,456,two`
```
one,two
456,poi
123,abc
789,xyz
```
Becomes:
```
one,two
poi,poi
abc,abc
789,xyz
```

`--replace-char -,_,;,:`
```
one,two
-,poi
;,abc
```
Becomes:
```
one,two
_,poi
:,abc
```

`--rename asd --columns two`
```
one,two
123,abc
```
Becomes:
```
one,asd
999,abc
```

`--split-on-delim |`
```
one,two
123|456,abc
```
Becomes:
```
one,asd,one_1
123,abc,456
```

`--cp --columns one`
```
one,two
123,abc
```
Becomes:
```
one,asd,one_1
123,abc,123
```

`--drop --columns three`
```
one,two,three
123,abc,3
```
Becomes:
```
one,asd
123,abc
```

`--stomp-alphas --columns two`
```
one,two,three
123,1,abc
123,a,xyz
123,ab2a,abc
```
Becomes:
```
one,two,three
123,1,abc
123,,xyz
123,2,abc
```

`--delete-where --columns three`
```
one,two,three
123,1,abc
123,a,xyz
123,2,abc
```

Becomes:
```
one,two,three
123,1,abc
123,2,abc
```

`--trim-whitespace`
```
one,two
 123 ,abc
1 23,abc
```
Becomes:
```
one,two
123,abc
1 23,abc
```
