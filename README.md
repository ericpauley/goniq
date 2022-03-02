I constantly need a version of `uniq` that doesn't require duplicates to be consecutive, so here it is.

# Installing

`go install github.com/ericpauley/goniq`

# Usage
Arguments mirror those of GNU `uniq`. Input and output are always on stdin/stdout.

```
Usage: goniq [-cdhiu] [parameters ...]
 -c, --count        prefix lines with occurence count
 -d, --repeated     only print duplicates
 -h, --help         show help
 -i, --ignore-case  ignore case
 -u, --unique       only print uniques
```