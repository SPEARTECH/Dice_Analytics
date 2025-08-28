package main

import (
	"C"
	"time"
	// "log"
	"encoding/json"
	// "flag"
	"fmt"
	"math/rand"
	"math"
)

func spinCalculate(userAmt float64, bet float64, number float64, winnings float64, history []float64) (float64, string, []float64) {
	allNumbers := make([]float64, 10000)
	for i := range allNumbers {
		allNumbers[i] = float64(i) / 100
	}
	result := allNumbers[rand.Intn(len(allNumbers)-1)]
	var status string
	if result > math.Floor(number) {
        userAmt += bet * (winnings - 1)
        history = append(history, math.Floor(userAmt*100)/100)
        status = "win"
	} else {
        userAmt -= bet
        history = append(history,   math.Floor(userAmt*100)/100)
        status = "loss"
	}
	if math.IsInf(userAmt, 0) {
        userAmt = 0  // Reset or handle as needed
    }


	return userAmt, status, history
}

func calculate(userBet float64, number float64, winnings float64, increase float64, rolls int, recoveryRolls int, maxLosingStreak int) ([]float64, map[string]float64, []int) {
	fmt.Printf("Simulating %d rolls over %.2f with %.2f%% increase and %.2f starting bet...\n", rolls, number, increase, userBet)
	history := make([]float64, 0)
	userAmt := 0.0
	bet := userBet
	var wins, losses int
	var winAmt, lossAmt, ratioAvg, ratio, percentage float64
	recoveryNumbers := make(map[float64]int)
	losingStreak := 0
	var streaks []int
	var status string

	for i := 0; i < rolls; i++ {
		userAmt, status, history = spinCalculate(userAmt, bet, number, winnings, history)
		// fmt.Println(result)
		if status == "loss" {
			losses++
			lossAmt += bet
			bet *= increase
			recoveryNumbers[bet] = recoveryRolls
			losingStreak++
			if losingStreak >= maxLosingStreak {
				bet = bet
				streaks = append(streaks, losingStreak)
				losingStreak = 0
			}
		} else {
			wins++
			winAmt += bet * (winnings - 1)
			if val, ok := recoveryNumbers[bet]; ok && val > 0 {
				recoveryNumbers[bet]--
			} else {
				bet = userBet
			}
			if losingStreak > 0 {
				streaks = append(streaks, losingStreak)
				losingStreak = 0
			}
		}
	}

	stats := map[string]float64{
		"percentage": percentage,
		"ratio_avg":  ratioAvg,
		"ratio":      ratio,
		"win_amt":    winAmt,
		"loss_amt":   lossAmt,
		"final_amt":  userAmt,
	}

	return history, stats, streaks
}

//export Calculate_go
func Calculate_go(userBet float64, number float64, winnings float64, increase float64, rolls int, recoveryRolls int, maxLosingStreak int) *C.char {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	history, stats, streaks := calculate(userBet, number, winnings, increase, rolls, recoveryRolls, maxLosingStreak)
    // # data, stats, streaks = calculate(user_bet, number, winnings, increase, rolls, recovery_rolls, max_losing_streak)

    // # bar_dict = {}
    // # for item in streaks:
    // #     if item in bar_dict:
    // #         bar_dict[item] = [bar_dict[item][0] + 1]
    // #     else:
    // #         bar_dict[item] = [1]

    // # bar_dict = dict(sorted(bar_dict.items()))
    // # del bar_dict[0]
 
    // # bar_chart = []
    // # streak_values = {}
    // # for field in bar_dict:
    // #     bar_chart.append({'name':str(field),'data':bar_dict[field]})
    // #     count = 0
    // #     streak_values[field] = []
    // #     for item in streaks:
    // #         if item == field:
    // #             streak_values[item].append(count)
    // #             count = 0
    // #         count += 1
    // #         # if count > len(streaks):
    // #         #     streak_values[item].append(count)
    // #         #     count = 0


    // # streak_distance = []
    // # for item in streak_values:
    // #     streak_distance.append({'name':str(item),'data':[round(sum(streak_values[item]) / len(streak_values[item]), 4)]})

    // # # Return the modified data as JSON
    // # # return jsonify({'chart': data,'stats':stats, 'bar_chart': bar_chart, 'streak_distance': streak_distance})
    // # return json.dumps({'chart': data,'stats':stats, 'bar_chart': bar_chart, 'streak_distance': streak_distance})

	result, err := json.Marshal(map[string]interface{}{"chart": history, "stats": stats, "streaks": streaks})
	if err != nil {
		return C.CString(fmt.Sprintf("Error: %v", err))
	}
	// fmt.Println("JSON Output:", string(result))  // Debugging line
	return C.CString(string(result))
}

func main() {

	// // Run the Calculate_go function
	// resultPtr := Calculate_go(1.0, 50.0, 2.0204, 1.0, 100, 1, 0)

	// // Convert the C string to a Go string to properly print it
	// result := C.GoString(resultPtr)

	// // Print the result
	// fmt.Println(result)
}
