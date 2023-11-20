# Hangman Game

This is a simple Hangman game implemented in Go. The game allows players to guess a word by suggesting letters or entering the complete word.

# Features

1. Game Initialization
When starting the game, users are prompted to choose a word file from a list of available options.
2. Word Display
The chosen word is displayed with underscores representing unrevealed letters.
3. Letter and Word Guessing
Players can make guesses by entering a single letter or a complete word.
The game updates the display with revealed letters when correct guesses are made.
4. Game Save and Restore
Players can save the game progress at any point by entering "STOP".
The game progress is saved to a file named "save.txt".
To resume a saved game, the player needs to use the "--startWith save.txt" flag.
5. Make guesses by entering a single letter or a complete word.
Continue playing until you guess the word correctly or run out of attempts.

# Game Controls

- 'cd main' to enter the main folder
- 'go run main.go' to run the program
- 'STOP' to stop the program
Enter a single letter or a complete word when prompted to make a guess.
Dependencies
This game relies on the terminal for user interaction, so make sure to run it in a compatible terminal environment.
Acknowledgments
This Hangman game was created with Go programming language.


Have fun playing Hangman!