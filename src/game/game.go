package game

import (
    "fmt"
    "math/rand"
)

type Game struct {
    Board [][]bool
    Neighbors [][]int
    Width int
    Height int
}

func InitGame(g *Game){
    for i := 0; i < g.Width; i++ {
        g.Board[i] = make([]bool, g.Height)
        g.Neighbors[i] = make([]int, g.Height)
	}
    FillBoard(g)
}

func FillBoard(g *Game){
    seedIndex := 0
    for i := 0; i < g.Width; i++{
        for j := 0; j < g.Height; j++ {
            rand.Seed(int64(seedIndex + i * j))
            randy := rand.Int()
            if randy % 2 == 1 {
                g.Board[i][j] = true
            } else {
                g.Board[i][j] = false
            }
            seedIndex++
        }
    }
}

func PrintBoard(g *Game) {
    for i := 0; i < g.Width; i++ {        
		fmt.Printf("%t\n", g.Board[i])
	}
}

func WriteText(g *Game, newline string) string{
    text := ""
    for i := 0; i < g.Width; i++{
        for j := 0; j < g.Height; j++ {
            if g.Board[i][j] == true {
                text += "O"
            } else {
                text += "-"
            }
        }
        text += newline
    }
    return text
}

func UpdateBoard (g *Game) {
    //count neighbors
    count := 0
    var x1, y1, x2, y2 int
    for i := 0; i < g.Width; i++{
        for j := 0; j < g.Height; j++ {
                        
            if i != 0 {
                x1 = i - 1
            } else {
                x1 = i
            }
            
            if i != (g.Width - 1) {
                x2 = i + 1
            } else {
                x2 = i
            }
            
            if j != 0 {
                y1 = j - 1
            } else {
                y1 = j
            }
            
            if j != (g.Height - 1) {
                y2 = j + 1
            } else {
                y2 = j
            }
            
            count = 0
            for k := x1; k <= x2; k++{
                for l := y1; l <= y2; l++ {
                    if g.Board[k][l] == true {
                        count++
                    }
                }
            }
            if g.Board[i][j] == true {
                count--
            }
            g.Neighbors[i][j] = count
        }
    }
    //reassign life
    for i := 0; i < g.Width; i++{
        for j := 0; j < g.Height; j++ {
            if g.Board[i][j] == true {
                if g.Neighbors[i][j] < 2 || g.Neighbors[i][j] > 3{
                    g.Board[i][j] = false
                }
            } else {
                if g.Neighbors[i][j] == 3 {
                    g.Board[i][j] = true
                }
            }
        }
    }
}


func IsAlive(g *Game) bool {
    for i := 0; i < g.Width; i++{
        for j := 0; j < g.Height; j++ {
            if g.Board[i][j] == true {
                return true
            }
        }
    }
    return false
}