# Poker

Super-fast poker hand evaluator for standard five-, six-, and seven-card hands
as well as the quirkier Omaha Hold 'em hands.

## Background

a.k.a., why build something like this when poker evaluators are already so
prolific?

First: I didn't build much.  Most of this is cobbled together from a variety of
poker evaluation libraries I've been looking at and playing with.

Second, the existing ecosystem for pure Go options was almost nonexistent,
especially when you look for robust, well-tested, and performant options.

I looked at some C options and then ran across
[chehsunliu/poker](https://github.com/chehsunliu/poker), and hoped to create a
faster version of that just by shimming in some of the C code.  Unfortunately,
the super-optimized C code I found is crazy-magic and not really compatible
with other approaches, so I instead tried to stick closely to chehsunliu's API
while replacing the core code.

And then I realized I didn't even really like those APIs, so I'm changing most
of it anyway.  For instance, using the insecure `math/rand` package for deck
shuffling rather than accepting any `rand.Source` implementation.

## Usage

Running `make` will compile `cmd/poker`, producing a working example in
`bin/poker`, which draws a random hand and outputs its rank and description.

The API is somewhat cryptic at first glance, but the system would be this fast
without a pretty low-level API.  Future plans include higher-level ways to get
at hand rankings.

- `poker.Evaluate` is the primary entrypoint for most use-cases.  Pass a `Card`
  slice in to `poker.Evaluate`, and get back a number
  - If you pass in five cards, they're evaluated as-is
  - If you pass in six or seven cards, all possible five-card permutations are
    computed to find the best score
  - The lower the number, the better the hand
- Call `poker.GetHandRank` on a hand's value to get its hand rank (e.g., Flush,
  Straight, Two pair, etc.)
- If you need to determine which five cards made up the best hand,
  `poker.BestHand` returns both a score and the combination of cards which had
  that score
- There are Omaha variations of these functions as well, since Omaha has
  somewhat unusual rules for how you have to use the four hole cards
  - `poker.EvaluateOmaha` returns the score for the best Omaha hand given four
    hole cards and five community cards
  - `poker.BestOmahaHand` is just like `poker.BestHand`, but with the same
    input as above: four hole cards and five community cards

Surprisingly, the `Best*` functions are only about 10% slower than their
`Evaluate*` counterparts, making them an excellent choice for any situation
where the caller might want to report more than simply the hand rank.

## Performance

This library doesn't do the "pregenerate all seven-card hands' values" approach
that makes some libraries insanely fast, though that could probably be added
without a ton of work.  So while five-card hand evaluation is more than fast
enough, seven-card evaluation is just a brute-forced solution: all 21 unique
five-card permutations are evaluated to get the seven-card hand's real rank.

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
