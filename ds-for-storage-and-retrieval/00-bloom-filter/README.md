# Bloom filter

## Starting values

```sh
➜ go run .
Elapsed time: 1.000442ms
Memory usage: 8000 bytes
False positive rate: 100.00%
```

## Implementation notes

### Storage

How do we want to store our bits? The naïve approach we were given uses an array of `uint64` numbers; specifically initialized with space for a thousand of them. At 8 bytes a pop, that array takes up 8000 bytes, or 8kB.

If we use a `big.Int` as a bitset, as [suggested in this SO answer](https://stackoverflow.com/a/53681508), we'll take up less space. For a `big.Int`, we get back `32` when we ask for `unsafe.Sizeof(b)` (where `var b big.Int`). We can use `SetBit` like `b.SetBit(&b, index, 1)` and then we can check the value by using `b.Bit(index)`.

How many bits should we use? We could use 2<sup>`m`</sup> bits, where `m` is the size of the returned hash (more on this later). The returned number from the hash function would be a number in the range from 0 to 2<sup>`m`</sup> - 1.

And how many bit arrays? Page 11 suggests we can get a lot of benefit from going with 2, as `p` goes from `p` to `sqrt(p)`. 
We later see, though, that it's more memory efficient to have overlapping filters, which is really one filter where we set multiple bits per object (still `k` number of hashing functions)

End of page 14: for a `p` of `0.0001` (chance of a false positive) and a `k` of `4`, the memory size should be `40 * S`, where `S` is the set size.

### Hashing

We'll need `k` hash functions. We need each of them to return different results for the same input. We'll be hashing strings, and want a binary number that's made up of `#` bits.
Ultimately, we need to set bits at the various indexes in our bitset.

Seeds can be determined at runtime, since the Bloom filter only exists in memory. If we wanted the filter to persist across restarts, we'd have to save our seeds as well.