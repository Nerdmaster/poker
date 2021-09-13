# Poker

Super-fast poker hand evaluator for standard five-, six-, and seven-card hands
as well as the quirkier Omaha Hold 'em hands.

## Usage

Running `make` will compile `cmd/poker`, producing a working example in
`bin/poker`, which creates a random five-person game of Texas Hold 'Em and
displays winners and their hands.

### Low-level APIs

The low-level API is somewhat cryptic at first glance:

- Everything starts with a `CardList`, which is just a slice of `Card`s. Create
  a cardlist with a string, e.g., `var cards = poker.ParseCards("Ah Qh Kh Th Jh")`
- `cards.Evaluate` returns a number:
  - If you have five cards, they're evaluated as-is
  - If you have six or seven cards, all possible five-card permutations are
    computed to find the best score
  - The lower the number, the better the hand
- Call `poker.GetHandRank` on a hand's value to get its hand rank (e.g., Flush,
  Straight, Two pair, etc.)
- If you need to determine which five cards made up the best hand,
  `cards.BestHand` returns both a score and the combination of cards which had
  that score
- There are Omaha variations of these functions as well, since Omaha has
  somewhat unusual rules for how you have to use the four hole cards
  - `cards.EvaluateOmaha` returns the score for the best Omaha hand given four
    hole cards and five community cards
  - `cards.BestOmahaHand` is just like `poker.BestHand`, but with the same
    input as above: four hole cards and five community cards

Surprisingly, the `Best*` functions are only about 10% slower than their
`Evaluate*` counterparts, making them an excellent choice for any situation
where the caller might want to report more than simply the hand rank.

### High-level APIs

The high-level APIs, on the other hand, should be pretty easy to use. Create a
deck, deal cards into a hand / community card list, and evaluate things.

It's easier to show an example and its output than to bother explaining this
all here. Fortunately, Go really kicks ass at that kind of documentation, so
check out the runnable examples in
[the official docs](https://pkg.go.dev/github.com/Nerdmaster/poker#section-documentation).

Or just look at the source for the example file(s).

That said, here's some basic info at a glance:

- Create a Deck: `var deck = poker.NewDeck(rand.NewSource(time.Now().UnixNano()))`
  - Using `math/rand`, especially with `time.Now()` is **not secure**. This is
    a simple example. Use a real source of randomness for this! The point here
    is that my poker package *allows any random source*.
- Create an empty hand and add a card to it: `var hand = poker.NewHand(nil); deck.Deal(hand)`
- Or create a hand from a list of drawn cards: `var hand = poker.NewHand(deck.Draw(5))`
- Evaluate a hand: `var res, err = hand.Evaluate()`

The `Evaluate` method takes an optional list of community cards. If those are
present, the hand to evaluate may be two cards for Texas Hold 'Em rules or four
cards for Omaha Hold 'Em rules.

The `HandResult` instance (`res` in the above example) can give you the raw
score, hand rank, best five cards sorted in a human-readable manner, and can
describe the hand in a human-friendly way, such as "Full House, Fours Over
Twos".

## Performance

Five-card evaluation is blazing fast, but seven-card evaluation is actually
just a brute-force approach, looking at all *twenty-one* possible five-card
permutations within the seven-card hand, so it's obviously 21x the cost.

Some libraries pregenerate every possible seven-card hand to have essentially
instant evaluations. This could be done here, but it's not necessary for the
vast majority of use-cases, and makes a lot more sense to be done in a
consuming project than in this relatively low-level package. I don't like the
idea of offering up a library which requires you to first generate a huge
binary blob before being able to use it. For a poker game server or something,
sure, that makes sense, and can be generated on first run, but for a library
meant to be used in other projects? Nope.

---

On my local system, which is pretty fast, running Go 1.17.1:

```
pkg: github.com/Nerdmaster/poker
BenchmarkNewCard-16             1000000000               1.298 ns/op           0 B/op          0 allocs/op
BenchmarkNewCardString-16       373777461               31.84 ns/op            0 B/op          0 allocs/op
BenchmarkEvalFiveFast-16        1000000000               6.422 ns/op           0 B/op          0 allocs/op
BenchmarkEvaluateFive-16        1000000000               6.888 ns/op           0 B/op          0 allocs/op
BenchmarkEvaluateSeven-16       81180350               145.0 ns/op             0 B/op          0 allocs/op
BenchmarkBestHandSeven-16       77212557               150.0 ns/op             0 B/op          0 allocs/op
BenchmarkEvaluateOmaha-16       26755777               416.1 ns/op             0 B/op          0 allocs/op
BenchmarkBestOmahaHand-16       28537060               417.3 ns/op             0 B/op          0 allocs/op
```

Breaking it down:

- Over 150 million five-card hands evaluated per second
- 6.5 million seven-card hands per second
- Over 2 million nine-card Omaha hands per second

I looked over the best pure go package I was able to find
([chehsunliu/poker](https://github.com/chehsunliu/poker)) so I could do some
benchmarking against it. Because its benchmarks actually evaluate several hands
per loop, it was easiest to shim my package into their benchmarking suite in
order to get a one-to-one comparison.

*This was run just minutes after the prior one. Same hardware, same version of
Go.  Same bash process, even.*

```
pkg: github.com/chehsunliu/poker
BenchmarkFiveNerdmaster-16      124406692               48.14 ns/op
BenchmarkFivePoker-16           43291016               141.2 ns/op
BenchmarkFiveJoker-16             183366             32632 ns/op
BenchmarkSixNerdmaster-16       16036488               375.5 ns/op
BenchmarkSixPoker-16             4487058              1330 ns/op
BenchmarkSixJoker-16               36391            164378 ns/op
BenchmarkSevenNerdmaster-16      4684904              1263 ns/op
BenchmarkSevenPoker-16            567606             10406 ns/op
BenchmarkSevenJoker-16             10000            588856 ns/op
```

My implementation absolutely crushes the "Joker" approach, and handily beats
chehsunliu's:

- No memory is allocated on the heap in any cases. If you use a package to
  implement a long-running poker game server, this can be critical.
- The seven-card case is just over 8x faster than chehsunliu's, which makes
  sense given the simpler logic I use
- Even the five-card case, which I thought would be equivalent, is nearly 3x
  faster than chehsunliu's

## Caveat

I'm amazing, and we all know this. But I have to be very clear here: **I can't
take *any* credit for the five-card eval's speed**. The best I can claim is
that I scoured the web to find the fastest implementation and then ported C,
C#, and even Java implementations. Sadly, I can't properly offer attribution,
because they're just a mish-mash of things posted on forums, stackoverflow,
etc.

## Background

a.k.a., why build something like this when poker evaluators are already so
prolific?

First: I didn't build *any* of the low-level evaluation code. The five-card
logic totally baffles me. As mentioned above, it was cobbled together from a
variety of poker evaluation libraries I've been looking at and playing with.
Hell, half of my tests exist because I needed to be sure I implemented the
weird-ass evaluation properly.

Second, the existing ecosystem for pure Go options was almost nonexistent,
especially when you look for robust, well-tested, and performant options.

I looked at some C options to see if I could just use something via a call from
Go. I'm probably stupid, but I just couldn't find anything simple enough or
portable enough.

So then I looked for pure Go projects. Surely **somebody** already ported
high-performance poker eval to Go, right?

I ran across what seemed like the most promising project,
[chehsunliu/poker](https://github.com/chehsunliu/poker). But it has some flaws
that I just can't accept, even in a hobby project like this:

- It doesn't let you specify a randomization source, meaning it's never a good
  fit for anything remotely secure
- Its evaluation logic is "fast enough", but still slower than I'd expect
- The seven-card hand eval is really awful. It doesn't pregenerate the list of
  permutations, which means a really slow brute-force compared to what I
  trivially whipped up.
- It instantiates a global object for the "master" deck, which is just...  so
  wild to me. Putting a global into the code to speed up deck creation? But the
  slow part of building a deck is *shuffling*, not filling it! And then doing
  seven-card eval the way it does...? Weird choice for optimizations, man.

I started off thinking I'd submit a PR to improve performance, but the API
being what it is, particularly the inability to customize the randomization
source, was too off-putting. Thus Nerdmaster's poker project begun!
