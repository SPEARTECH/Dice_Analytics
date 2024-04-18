package main

import (
	"C"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand" // Added import for the rand package
)

func spinCalculate(bet float64, number float64, winnings float64) (float64, string) {
	allNumbers := make([]float64, 10000)
	for i := range allNumbers {
		allNumbers[i] = float64(i) / 100
	}
	result := allNumbers[rand.Intn(len(allNumbers))]
	var status string
	if result > number {
		status = "win"
	} else {
		status = "loss"
	}
	return result, status
}

func calculate(userBet float64, number float64, winnings float64, increase float64, rolls int, recoveryRolls int, maxLosingStreak int) ([]float64, map[string]float64, []int) {
	fmt.Printf("Simulating %d rolls over %.2f with %.2f%% increase and %.2f starting bet...\n", rolls, number, increase, userBet)
	var history []float64
	userAmt := 0.0
	bet := userBet
	var wins, losses, winAmt, lossAmt, ratioAvg, ratio float64
	recoveryNumbers := make(map[float64]int)
	losingStreak := 0
	var streaks []int
	for i := 0; i < rolls; i++ {
		_, status := spinCalculate(bet, number, winnings)
		if status == "loss" {
			losses++
			lossAmt += bet
			bet *= increase
			recoveryNumbers[bet] = recoveryRolls
			losingStreak += recoveryRolls
			if losingStreak >= maxLosingStreak*recoveryRolls && maxLosingStreak != 0 {
				bet = userBet
				streaks = append(streaks, losingStreak)
				losingStreak = 0
			} else {
				streaks = append(streaks, 0)
			}
		} else {
			if losingStreak > 0 {
				streaks = append(streaks, losingStreak)
				losingStreak = 0
			} else {
				streaks = append(streaks, 0)
			}
			wins++
			winAmt += bet * (winnings - 1)
			if val, ok := recoveryNumbers[bet]; ok {
				recoveryNumbers[bet]--
				if val <= 0 {
					recoveryNumbers[bet] = 0
					bet = userBet
				}
			}
		}
	}
	var percentage float64 // Declaration moved here
	if losses != 0 {
		percentage = (wins / (wins + losses)) * 100
		ratioAvg = ((winAmt / wins) / (lossAmt / losses))
		ratio = winAmt / lossAmt
	} else {
		percentage = 100
		ratio = 1
	}
	stats := map[string]float64{
		"percentage":                 percentage,
		"ratio_avg":                  ratioAvg,
		"ratio":                      ratio,
		"win_amt":                    winAmt,
		"loss_amt":                   lossAmt,
		"final_amt":                  userAmt,
		"profitability_ratio_avg":    (percentage + (ratioAvg * 100)) / 100 - 1,
		"profitability_ratio":        (percentage + (ratio * 100)) / 100 - 1,
	}

	return history, stats, streaks
}

func main() {
	// Parse command-line arguments
	var (
		number          = flag.Float64("number", 5.0, "Number")
		winnings        = flag.Float64("winnings", 2.0, "Winnings")
		increase        = flag.Float64("increase", 1.5, "Increase")
		userBet         = flag.Float64("user_bet", 10.0, "User Bet")
		rolls           = flag.Int("rolls", 1000, "Rolls")
		recoveryRolls   = flag.Int("recovery_rolls", 3, "Recovery Rolls")
		maxLosingStreak = flag.Int("max_losing_streak", 5, "Max Losing Streak")
	)
	flag.Parse()

	// Calculate based on command-line arguments
	history, stats, streaks := calculate(*userBet, *number, *winnings, *increase, *rolls, *recoveryRolls, *maxLosingStreak)

	// Prepare and print the output
	result, err := json.Marshal(map[string]interface{}{"chart": history, "stats": stats, "streaks": streaks})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(result))
}
