## Ducky

> It tapes things together

![Ducky Logo](https://notbad.software/img/ducky_logo.jpg "Picture of a giant rubber duck")

Ducky merges rows where the value of one cell is equal to the value of another on the line before it.

input.csv
```
id,start,end
one,9am,11am
one,11am,5pm
two,9am,11am
two,11am,5pm
```

ducky --match "id:id,end:start" < input.csv

Becomes
```
id,start,end,ducky_taped
one,9am,5pm,true
two,9am,5pm,true
```

Ducky expects input to be sorted and will only merge adjacent rows.

input.csv
```
id,start,end
one,9am,11am
one,11am,2pm
beep,bonk,bork
one,2pm,5pm
```

ducky --match "id:id,end:start" < input.csv

Becomes
```
id,start,end,ducky_taped
one,9am,2pm,true
beep,bonk,bork,false
one,2pm,5pm,false
```

If multiple adjacent lines result in a match, they will all be merged together.

input.csv
```
id,start,end
one,9am,11am
one,11am,2pm
one,2pm,5pm
```

ducky --match "id:id,end:start" < input.csv

Becomes
```
id,start,end,ducky_taped
one,9am,5pm,true
```

You can also specify that a column must match a literal value.

input.csv
```
id,paycode,start,end
one,foo,9am,11am
one,bar,11am,5pm
one,baz,9am,11am
one,quux,11am,5pm
```

ducky --match "id:id,end:start" --match-literal-right "paycode:bar"

Becomes
```
one,bar,9am,5pm,true
one,baz,9am,11am,false
one,quux,11am,5pm,false
```

input.csv
```
id,paycode,start,end
one,foo,9am,11am
one,bar,11am,5pm
one,baz,9am,11am
one,quux,11am,5pm
```

ducky --match "id:id,end:start" --match-literal-left "paycode:foo"

Becomes
```
one,bar,9am,5pm,true
one,baz,9am,11am,false
one,quux,11am,5pm,false
```

input.csv
```
id,start,end
one,9am,11am
one,12am,5pm
two,9am,11am
two,11am,5pm
```

ducky --match "id:id" --inverse-match "end:start"

Becomes
```
id,start,end
one,9am,5pm,true
two,9am,11am,false
two,11am,5pm,false
```

input.csv
```
id,start,end
one,9am,11am
one,11am,5pm
two,9am,11am
two,11am,never
```

ducky --match "id:id,end:start" --inverse-match-literal-right "end:never"

Becomes
```
id,start,end
one,9am,5pm,true
two,9am,11am,false
two,11am,never,false
```
