package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
    "fmt"
    "io/ioutil"
    "os/exec"
	"math"
)

func main() {
	r := gin.Default()

	// Routes
	r.GET("/", index)
	r.POST("/api/simulate", simulate)

	// Start the server
	go openChrome("http://127.0.0.1:8080") // Open Chrome with your server URL
	r.Run(":8080")
}

func index(c *gin.Context) {
	// Load the HTML template
	htmlPath := ".\\templates\\index.html"
	htmlContent, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to load HTML template")
		return
	}

	// Render the HTML template
	c.Data(http.StatusOK, "text/html; charset=utf-8", htmlContent)
}

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
		if status == "loss" {
			losses++
			lossAmt += bet
			bet *= increase
			recoveryNumbers[bet] = recoveryRolls
			losingStreak++
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
		// fmt.Printf("Roll %d: bet=%.2f, status=%s, losingStreak=%d\n", i, bet, status, losingStreak) // Debug each roll
	}

	// Catch any remaining streak
	if losingStreak > 0 {
		streaks = append(streaks, losingStreak)
	}

	fmt.Println("Final Streaks:", streaks) // Debug final streak list

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

func simulate(c *gin.Context) {
	var requestData map[string]interface{}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userBet, _ := requestData["user_bet"].(float64)
	number, _ := requestData["number"].(float64)
	winnings, _ := requestData["winnings"].(float64)
	increase, _ := requestData["increase"].(float64)
	rolls, _ := strconv.Atoi(fmt.Sprintf("%.0f", requestData["rolls"]))
	recoveryRolls, _ := strconv.Atoi(fmt.Sprintf("%.0f", requestData["recovery_rolls"]))
	maxLosingStreak, _ := strconv.Atoi(fmt.Sprintf("%.0f", requestData["max_losing_streak"]))

	history, stats, streaks := calculate(userBet, number, winnings, increase, rolls, recoveryRolls, maxLosingStreak)

	// Debugging output for streaks
	fmt.Println("Streaks:", streaks)

	barDict := make(map[int]int)
	for _, item := range streaks {
		if item != 0 {
			barDict[item]++
		}
	}

	// Debugging output for barDict
	fmt.Println("BarDict:", barDict)

	var barChart []map[string]interface{}
	for key, value := range barDict {
		barChart = append(barChart, map[string]interface{}{
			"name": strconv.Itoa(key),
			"data": []int{value},
		})
	}

	// Debugging output for barChart
	fmt.Println("BarChart:", barChart)

	c.JSON(http.StatusOK, gin.H{
		"history":   history,
		"stats":     stats,
		"streaks":   streaks,
		"bar_chart": barChart,
	})
}

 
func openChrome(url string) {
	cmd := exec.Command("C:/Program Files/Google/Chrome/Application/chrome.exe", "--app=" + url, "--disable-pinch", "--disable-extensions", "--guest")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to open Chrome:", err)
	}
}
