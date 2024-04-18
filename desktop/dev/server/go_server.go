package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
    "fmt"
    "io/ioutil"
    "os/exec"
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
 
func spinCalculate(userAmt float64, bet float64, number float64, winnings float64, history []float64) (float64, string, float64, []float64) {
	allNumbers := make([]float64, 10000)
	for i := 0; i < 10000; i++ {
		allNumbers[i] = float64(i) / 100
	}
	result := allNumbers[rand.Intn(10000)]
	var status string
	if result > number {
		userAmt += bet * (winnings - 1)
		history = append(history, userAmt)
		status = "win"
	} else {
		userAmt -= bet
		history = append(history, userAmt)
		status = "loss"
	}
	return userAmt, status, result, history
}

func calculate(userBet float64, number float64, winnings float64, increase float64, rolls int, recoveryRolls int, maxLosingStreak int) ([]float64, map[string]float64,  []int) {
    fmt.Printf("Simulating %d rolls over %.2f with %.2f%% increase and $%.2f starting bet...\n", rolls, number, increase, userBet)
    userAmt := 0.0
    bet := userBet
    wins := 0
    losses := 0
    var percentage float64
    var winAmt, lossAmt, ratio float64
    recoveryNumbers := make(map[float64]int)
    losingStreak := 0
    var streaks []int
    var history []float64
    for i := 0; i < rolls; i++ {
        i++
        userAmt, status, result, history := spinCalculate(userAmt, bet, number, winnings, history)
        // fmt.Println(userAmt)
        // fmt.Println(result)
        // fmt.Println(output)
        userAmt = userAmt
        result = result
        history = history

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
                if val > 0 {
                    recoveryNumbers[bet]--
                } else {
                    recoveryNumbers[bet] = 0
                    bet = userBet
                }
            }
        }
    }
    if losses != 0 {
        percentage = float64(wins) / float64(wins+losses) * 100
        ratio = winAmt / lossAmt
    } else {
        percentage = 100
        ratio = 1
    }

    stats := map[string]float64{
        "percentage":              percentage,
        "ratio":                   ratio,
        "win_amt":                 winAmt,
        "loss_amt":                lossAmt,
        "final_amt":               userAmt,
        "profitability_ratio_avg": (percentage) / 100 - 1,
        "profitability_ratio":     (percentage + ratio*100) / 100 - 1,
    }

    return  history, stats, streaks
    //return map[string]float64{}, []float64{}, []int{}
}

func simulate(c *gin.Context) {
	var requestData map[string]interface{}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Extract parameters from the request
	number, _ := requestData["number"].(float64)
	winnings, _ := requestData["winnings"].(float64)
	increase, _ := requestData["increase"].(float64)
	userBet, _ := requestData["user_bet"].(float64)
	rolls, _ := strconv.Atoi(fmt.Sprintf("%.0f", requestData["rolls"])) // Convert to string and then to integer
	recoveryRolls, _ := strconv.Atoi(fmt.Sprintf("%.0f", requestData["recovery_rolls"])) // Convert to string and then to integer
	maxLosingStreak, _ := strconv.Atoi(fmt.Sprintf("%.0f", requestData["max_losing_streak"])) // Convert to string and then to integer

	// Call the calculation function
	history, stats, streaks := calculate(userBet, number, winnings, increase, rolls, recoveryRolls, maxLosingStreak)

	// Calculate bar_chart
	barDict := make(map[int]int)
	for _, item := range streaks {
		if item != 0 {
			barDict[item]++
		}
	}

	var barChart []map[string]interface{}
	for key, value := range barDict {
		barChart = append(barChart, map[string]interface{}{
			"name": strconv.Itoa(key),
			"data": value,
		})
	}

	// Calculate streak_distance
	streakValues := make(map[int][]int)
	for _, item := range streaks {
		if item != 0 {
			streakValues[item] = append(streakValues[item], 1)
		}
	}

	var streakDistance []map[string]interface{}
	for key, value := range streakValues {
		var sum int
		for _, v := range value {
			sum += v
		}
		average := float64(sum) / float64(len(value))
		streakDistance = append(streakDistance, map[string]interface{}{
			"name": strconv.Itoa(key),
			"data": []float64{average},
		})
	}

	// Format and send the response
	c.JSON(http.StatusOK, gin.H{
		"history":         history,
		"stats":           stats,
		"streaks":         streaks,
		"bar_chart":       barChart,
		"streak_distance": streakDistance,
		// You can include additional data if needed
	})
}
 
func openChrome(url string) {
	cmd := exec.Command("C:/Program Files/Google/Chrome/Application/chrome.exe", "--app=" + url, "--disable-pinch", "--disable-extensions", "--guest")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to open Chrome:", err)
	}
}
