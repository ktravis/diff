# diff [![GoDoc](https://godoc.org/github.com/gin-gonic/gin?status.svg)](https://godoc.org/github.com/gin-gonic/gin)
Go implementation of patience diff, LCS, and merge

```shell
go get github.com/ktravis/diff
```

---

Patience diff was developed by Bram Cohen, see [here](http://alfedenzo.livejournal.com/170301.html) and
[here](http://bramcohen.livejournal.com/73318.html) for explanations.
The diffs produced are generally more human-readable, i.e.,

the patience diff:

```
 void func1() {
     x += 1
 }

+void functhreehalves() {
+    x += 1.5
+}
+
 void func2() {
     x += 2
 }
```

vs the traditional LCS-based diff:

```
 void func1() {
     x += 1
+}
+
+void functhreehalves() {
+    x += 1.5
 }
 
 void func2() {
     x += 2
 }
```
