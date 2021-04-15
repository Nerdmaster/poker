# Poker

Poker is an implementation of a variety of poker evaluation libraries I've been
looking at and playing with.  The existing ecosystem for pure Go options was
basically nonexistent, so I looked at some C options and then tried to adapt
[chehsunliu/poker](https://github.com/chehsunliu/poker) to use more performant
APIs.  After some barriers, I found it far easier to just replicate that
project's API while doing most of the heavy lifting in my own custom code.
This removes  unnecessary dependencies in addition to the performance boosts
(this project has zero dependencies, while chehsunliu's has like half a dozen,
and one of them is only there for benchmarking purposes).

This library doesn't do the "pregenerate all seven-card hands' values" approach
that makes some libraries insanely fast, though that could probably be added
without a ton of work.  So while five-card hand evaluation is more than fast
enough, seven-card evaluation is just a brute-forced solution: all 21 unique
five-card permutations are evaluated to get the seven-card hand's real rank.

Running `make` will produce a working example in `bin/poker`, which draws a
random hand and outputs its rank and description.

## Performance

On my local system, which is pretty fast, this suite can do about 5 million
seven-card hands per second, and over 100 million five-card hands per second:

```
pkg: github.com/Nerdmaster/poker
BenchmarkEvalFiveFast-16        165811814                7.08 ns/op            0 B/op          0 allocs/op
BenchmarkEvaluateFive-16        134446173                8.73 ns/op            0 B/op          0 allocs/op
BenchmarkEvaluateSeven-16        5980243               195 ns/op               0 B/op          0 allocs/op
```

This isn't directly comparable to chehsunliu/poker, because each benchmark loop
in that code is actually evaluating several hands.  In order to compare with
100% equivalence, I chose instead to add my code into the chehsunliu's
benchmark suite:

```
pkg: github.com/chehsunliu/poker
BenchmarkFiveNerdmaster-16      173458536               68.7 ns/op             0 B/op          0 allocs/op
BenchmarkFivePoker-16           69460522               172 ns/op               0 B/op          0 allocs/op
BenchmarkFiveJoker-16             375366             32068 ns/op           14433 B/op        657 allocs/op
BenchmarkSixNerdmaster-16       22927119               520 ns/op               0 B/op          0 allocs/op
BenchmarkSixPoker-16             8416456              1419 ns/op             288 B/op          9 allocs/op
BenchmarkSixJoker-16               74277            161396 ns/op           67972 B/op       2923 allocs/op
BenchmarkSevenNerdmaster-16      6978400              1715 ns/op               0 B/op          0 allocs/op
BenchmarkSevenPoker-16           1000000             11050 ns/op            2304 B/op         72 allocs/op
BenchmarkSevenJoker-16             20720            579313 ns/op          265405 B/op      10231 allocs/op
```

My implementation absolutely crushes those other two:

- No memory is allocated on the heap in any cases
- Worst-case improvement is in the five-card case, at 2.5x faster
- The seven-card case is over 7x faster

I can't take much credit for this, though - most of the performance comes from
porting C, C#, and even Java implementations I found all over the net which I
can't properly give attribution to, because they're just a mish-mash of things
posted on forums, stackoverflow, etc.
