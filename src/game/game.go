package game

import (
    "fmt"
    "math/rand"
    "time"
)

type Game struct {
    Board [][]bool
    Neighbors [][]int
    Width int
    Height int
    Generations int64
}

func InitGame(g *Game) {
    for i := 0; i < g.Width; i++ {
        g.Board[i] = make([]bool, g.Height)
        g.Neighbors[i] = make([]int, g.Height)
	}
    FillBoard(g)
}

func FillBoard(g *Game){
    for i := 0; i < g.Width; i++{
        for j := 0; j < g.Height; j++ {
            rand.Seed((time.Now()).UnixNano())
            randy := rand.Int()
            if randy % 8 == 1 {
                g.Board[i][j] = false
            } else {
                g.Board[i][j] = true
            }
        }
    }
    g.Generations = 1
}

func PrintBoard(g *Game) {
    for i := 0; i < g.Width; i++ {        
		fmt.Printf("%t\n", g.Board[i])
	}
}

func WriteText(g *Game, aliveCells string, deadCells string, newline string) string{
    text := ""
    for i := 0; i < g.Width; i++{
        for j := 0; j < g.Height; j++ {
            if g.Board[i][j] == true {
                text += aliveCells
            } else {
                text += deadCells
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
    g.Generations++
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

func RunGame(g *Game) int64 {
    time.Sleep(1000 * time.Millisecond)
    UpdateBoard(g)
    return g.Generations
}