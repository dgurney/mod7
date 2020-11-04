# DISCONTINUED

Further development will happen in [here for the generator library](https://github.com/dgurney/unikey), and [here for the cli version](https://github.com/dgurney/unikey-mod7)

# mod7
Key generator for any Microsoft products that use the mod7 algorithm. Caveat emptor: despite being functional, this codebase is rather old, so you may find some stupid code here and there.

## Usage
```
  -a    Generate all kinds of keys.
  -bench
        Benchmark generation and validation of keys.
  -bv string
        Batch validate a key file. The key file should be a plain text file (with a .txt extension) with 1 key per line.
  -d    Generate a 10-digit key (aka CD Key).
  -e    Generate an 11-digit CD key.
  -o    Generate an OEM key.
  -r int
        Generate n keys. (default 1)
  -t    Show how long the generation or batch validation took.
  -v string
        Validate a CD or OEM key.
  -ver
        Show version information and exit
```
