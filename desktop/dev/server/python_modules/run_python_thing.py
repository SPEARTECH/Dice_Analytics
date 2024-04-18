import numpy as np
import json

def spin_calculate(user_amt, bet, number, winnings, history):
    results = np.random.rand()  # Generate a random number between 0 and 1
    result = results * 100  # Scale the random number to be between 0 and 100

    win_mask = result > int(number)  # Mask for win cases
    user_amt += np.where(win_mask, bet * (winnings - 1), -bet)  # Update user amount based on win or loss
    outcome = np.where(win_mask, bet * (winnings - 1), -bet)  # Calculate the outcome for each spin
    status = np.where(win_mask, 'win', 'loss')  # Determine the status for each spin

    # Append the updated user amount to history
    history += list(user_amt)

    return user_amt, status, result, outcome, history

def calculate(user_bet, number, winnings, increase, rolls, recovery_rolls, max_losing_streak):
    print(f'Simulating {rolls} rolls over {number} with {increase}% increase and {user_bet} starting bet...')
    history = []
    user_amt = np.zeros(rolls)  # Initialize user amount array
    bet = user_bet
    wins = 0
    losses = 0
    percentage = 0
    win_amt = 0
    loss_amt = 0
    ratio = 0
    recovery_numbers = {}
    losing_streak = 0
    streaks = np.zeros(rolls)  # Initialize streaks array

    for i in range(rolls):
        user_amt[i], status, result, outcome, history = spin_calculate(user_amt[i], bet, number, winnings, history)

        if status == 'loss':
            losses += 1
            loss_amt += bet
            bet *= increase
            recovery_numbers[bet] = recovery_rolls
            losing_streak += 1 * recovery_rolls
            if losing_streak >= max_losing_streak * recovery_rolls and max_losing_streak != 0:
                bet = user_bet
                streaks[i] = losing_streak
                losing_streak = 0
        else:
            if losing_streak > 0:
                streaks[i] = losing_streak
                losing_streak = 0
            wins += 1
            win_amt += outcome
            if bet in recovery_numbers:
                recovery_numbers[bet] -= 1
                if recovery_numbers[bet] <= 0:
                    recovery_numbers[bet] = 0
                    bet = user_bet

    if losses != 0:
        percentage = round(wins / (wins + losses) * 100, 2)
        ratio = round(win_amt / loss_amt, 2)
    else:
        percentage = 100
        ratio = 1

    stats = {
        'percentage': percentage,
        'ratio': ratio,
        'win_amt': win_amt,
        'loss_amt': loss_amt,
        'final_amt': int(user_amt[-1]),
        'profitability_ratio': round((percentage + (ratio * 100)) / 100 - 1, 2),
    }

    return history, stats, streaks


def main(number, winnings, increase, user_bet, rolls, recovery_rolls, max_losing_streak):
    data, stats, streaks = calculate(user_bet, number, winnings, increase, rolls, recovery_rolls, max_losing_streak)

    # Process streaks for visualization
    streak_values, streak_counts = np.unique(streaks[streaks != 0], return_counts=True)
    bar_chart = [{'name': str(val), 'data': [count]} for val, count in zip(streak_values, streak_counts)]
    streak_distance = [{'name': str(val), 'data': [np.mean(np.where(streaks == val)[0])]} for val in streak_values]

    return json.dumps({'chart': data, 'stats': stats, 'bar_chart': bar_chart, 'streak_distance': streak_distance})
