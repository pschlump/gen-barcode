# Generate bar codes from the command line

A simple tool to generate bar codes in different formats from the command line.

# Examples

2 of 5

```
$ ./gen-barcode --text 00223 --format 2-of-5  --output output/2-of-5-00233.png
```

![2 of 5 barcode](https://github.com/pschlump/gen-barcode/blob/main/output/2-of-5-00233.png?raw=true)

```
$ ./gen-barcode --text 002287 --format 2-of-5 --interleaved --output output/2-of-5-002387.png --height 50 --width 200
```

![2 of 5 barcode](https://github.com/pschlump/gen-barcode/blob/main/output/2-of-5-0023387.png?raw=true)
