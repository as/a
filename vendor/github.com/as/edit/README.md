
# edit 
Edit is an implementation of the Acme/Sam command language 

# usage
```
ed, _ := text.Open(text.NewBuffer())
ed.Insert([]byte("Removing vowels isnt the best way to name things"), 0)

cmd, _ := edit.Compile(",x,[aeiou],d")
cmd.Run(ed)

fmt.Printf("%s\n", ed.Bytes())
// Rmvng vwls snt th bst wy t nm thngs

```

# example
See example/example.go

# reference
Rob Pike pioneered structural regular expressions in the 1980s. The original implementations can be found in his Sam and Acme text editors. 

http://doc.cat-v.org/bell_labs/structural_regexps/

http://doc.cat-v.org/bell_labs/sam_lang_tutorial/

[![Go Report Card](https://goreportcard.com/badge/github.com/as/edit)](https://goreportcard.com/report/github.com/as/edit)

# appendix

This implementation now runs 10-1000 times faster

Benchmark before coalescing (2017.09.17)

```
goos: windows
goarch: amd64
BenchmarkChange128KBto64KB-4             	       1	3531749200 ns/op
BenchmarkChange128KBto128KB-4            	       1	3784740700 ns/op
BenchmarkChange128KBto128KBNest4x2x1-4   	       1	3642752800 ns/op
BenchmarkChange128KBto128KBx16x4x1-4     	       1	3589181900 ns/op
```

After coalescing (current)

```
goos: windows
goarch: amd64
pkg: github.com/as/edit
BenchmarkChange128KBto64KB-4                           2         530753100 ns/op
BenchmarkChange128KBto128KB-4                        200           6529711 ns/op
BenchmarkChange128KBto128KBNest4x2x1-4               200           6450888 ns/op
BenchmarkChange128KBto128KBx16x4x1-4                 200           6333687 ns/op
BenchmarkDelete128KB-4                            200000             11720 ns/op
BenchmarkDelete128KBx64-4                             20          93447760 ns/op
BenchmarkDelete128KBx8-4                              20          68008640 ns/op
BenchmarkDelete128KBx1-4                               5         263001000 ns/op
BenchmarkDelete256KBx1-4                               2         519648200 ns/op
BenchmarkDelete512KBx1-4                               1        1032597800 ns/op
```

