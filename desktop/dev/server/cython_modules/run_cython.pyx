import os
import sys
import random
import numpy as np
import random
import json


cimport numpy as np
from libc.stdlib cimport rand, RAND_MAX
from libc.math cimport floor
from libc.stdlib cimport atof, atoi

cpdef tuple spin_calculate(double user_amt, double bet, int number, double winnings, list history):
    cdef int array_size = 10000
    cdef np.ndarray[np.float64_t, ndim=1] all_numbers = np.arange(0, array_size, dtype=np.float64) / 100
    cdef int idx = rand() % array_size  # Safe way to generate index
    cdef double result = all_numbers[idx]

    if result > number:
        user_amt += bet * (winnings - 1)
        history.append(user_amt)
        return user_amt, 'win', result, history
    else:
        user_amt -= bet
        history.append(user_amt)
        return user_amt, 'loss', result, history

cpdef tuple calculate(double user_bet, int number, double winnings, double increase, int rolls, int recovery_rolls, int max_losing_streak):
    cdef int i, losses = 0, wins = 0
    cdef double user_amt = 0, bet = user_bet, win_amt = 0, loss_amt = 0, percentage, ratio, ratio_avg
    cdef list history = []
    cdef dict recovery_numbers = {}
    cdef list streaks = []
    cdef int losing_streak = 0
    cdef str status
    cdef double result

    for i in range(rolls):
        user_amt, status, result, history = spin_calculate(user_amt, bet, number, winnings, history)

        if status == 'loss':
            losses += 1
            loss_amt += bet
            bet *= increase
            recovery_numbers[bet] = recovery_rolls
            losing_streak += recovery_rolls
            if losing_streak >= max_losing_streak * recovery_rolls and max_losing_streak != 0:
                bet = user_bet
                streaks.append(losing_streak)
                losing_streak = 0
            else:
                streaks.append(0)
        else:
            if losing_streak > 0:
                streaks.append(losing_streak)
                losing_streak = 0
            else:
                streaks.append(0)
            wins += 1
            win_amt += bet * (winnings - 1)
            if bet in recovery_numbers:
                recovery_numbers[bet] -= 1
                if recovery_numbers[bet] <= 0:
                    recovery_numbers[bet] = 0
                    bet = user_bet

    if losses != 0:
        percentage = round(wins / (wins + losses) * 100, 2)
        ratio_avg = round(win_amt / wins / (loss_amt / losses), 2)
        ratio = round(win_amt / loss_amt, 2)
    else:
        percentage = 100
        ratio = 1

    cdef dict stats = {
        'percentage': percentage,
        'ratio_avg': ratio_avg,
        'ratio': ratio,
        'win_amt': win_amt,
        'loss_amt': loss_amt,
        'final_amt': int(user_amt),
        'profitability_ratio_avg': round((percentage + (ratio_avg * 100)) / 100 - 1, 2),
        'profitability_ratio': round((percentage + (ratio * 100)) / 100 - 1, 2),
    }

    return history, stats, streaks

cpdef str main(double number, double winnings, double increase, double user_bet, int rolls, int recovery_rolls, int max_losing_streak):
    # Assuming calculate returns Python objects which do not need typing in Cython
    data, stats, streaks = calculate(user_bet, number, winnings, increase, rolls, recovery_rolls, max_losing_streak)

    cdef dict bar_dict = {}
    cdef int item
    for item in streaks:
        if item in bar_dict:
            bar_dict[item].append(bar_dict[item][0] + 1)
        else:
            bar_dict[item] = [1]

    # Sorting dictionary by keys (converted to a list of tuples and sorted)
    bar_dict = {k: v for k, v in sorted(bar_dict.items())}

    # Assuming there might be a key '0' that needs deletion
    if 0 in bar_dict:
        del bar_dict[0]

    cdef list bar_chart = []
    cdef dict streak_values = {}
    cdef int count
    cdef int field
    for field in bar_dict:
        bar_chart.append({'name': str(field), 'data': bar_dict[field]})
        streak_values[field] = []
        count = 0
        for item in streaks:
            if item == field:
                streak_values[field].append(count)
                count = 0
            count += 1

    cdef list streak_distance = []
    for field in streak_values:
        if len(streak_values[field]) > 0:
            average_distance = sum(streak_values[field]) / len(streak_values[field])
        else:
            average_distance = 0
        streak_distance.append({'name': str(field), 'data': [round(average_distance, 4)]})

    # Converting the result to JSON is typically a Python-level operation
    return json.dumps({'chart': data, 'stats': stats, 'bar_chart': bar_chart, 'streak_distance': streak_distance})
