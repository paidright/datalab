## Trogdor

> I just wanted to call something Trogdor

![Trogdor Logo](https://notbad.software/img/trogdor_logo.jpg "Hand drawn picture of a muscly dragon breathing fire")

Trogdor takes a set of CSV files and re-arranges their columns.

This is particularly useful as part of a sorting pipeline. ie: you can trogdor something then run it through a dumb sort function to sort on the first column(s)

```
trogdor --columns baz,foo < input.csv
```

Will take input.csv:

```
foo,bar,baz
1,2,3
4,5,6
```

and turn it into:

```
baz,foo,bar
3,1,2
6,4,5
```
