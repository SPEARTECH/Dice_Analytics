import sys
import os
 
# getting the name of the directory
# where the this file is present.
current = os.path.dirname(os.path.realpath(__file__))
 
# Getting the parent directory name
# where the current directory is present.
parent = os.path.dirname(current)
 
# adding the parent directory to 
# the sys.path.
sys.path.append(parent)
 
# now we can import the module in the parent
# directory.
print(os.getcwd())
import run_python

def main(number,winnings,increase,user_bet,rolls,recovery_rolls,max_losing_streak):
    response = run_python.main(number,winnings,increase,user_bet,rolls,recovery_rolls,max_losing_streak)
    return response
