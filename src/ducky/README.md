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
id,start,end
one,9am,5pm
two,9am,5pm
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
id,start,end
one,9am,2pm
beep,bonk,bork
one,2pm,5pm
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
id,start,end
one,9am,5pm
```
