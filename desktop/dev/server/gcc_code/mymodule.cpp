#include <iostream>
#include <vector>
#include <map>
#include <random>
#include <cmath>
#include <iomanip>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

double spin_calculate(double &user_amt, double bet, double number, double winnings, std::vector<double> &history) {
    std::vector<double> all_numbers(10000);
    for (int i = 0; i < 10000; ++i) {
        all_numbers[i] = static_cast<double>(i) / 100;
    }
    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> dis(0, 9999);
    double result = all_numbers[dis(gen)];
    if (result > number) {
        user_amt += bet * (winnings - 1);
        history.push_back(user_amt);
        return user_amt;
    } else {
        user_amt -= bet;
        history.push_back(user_amt);
        return user_amt;
    }
}

std::tuple<std::vector<double>, std::map<std::string, double>, std::vector<int>> calculate(double user_bet, double number, double winnings, double increase, int rolls, int recovery_rolls, int max_losing_streak) {
    std::cout << "Simulating " << rolls << " rolls over " << number << " with " << increase << "% increase and " << user_bet << " starting bet..." << std::endl;
    std::vector<double> history;
    double user_amt = 0;
    double bet = user_bet;
    int wins = 0;
    int losses = 0;
    double percentage = 0;
    double win_amt = 0;
    double loss_amt = 0;
    double ratio = 0;
    std::map<double, int> recovery_numbers;
    int losing_streak = 0;
    std::vector<int> streaks;
    for (int i = 0; i < rolls; ++i) {
        ++i;
        user_amt = spin_calculate(user_amt, bet, number, winnings, history);
        if (user_amt < 0) {
            losses += 1;
            loss_amt += bet;
            bet *= increase;
            recovery_numbers[bet] = recovery_rolls;
            losing_streak += 1 * recovery_rolls;
            if (losing_streak >= max_losing_streak * recovery_rolls && max_losing_streak != 0) {
                bet = user_bet;
                streaks.push_back(losing_streak);
                losing_streak = 0;
            } else {
                streaks.push_back(0);
            }
        } else {
            if (losing_streak > 0) {
                streaks.push_back(losing_streak);
                losing_streak = 0;
            } else {
                streaks.push_back(0);
            }
            wins += 1;
            win_amt += bet * (winnings - 1);
            if (recovery_numbers.find(bet) != recovery_numbers.end()) {
                --recovery_numbers[bet];
                if (recovery_numbers[bet] <= 0) {
                    recovery_numbers[bet] = 0;
                    bet = user_bet;
                }
            }
        }
    }
    if (losses != 0) {
        percentage = round((wins / (wins + losses)) * 100, 2);
        double ratio_avg = round((win_amt / wins) / (loss_amt / losses), 2);
        ratio = round(win_amt / loss_amt, 2);
    } else {
        percentage = 100;
        ratio = 1;
    }

    std::map<std::string, double> stats = {
        {"percentage", percentage},
        {"ratio_avg", ratio_avg},
        {"ratio", ratio},
        {"win_amt", win_amt},
        {"loss_amt", loss_amt},
        {"final_amt", user_amt},
        {"profitability_ratio_avg", round((percentage + (ratio_avg * 100)) / 100 - 1, 2)},
        {"profitability_ratio", round((percentage + (ratio * 100)) / 100 - 1, 2)}
    };

    return std::make_tuple(history, stats, streaks);
}

int main() {
    double number = 10.0;
    double winnings = 2.0;
    double increase = 1.5;
    double user_bet = 5.0;
    int rolls = 1000;
    int recovery_rolls = 5;
    int max_losing_streak = 10;

    auto [data, stats, streaks] = calculate(user_bet, number, winnings, increase, rolls, recovery_rolls, max_losing_streak);

    std::map<int, int> bar_dict;
    for (auto item : streaks) {
        if (bar_dict.find(item) != bar_dict.end()) {
            bar_dict[item] += 1;
        } else {
            bar_dict[item] = 1;
        }
    }

    std::vector<std::map<std::string, int>> bar_chart;
    std::map<int, std::vector<int>> streak_values;
    for (auto &field : bar_dict) {
        bar_chart.push_back({{"name", field.first}, {"data", field.second}});
        int count = 0;
        streak_values[field.first] = {};
        for (auto &item : streaks) {
            if (item == field.first) {
                streak_values[item].push_back(count);
                count = 0;
            }
            ++count;
        }
    }

    std::vector<std::map<std::string, double>> streak_distance;
    for (auto &item : streak_values) {
        double sum = 0.0;
        for (auto &val : item.second) {
            sum += val;
        }
        double avg = sum / item.second.size();
        streak_distance.push_back({{"name", item.first}, {"data", avg}});
    }

    json result;
    result["chart"] = data;
    result["stats"] = stats;
    result["bar_chart"] = bar_chart;
    result["streak_distance"] = streak_distance;

    std::cout << result.dump() << std::endl;

    return 0;
}
