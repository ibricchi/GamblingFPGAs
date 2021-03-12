package server

import (
	"errors"
	"fmt"

	"github.com/chehsunliu/poker"
)

// The active game
var pokerGameStart gameStart
var pokerGame game

// Used when a game is open but not yet started.
type gameStart struct {
	open               bool
	players            []player
	initialPlayerMoney int
	smallBlindValue    int
}

/*
	gameCards is a slice of cards that will appear in the game.
 	currentRound is an integer between 1 and 4.
 	currentPlayer refers to the player slice index.
 	The player slice is sorted so that index 0 refers to the first player (small blind).
*/
type game struct {
	active                    bool
	hasEnded                  bool
	deck                      *poker.Deck
	communityCards            []poker.Card
	players                   []player
	currentRound              int
	currentPlayer             int
	lastBetAmountCurrentRound int
	smallBlindAmount          int
}

/*
	Expects a slice of players that only have the name attribute initialised.
 	All other attributes will be overriden.
*/
func initGame(players []player, initialPlayerMoney int, smallBlindAmount int) (game, error) {
	if len(players) < 2 {
		return game{}, errors.New("server: poker: Need at least 2 players to play a game")
	}

	deck := poker.NewDeck()

	// Give each player two cards and initial money
	for i := range players {
		players[i].Hand = deck.Draw(2)
		players[i].MoneyAvailableAmount = initialPlayerMoney
	}

	// Determine which other cards will appear in game
	communityCards := deck.Draw(5)

	players = sortPlayersAccordingToRandomBlind(players)

	players = allocateRelativeCardScores(players, communityCards)

	return game{
		active:                    true,
		hasEnded:                  false,
		deck:                      deck,
		communityCards:            communityCards,
		players:                   players,
		currentRound:              1,
		currentPlayer:             0,
		lastBetAmountCurrentRound: 0,
		smallBlindAmount:          smallBlindAmount,
	}, nil
}

/*
	Go to next player.
	Go to next round if last player of this round.
*/
func (g *game) next() {
	if g.currentPlayer != len(g.players)-1 {
		g.currentPlayer++
	} else if g.currentRound < 4 {
		g.currentRound++
		g.lastBetAmountCurrentRound = 0
	} else {
		g.hasEnded = true
	}
}

func (g *game) updateWithFPGAData(player *player, data incomingFPGAData) error {
	player.ShowCardsMe = data.ShowCardsMe
	player.ShowCardsEveryone = data.ShowCardsEveryone

	if !data.IsActiveData {
		return nil
	}

	if data.NewTryPeak {
		return fmt.Errorf("server: poker: updateGameWithFPGAData: NewTryPeak not yet implemented!")
	}

	// Player can't do anything
	if player.AllInCurrentRound {
		g.next()
	}

	if !isMoveAnAvailableNextMove(data.NewMoveType) {
		return fmt.Errorf("server: poker: move %v is not one of the available moves", data.NewMoveType)
	}

	switch data.NewMoveType {
	case "fold":
		player.HasFolded = true
	case "check":
		// Do nothing
	case "bet":
		if err := player.bet(data.NewBetAmount); err != nil {
			return fmt.Errorf("server: poker: failed to place bet")
		}
	case "call":
		player.call()
	case "raise":
		if err := player.raise(data.NewBetAmount); err != nil {
			return fmt.Errorf("server: poker: failed to place raise")
		}
	}

	g.next()

	return nil
}