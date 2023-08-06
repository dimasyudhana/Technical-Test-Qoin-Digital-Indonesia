package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type Game struct {
	GameID        int
	Pemain        []Pemain
	Round         int
	HasilPerRound []HasilPerRound
}

type Pemain struct {
	PemainID int
	Skor     int
	Dadu     Dadu
}

type Dadu struct {
	SliceDaduSementara []int
	JumlahDaduSaatIni  int
	JumlahDaduAngka6   int
	JumlahDaduAngka1   int
}

type HasilPerRound struct {
	SkorPerPemain []int
	JumlahDadu    []any
	NilaiDadu     [][]int
}

func main() {
	jumlahPemain := 0
	jumlahDadu := 0

	fmt.Println("Selamat datang di Lempar Dadu Game!")
	fmt.Print("Masukkan jumlah pemain: ")
	fmt.Scan(&jumlahPemain)
	fmt.Print("Masukkan jumlah dadu per pemain: ")
	fmt.Scan(&jumlahDadu)

	fmt.Printf("Pemain = %d, Dadu = %d\n", jumlahPemain, jumlahDadu)

	game := New(jumlahPemain, jumlahDadu)
	game.PlayRounds()
}

func New(jumlahPemain, jumlahDadu int) *Game {
	game := &Game{
		Pemain: make([]Pemain, jumlahPemain),
	}
	for i := 0; i < jumlahPemain; i++ {
		game.Pemain[i] = Pemain{
			PemainID: i + 1,
			Dadu: Dadu{
				JumlahDaduSaatIni: jumlahDadu,
			},
		}
	}
	return game
}

func (g *Game) PlayRounds() {
	round := 1
	g.HasilPerRound = []HasilPerRound{}
	for {
		fmt.Println("=======================")
		fmt.Printf("Giliran \x1b[1m%d\x1b[0m lempar dadu:\n", round)
		for i := range g.Pemain {
			g.GameID++
			g.Pemain[i].Dadu.SliceDaduSementara = g.lemparDadu(g.Pemain[i].Dadu.JumlahDaduSaatIni, &g.Pemain[i].Dadu)
			if len(g.Pemain[i].Dadu.SliceDaduSementara) != 0 {
				fmt.Printf("\tPemain #%d (%d): %v\n",
					g.Pemain[i].PemainID,
					g.Pemain[i].Dadu.JumlahDaduAngka6,
					g.formatHasilDaduSementara(g.Pemain[i].Dadu.SliceDaduSementara),
				)
			} else {
				fmt.Printf("\tPemain #%d (%d): _ (Berhenti bermain karena tidak memiliki dadu)\n",
					g.Pemain[i].PemainID,
					g.Pemain[i].Dadu.JumlahDaduAngka6,
				)
			}
		}

		maxPlayerIndex, minPlayerIndex := g.getMaxAndMinPlayerIndex()
		g.moveDiceOneFromMaxIndexToMinIndexPlayer(maxPlayerIndex, minPlayerIndex)
		roundResult := g.evaluateRound()
		g.HasilPerRound = append(g.HasilPerRound, roundResult)

		fmt.Println("Setelah evaluasi:")
		for i := range g.Pemain {
			g.dadu6(i)
			nextPlayerIndex := (i + 1) % len(g.Pemain)
			g.dadu1(i, nextPlayerIndex)
			if len(g.Pemain[i].Dadu.SliceDaduSementara) != 0 {
				fmt.Printf("\tPemain #%d (%d): %v\n",
					g.Pemain[i].PemainID,
					g.Pemain[i].Dadu.JumlahDaduAngka6,
					g.formatHasilDaduSementara(g.Pemain[i].Dadu.SliceDaduSementara),
				)
			} else {
				fmt.Printf("\tPemain #%d (%d): _ (Berhenti bermain karena tidak memiliki dadu)\n",
					g.Pemain[i].PemainID,
					g.Pemain[i].Dadu.JumlahDaduAngka6,
				)
			}
		}

		remainingPlayers := 0
		var remainingPlayerID int
		for _, player := range g.Pemain {
			if len(player.Dadu.SliceDaduSementara) > 0 {
				remainingPlayers++
				remainingPlayerID = player.PemainID
			}
		}
		if remainingPlayers == 1 {
			fmt.Printf("Game berakhir karena hanya pemain #%d yang memiliki dadu.\n", remainingPlayerID)

			// Find the player with the highest score
			highestScore := 0
			winningPlayerID := 0
			for _, player := range g.Pemain {
				if player.Dadu.JumlahDaduAngka6 > highestScore {
					highestScore = player.Dadu.JumlahDaduAngka6
					winningPlayerID = player.PemainID
				}
			}

			fmt.Printf("Game dimenangkan oleh pemain #%d karena memiliki poin lebih banyak dari pemain lainnya.\n", winningPlayerID)
			break
		}

		round++
	}
}

func (g *Game) lemparDadu(jumlahDadu int, dadu *Dadu) []int {
	result := make([]int, jumlahDadu)
	for i := 0; i < jumlahDadu; i++ {
		lempar := rand.Intn(6) + 1
		switch lempar {
		case 6:
			dadu.JumlahDaduAngka6++
		case 1:
			dadu.JumlahDaduAngka1++
		}

		result[i] = lempar
	}
	return result
}

func (g *Game) filterHasilDaduSementara(dadu []int, f func(int) bool) []int {
	var result []int
	for _, nilaiDadu := range dadu {
		if f(nilaiDadu) {
			result = append(result, nilaiDadu)
		}
	}
	return result
}

func (g *Game) dadu6(playerIndex int) {
	g.Pemain[playerIndex].Dadu.SliceDaduSementara = g.filterHasilDaduSementara(g.Pemain[playerIndex].Dadu.SliceDaduSementara, func(d int) bool {
		return d != 6
	})
	g.Pemain[playerIndex].Dadu.JumlahDaduSaatIni = len(g.Pemain[playerIndex].Dadu.SliceDaduSementara)
	g.Pemain[playerIndex].Dadu.JumlahDaduAngka6 += g.Pemain[playerIndex].Dadu.JumlahDaduSaatIni - len(g.Pemain[playerIndex].Dadu.SliceDaduSementara)
}

func (g *Game) dadu1(currentPlayerIndex, nextPlayerIndex int) {
	currentPlayer := &g.Pemain[currentPlayerIndex]
	nextPlayer := &g.Pemain[nextPlayerIndex]
	if currentPlayer.Dadu.JumlahDaduAngka1 > 0 {
		currentPlayer.Dadu.SliceDaduSementara = g.removeOneValue(currentPlayer.Dadu.SliceDaduSementara, 1)
		currentPlayer.Dadu.JumlahDaduAngka1--
		nextPlayer.Dadu.SliceDaduSementara = append(nextPlayer.Dadu.SliceDaduSementara, 1)
		nextPlayer.Dadu.JumlahDaduSaatIni++
		g.dadu1(nextPlayerIndex, (nextPlayerIndex+1)%len(g.Pemain))
	}
}

func (g *Game) removeOneValue(slice []int, value int) []int {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (g *Game) formatHasilDaduSementara(dadu []int) string {
	daduStrings := make([]string, len(dadu))
	for i, nilaiDadu := range dadu {
		if nilaiDadu == 1 || nilaiDadu == 6 {
			daduStrings[i] = fmt.Sprintf("\x1b[1m%d\x1b[0m", nilaiDadu)
		} else {
			daduStrings[i] = fmt.Sprintf("%d", nilaiDadu)
		}
	}
	return strings.Join(daduStrings, ",")
}

func (g *Game) getMaxAndMinPlayerIndex() (maxIndex, minIndex int) {
	maxID := g.Pemain[0].PemainID
	minID := g.Pemain[0].PemainID

	for i := 1; i < len(g.Pemain); i++ {
		if g.Pemain[i].PemainID > maxID {
			maxID = g.Pemain[i].PemainID
			maxIndex = i
		}
		if g.Pemain[i].PemainID < minID {
			minID = g.Pemain[i].PemainID
			minIndex = i
		}
	}

	return maxIndex, minIndex
}

func (g *Game) moveDiceOneFromMaxIndexToMinIndexPlayer(maxPlayerIndex, minPlayerIndex int) {
	maxPlayer := &g.Pemain[maxPlayerIndex]
	minPlayer := &g.Pemain[minPlayerIndex]
	numDiceOne := maxPlayer.Dadu.JumlahDaduAngka1
	maxPlayer.Dadu.SliceDaduSementara = g.filterHasilDaduSementara(maxPlayer.Dadu.SliceDaduSementara, func(d int) bool {
		return d != 1
	})
	maxPlayer.Dadu.JumlahDaduAngka1 = 0
	minPlayer.Dadu.SliceDaduSementara = append(minPlayer.Dadu.SliceDaduSementara, g.generateDaduSatu(numDiceOne)...)
	minPlayer.Dadu.JumlahDaduSaatIni += numDiceOne
}

func (g *Game) generateDaduSatu(numDice int) []int {
	daduSatu := make([]int, numDice)
	for i := 0; i < numDice; i++ {
		daduSatu[i] = 1
	}
	return daduSatu
}

func (g *Game) evaluateRound() HasilPerRound {
	hasilPerRound := HasilPerRound{
		SkorPerPemain: make([]int, len(g.Pemain)),
		JumlahDadu:    make([]any, len(g.Pemain)),
		NilaiDadu:     make([][]int, len(g.Pemain)),
	}

	for i := range g.Pemain {
		hasilPerRound.SkorPerPemain[i] = g.Pemain[i].Dadu.JumlahDaduAngka6
		hasilPerRound.JumlahDadu[i] = g.Pemain[i].Dadu.JumlahDaduSaatIni
		hasilPerRound.NilaiDadu[i] = append([]int{}, g.Pemain[i].Dadu.SliceDaduSementara...)
	}

	return hasilPerRound
}
