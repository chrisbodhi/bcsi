---
title: Bloom Filters
---

<section>

In this exercise, you will implement and optimize a Bloom filter. You are given the following files to start:

- `main.go` instantiates a Bloom filter, tests its correctness, and measures its quality.
- `bloom.go` defines the Bloom filter interface and provides a trivial implementation. It even performs very well on speed and memory usage! But the 100% false positive rate leaves some room for improvement.

## Objectives {-}

Your goal is to implement a working (nontrivial!) Bloom filter and optimize the following three criteria in tandem:

- Speed
- Memory usage
- False positive rate

Some of these will be in conflict (for example, it's easy to reduce the false positive rate by using significantly more memory), so part of your task is to find a good tradeoff between the three criteria. When considering tradeoffs, it can be helpful to keep specific use cases in mind.

Assuming that `/usr/share/dict/words` (the file used for testing / measurement) is around 2mb, a good initial target to aim for is 100kb of memory and a 15% false positive rate, without worrying about speed.

## Resources {-}

Here are some resources you can consult when getting started:

- [Bloom Filters by Example](https://llimllib.github.io/bloomfilter-tutorial/)
- [Why Bloom filters work the way they do](http://www.michaelnielsen.org/ddi/why-bloom-filters-work-the-way-they-do/)

## Suggestions {-}

Here are some directions you might consider exploring:

- Which hash function(s) should you use?
- How many hash functions should you use?
- Should you use a separate bit vector for each hash function, or a single combined bit vector?
- What bit vector sizes work well (e.g. prime numbers)?

## Stretch Goals {-}

Finally, if you're feeling really adventurous and would like some stretch goals:

- Investigate the quantitative relationship between "number of elements added" and "false positive rate"
- Consider how you'd extend the Bloom filter idea to support cardinality and/or deletion
- Read about and experiment with [Cuckoo Filters](https://bdupras.github.io/filter-tutorial/)

