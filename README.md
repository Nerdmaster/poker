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

On my local system, which is pretty fast, this suite can do about 30 million
seven-card hands per second, and nearly a billion five-card hands per second:

```
pkg: github.com/Nerdmaster/poker
BenchmarkEvalFiveFast-16        837644409                6.70 ns/op
BenchmarkEvaluateFive-16        703266120                8.47 ns/op
BenchmarkEvaluateSeven-16       30624259               194 ns/op
```

This isn't directly comparable to chehsunliu/poker, because each benchmark loop
in that code is actually evaluating several hands.  In order to compare with
100% equivalence, I chose instead to add my code into the chehsunliu's
benchmark suite:

```
pkg: github.com/chehsunliu/poker
BenchmarkFiveNerdmaster-16      177566901               67.3 ns/op
BenchmarkFivePoker-16           70615275               168 ns/op
BenchmarkFiveJoker-16             359961             31131 ns/op
BenchmarkSixNerdmaster-16       23403450               507 ns/op
BenchmarkSixPoker-16             8662846              1376 ns/op
BenchmarkSixJoker-16               76149            157771 ns/op
BenchmarkSevenNerdmaster-16      7137416              1690 ns/op
BenchmarkSevenPoker-16           1000000             10886 ns/op
BenchmarkSevenJoker-16             21188            566397 ns/op
```

My implementation absolutely crushes those other two, with the worst-case being
an improvement of nearly 3x for five- and six-card hands, while seven-card
evaluations are over 7x faster.

I can't take much credit for this, though - most of the performance comes from
porting C, C#, and even Java implementations I found all over the net which I
can't properly give attribution to, because they're just a mish-mash of things
posted on forums, stackoverflow, etc.
