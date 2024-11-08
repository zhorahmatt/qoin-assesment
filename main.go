package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Player structure to hold dice and score
type Player struct {
	Points      int
	CurrentDice []int
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed for random number generator

	// Initialize the number of players and dice per player
	arg1 := os.Args[1]
	arg2 := os.Args[2]

	//convert to int
	numPlayers, err1 := strconv.Atoi(arg1)
	dicePerPlayer, err2 := strconv.Atoi(arg2)

	if err1 != nil {
		fmt.Printf("Error convertion '%s' to integer: %v\n", arg1, err1)
		return
	}
	if err2 != nil {
		fmt.Printf("Error convertion '%s' to integer: %v\n", arg2, err2)
		return
	}

	// Initialize players map to hold each player's dice and score
	players := make(map[int]*Player)
	for i := 1; i <= numPlayers; i++ {
		dice := make([]int, dicePerPlayer)
		for j := range dice {
			dice[j] = rollDice()
		}
		players[i] = &Player{Points: 0, CurrentDice: dice}
	}

	// Display the initial state before any turns are taken
	fmt.Println("==== Initial State ====")
	displayResults(players)

	// Game loop
	for {
		// Display the round start
		fmt.Println("==== Round Start ====")

		// Create a list of active players who have dice left
		activePlayers := getActivePlayers(players)

		// Temporary map to hold dice that are passed to other players during the round
		passedDice := make(map[int][]int)

		// Process each player who still has dice
		for i := 1; i <= numPlayers; i++ {
			// Skip players with no dice left
			if len(players[i].CurrentDice) == 0 {
				continue
			}

			player := players[i]
			// Determine next player who has dice (cycling through active players)
			nextPlayerID := findNextActivePlayer(i, activePlayers)

			// Store new dice for the player after evaluation (removing passed dice)
			var newDice []int

			// Evaluate each die for the player
			for _, die := range player.CurrentDice {
				switch die {
				case 1:
					// Pass a die to the next active player
					newDie := rollDice()
					passedDice[nextPlayerID] = append(passedDice[nextPlayerID], newDie)
					fmt.Printf("Player #%d passes a die to Player #%d: New die = %d\n", i, nextPlayerID, newDie)
				case 6:
					// Remove this die and increase the player's score
					player.Points++
				default:
					// Keep other dice
					newDice = append(newDice, die)
				}
			}

			// Set new dice after evaluation
			player.CurrentDice = newDice
		}

		// After all players have acted, add the passed dice to the next active player's list
		for playerID, dice := range passedDice {
			players[playerID].CurrentDice = append(players[playerID].CurrentDice, dice...)
		}

		// Re-roll all remaining dice for each player who still has dice
		for i := 1; i <= numPlayers; i++ {
			if len(players[i].CurrentDice) > 0 {
				for j := range players[i].CurrentDice {
					players[i].CurrentDice[j] = rollDice()
				}
			}
		}

		// Display the updated state after all players have had their turn
		displayResults(players)

		// Check for game-ending condition: only one player with dice
		if checkGameEnd(players) {
			break
		}
	}

	fmt.Println("Game Over")
	displayResults(players)
}

// Function to roll a 6-sided die
func rollDice() int {
	return rand.Intn(6) + 1
}

// Function to display the current game results
func displayResults(players map[int]*Player) {
	for i, player := range players {
		fmt.Printf("Player #%d (Points: %d, Dice: %v)\n", i, player.Points, player.CurrentDice)
	}
	fmt.Println()
}

// Function to check if only one player has dice left
func checkGameEnd(players map[int]*Player) bool {
	activePlayers := 0
	for _, player := range players {
		if len(player.CurrentDice) > 0 {
			activePlayers++
		}
	}
	return activePlayers <= 1
}

// Function to get the list of players who still have dice
func getActivePlayers(players map[int]*Player) []int {
	var activePlayers []int
	for i, player := range players {
		if len(player.CurrentDice) > 0 {
			activePlayers = append(activePlayers, i)
		}
	}
	return activePlayers
}

// Function to find the next active player
func findNextActivePlayer(currentPlayerID int, activePlayers []int) int {
	// Find the index of the current player in the active list
	for i, id := range activePlayers {
		if id == currentPlayerID {
			// Get the next player in the list, wrapping around if necessary
			if i+1 < len(activePlayers) {
				return activePlayers[i+1]
			} else {
				return activePlayers[0]
			}
		}
	}
	return currentPlayerID // Fallback, shouldn't happen
}
