## Stanley

> Stanley joins your data and makes it nice and flat

![Stanley Logo](https://notbad.software/img/stanley_logo.jpg "Picture of a Stanley No. 5 bench plane used in woodworking")

The `--left` file is buffered in memory. In the use case of this tool it's generally expected that this is a relatively small file mapping one variable to another.

The right side of the join is streamed from stdin. It may be arbitrarily large.

In this way the memory usage of Stanley is constant-ish in that it is constant yet bounded to the size of the map holding the data from `--left`

`left.csv`
```
id,foo
1,a
2,x
```

`right.csv`
```
id,bar
1,b
2,y
1,0
```

```
stanley --left left.csv --join-key id < right.csv
```

Will produce:

```
id,foo,bar
1,a,0
2,x,y
1,a,b
```
