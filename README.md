I constantly need a version of `uniq` that doesn't require duplicates to be consecutive, so here it is.

# Installing

`go get -u github.com/ericpauley/goniq`

# Usage
Arguments mirror those of GNU `uniq`. Input and output are always on stdin/stdout.

```
Usage: goniq [-cdhiu] [-f value] [-s value] [-w value]
 -c, --count        prefix lines with occurence count
 -d, --repeated     only print duplicates
 -f, --skip-fields=value
                    skip n fields
 -h, --help         show help
 -i, --ignore-case  ignore case
 -s, --skip-chars=value
                    skip characters (done after skipping fields)
 -u, --unique       only print uniques
 -w, --check-chars=value
                    compare max N characters per line
```