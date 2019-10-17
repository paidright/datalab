## Colcat

> He concatenates columns

```
      /\_____/\
     /  o   o  \
    ( ==  ^  == )
     )_________(
    (  |Colin|  )
   ( (  )   (  ) )
  (__(__)___(__)__)
  G'day. Name's Col.
```

`input.csv`
```
hi,foo,bar
1,2,3
a,b,c
```

`targets.csv`
```
dest,sep,sources
baz,-,foo:bar
```

```
colcat --targetFile targets.csv < input.csv
```

Will produce:

```
hi,baz
a,b-c
1,2-3
```
