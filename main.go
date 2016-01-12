package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
)

var pick []int
var special int
var startingCash int64
var cash int64
var winnings int64
var jackpot int64
var history History
var structure LottoStructure
var paytable Paytable
var r *rand.Rand
var wonJackpot bool

func main() {
	paytable = Paytable{
		Payouts: map[Match]int64{
			Match{0, true}:  4,
			Match{1, true}:  4,
			Match{2, true}:  7,
			Match{3, false}: 7,
			Match{3, true}:  100,
			Match{4, false}: 100,
			Match{4, true}:  50000,
			Match{5, false}: 1000000,
			Match{5, true}:  -1,
		},
	}

	structure = LottoStructure{
		69,
		26,
		1,
	}

	fmt.Println("Lotto Sim")
	app := cli.NewApp()
	app.Name = "LottoSim"
	app.Usage = "lottoSim -p 1 2 3 4 5 [4]"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "cash",
			Value: 10000,
			Usage: "cash to play with",
		},
		cli.IntFlag{
			Name:  "n1",
			Value: 1,
			Usage: "lotto numbers you want to play",
		},
		cli.IntFlag{
			Name:  "n2",
			Value: 2,
			Usage: "lotto numbers you want to play",
		},
		cli.IntFlag{
			Name:  "n3",
			Value: 3,
			Usage: "lotto numbers you want to play",
		},
		cli.IntFlag{
			Name:  "n4",
			Value: 4,
			Usage: "lotto numbers you want to play",
		},
		cli.IntFlag{
			Name:  "n5",
			Value: 5,
			Usage: "lotto numbers you want to play",
		},
		cli.IntFlag{
			Name:  "special",
			Value: 6,
			Usage: "special number you want to play",
		},
		cli.IntFlag{
			Name:  "jackpot",
			Value: 100000000,
			Usage: "proposed jackpot",
		},
	}

	app.Action = func(c *cli.Context) {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		pick = []int{c.Int("n1"), c.Int("n2"), c.Int("n3"), c.Int("n4"), c.Int("n5")}
		special = c.Int("special")
		cash = int64(c.Int("cash"))
		startingCash = cash
		winnings = int64(0)
		jackpot = int64(c.Int("jackpot"))
		history = History{
			Plays: 0,
			Wins: map[Match]int64{
				Match{0, true}:  0,
				Match{1, true}:  0,
				Match{2, true}:  0,
				Match{3, false}: 0,
				Match{3, true}:  0,
				Match{4, false}: 0,
				Match{4, true}:  0,
				Match{5, false}: 0,
				Match{5, true}:  0,
			},
		}
		rand.Seed(time.Now().Unix())
		for cash > 0 && !wonJackpot {
			play()
		}

		fmt.Printf("Plays: %s\n", humanize.Comma(history.Plays))
		fmt.Printf("Starting Cash: $%s\n", humanize.Comma(startingCash))
		fmt.Printf("Winnings: $%s\n", humanize.Comma(winnings))
		fmt.Printf("Remaining: $%s\n", humanize.Comma(cash))
		fmt.Printf("PB:   %s\n", humanize.Comma(history.Wins[Match{0, true}]))
		fmt.Printf("1+PB: %s\n", humanize.Comma(history.Wins[Match{1, true}]))
		fmt.Printf("2+PB: %s\n", humanize.Comma(history.Wins[Match{2, true}]))
		fmt.Printf("3:    %s\n", humanize.Comma(history.Wins[Match{3, false}]))
		fmt.Printf("3+PB: %s\n", humanize.Comma(history.Wins[Match{3, true}]))
		fmt.Printf("4:    %s\n", humanize.Comma(history.Wins[Match{4, false}]))
		fmt.Printf("4+PB: %s\n", humanize.Comma(history.Wins[Match{4, true}]))
		fmt.Printf("5:    %s\n", humanize.Comma(history.Wins[Match{5, false}]))
		fmt.Printf("5+PB: %s\n", humanize.Comma(history.Wins[Match{5, true}]))
	}

	app.Run(os.Args)
}

func play() {
	history.Plays = history.Plays + 1
	numbers := []int{}
	for len(numbers) <= 5 {
		num := r.Intn(structure.NumberRange)
		if !Contain(num, numbers) {
			numbers = append(numbers, num)
		}
	}
	cash = cash - structure.PricePerTicket
	powerball := r.Intn(structure.SpecialRange)
	matchesPowerBall := powerball == special
	matches := 0
	for i := 0; i < len(pick); i++ {
		if Contain(pick[i], numbers) {
			matches = matches + 1
		}
	}
	payout, ok := paytable.Payouts[Match{matches, matchesPowerBall}]

	if ok {
		if payout == -1 {
			wonJackpot = true
			payout = jackpot
		}

		cash = cash + payout
		winnings = winnings + payout
		history.Wins[Match{matches, matchesPowerBall}] = history.Wins[Match{matches, matchesPowerBall}] + 1
	}
}

func Contain(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type LottoStructure struct {
	NumberRange    int
	SpecialRange   int
	PricePerTicket int64
}

type Match struct {
	MatchingNumbers int
	MatchesSpecial  bool
}

type Paytable struct {
	Payouts map[Match]int64
}

type History struct {
	Wins  map[Match]int64
	Plays int64
}

func (h History) String() string {
	out := ""
	for key, val := range h.Wins {
		out = out + fmt.Sprintf("%s: %d Wins\n", key, val)
	}
	return out
}

func (m Match) String() string {
	out := fmt.Sprintf("%d", m.MatchingNumbers)
	if m.MatchesSpecial {
		out = out + "+ Powerball"
	}
	return out
}
