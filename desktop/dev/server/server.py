from http.server import HTTPServer, SimpleHTTPRequestHandler
import os
import sys
import random
from flask import Flask, render_template, render_template_string, request, jsonify
# import numpy as np
import random
import json

app = Flask(__name__)
# history = []

# Routes
@app.route('/')
def index():
    html = """
   
    """

    file_path = f'{os.path.dirname(os.path.realpath(__file__))}\\templates\\index.html'

    with open(file_path, 'r') as file:
        html = ''
        for line in file:
            html += line
            
        # return render('index.html')
        return render_template_string(html)

def spin_calculate(user_amt, bet, number, winnings, history):
    all_numbers = []
    for i in range(0,10000):
        all_numbers.append(i)
    status = ''
    result = all_numbers[random.randint(0,len(all_numbers)-1)]/100

    # Generate the array with numpy
    # all_numbers = np.arange(0, 10000) / 100

    # # Select a random element
    # result = np.random.choice(all_numbers)
    
    if result > int(number):
        user_amt = user_amt + (bet * (winnings - 1))
        history.append(user_amt)
        status = 'win'
        outcome = int(bet * (winnings - 1))
        # print(f"\t-${outcome} \tTotal: {int(user_amt)}.00\t {result}\t FIRST TWELVE\t Bet:{bet}")
    else:
        user_amt = user_amt - (bet)
        history.append(user_amt)
        status = 'loss'
        outcome = int(bet)
        # print(f"\t-${outcome} \tTotal: {int(user_amt)}.00\t {result}\t ZEROS\t\t Bet:{bet}")

    return user_amt, status, result, history

def calculate(user_bet, number, winnings, increase, rolls, recovery_rolls, max_losing_streak):
    print(f'Simulating {rolls} rolls over {number} with {increase}% increase and ${user_bet} starting bet...')
    history = []
    user_amt = 0
    bet = user_bet
    wins = 0
    losses = 0
    percentage = 0
    win_amt = 0
    loss_amt = 0
    ratio = 0
    recovery_numbers = {}
    # recovery_count = 0
    recovery_rolls = int(recovery_rolls)
    losing_streak = 0
    streaks = []
    for i in range(rolls):
        i = i + 1
        # sys.stdout.write(f"\Roll {i}/{rolls}")
        user_amt, status, result, history = spin_calculate(user_amt, bet, number, winnings, history)

        if status == 'loss':
            losses = losses + 1
            loss_amt = loss_amt + bet 
            bet = bet * increase       
            # if bet in recovery_numbers:
            #     recovery_numbers[bet] += recovery_rolls
            # else:
            #     recovery_numbers[bet] = recovery_rolls
            recovery_numbers[bet] = recovery_rolls
            # recovery_count = recovery_rolls
            losing_streak += 1 * recovery_rolls
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
            wins = wins + 1
            win_amt = win_amt + (bet * (winnings - 1)) 
            # recovery_count = recovery_count - 1
            # if recovery_count <= 0:
            #     bet = user_bet 
            #     recovery_count = 0
            if bet in recovery_numbers:
                recovery_numbers[bet] -= 1
                if recovery_numbers[bet] <= 0:
                    recovery_numbers[bet] = 0
                    # bet = bet / increase
                    # if bet <= user_bet:
                    #     bet = user_bet
                    bet = user_bet    
    if losses != 0:
        percentage = round(wins/(wins+(losses))*100,2)
        ratio_avg = round((win_amt/wins)/(loss_amt/losses),2)
        ratio = round((win_amt)/(loss_amt),2)
    else:
        percentage = 100
        ratio = 1
    # sys.stdout.flush()

    # print(f'\nWin %: {percentage}\nW/L ratio (Avg): {ratio_avg}/1.0\nW/L ratio (Total): {ratio}/1.0\nWin $: {win_amt}\nLoss $: {loss_amt}\nFinal $: ${int(user_amt)}')
    
    stats = {
        'percentage': percentage,
        'ratio_avg':ratio_avg,
        'ratio':ratio,
        'win_amt':win_amt,
        'loss_amt':loss_amt,
        'final_amt':int(user_amt),
        'profitability_ratio_avg':round((percentage + (ratio_avg * 100))/100-1,2),
        'profitability_ratio':round((percentage + (ratio * 100))/100-1,2),
    }

    return history, stats, streaks

@app.route('/api/example_api_endpoint', methods=['GET'])
def example_api_endpoint():
    # Get the data from the request
    # data = request.json.get('data') # for POST requests with data
    data = {'welcome to':'Raptor!'}

    # Perform data processing

    # Return the modified data as JSON
    return jsonify({'result': data})

@app.route('/api/simulate', methods=['POST'])
def simulate():
    # Get the data from the request
    # data = request.json.get('data') # for POST requests with data
    number = request.json['number']
    winnings = request.json['winnings']
    increase = request.json['increase']
    user_bet = request.json['user_bet']
    rolls = request.json['rolls']
    recovery_rolls = request.json['recovery_rolls']
    max_losing_streak = request.json['max_losing_streak']

    # Perform data processing

    data, stats, streaks = calculate(user_bet, number, winnings, increase, rolls, recovery_rolls, max_losing_streak)

    bar_dict = {}
    for item in streaks:
        if item in bar_dict:
            bar_dict[item] = [bar_dict[item][0] + 1]
        else:
            bar_dict[item] = [1]

    bar_dict = dict(sorted(bar_dict.items()))
    del bar_dict[0]

    bar_chart = []
    streak_values = {}
    for field in bar_dict:
        bar_chart.append({'name':str(field),'data':bar_dict[field]})
        count = 0
        streak_values[field] = []
        for item in streaks:
            if item == field:
                streak_values[item].append(count)
                count = 0
            count += 1
            # if count > len(streaks):
            #     streak_values[item].append(count)
            #     count = 0


    streak_distance = []
    for item in streak_values:
        streak_distance.append({'name':str(item),'data':[round(sum(streak_values[item]) / len(streak_values[item]), 4)]})

    # Return the modified data as JSON
    # return jsonify({'chart': data,'stats':stats, 'bar_chart': bar_chart, 'streak_distance': streak_distance})
    return json.dumps({'chart': data,'stats':stats, 'bar_chart': bar_chart, 'streak_distance': streak_distance})



if __name__ == '__main__':   
    pid_file = f'{os.path.expanduser("~")}/flask_server.pid'
    with open(pid_file, 'w') as f:
        f.write(str(os.getpid()))  # Write the PID to the file

    app.run(debug=True)    
    