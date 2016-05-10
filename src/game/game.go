package game

import (
    "fmt"
    "math/rand"
    "time"
    "encoding/json"
    "io/ioutil"
    "os"
)

type Game struct {
    Board [][]bool
    Neighbors [][]int
    Width int
    Height int
    Generations int64
    CurrentLongest bool
}

func InitGame(g *Game) {
    for i := 0; i < g.Height; i++ {
        g.Board[i] = make([]bool, g.Width)
        g.Neighbors[i] = make([]int, g.Width)
	}
    FillBoard(g)
}

func FillBoard(g *Game){
    for i := 0; i < g.Height; i++{
        for j := 0; j < g.Width; j++ {
            rand.Seed((time.Now()).UnixNano() + int64(j) + int64(i))
            randy := rand.Int()
            if randy % 2 == 1 {
                g.Board[i][j] = false
            } else {
                g.Board[i][j] = true
            }
        }
    }
    g.Generations = 1
}

func RunGame(g *Game) {
    time.Sleep(1000 * time.Millisecond)
    UpdateBoard(g)
}


func PrintBoard(g *Game) {
    for i := 0; i < g.Height; i++ {        
		fmt.Printf("%t\n", g.Board[i])
	}
}


func WriteText(g *Game, aliveCells string, deadCells string, newline string) string{
    text := ""
    for i := 0; i < g.Height; i++{
        for j := 0; j < g.Width; j++ {
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

func CreateFile(g *Game, filename string) error {
    text, err := json.MarshalIndent(g, "", "    ")
    if err != nil {
        return err
    }
    
    _, err2 := os.Stat(filename)
    if os.IsExist(err2) {
        err3 := os.Remove(filename)
        if err3 != nil {
            return err3
        }
    }    
        
    return ioutil.WriteFile(filename, text, 0600) 
}

func LoadFile(filename string, g *Game){
    text, _ := ioutil.ReadFile(filename)   
    json.Unmarshal(text, g)
}

func UpdateBoard (g *Game) {
    //count neighbors
    count := 0
    var x1, y1, x2, y2 int
    for i := 0; i < g.Height; i++{
        for j := 0; j < g.Width; j++ {
                        
            if i != 0 {
                x1 = i - 1
            } else {
                x1 = i
            }
            
            if i != (g.Height - 1) {
                x2 = i + 1
            } else {
                x2 = i
            }
            
            if j != 0 {
                y1 = j - 1
            } else {
                y1 = j
            }
            
            if j != (g.Width - 1) {
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
    for i := 0; i < g.Height; i++{
        for j := 0; j < g.Width; j++ {
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
    for i := 0; i < g.Height; i++{
        for j := 0; j < g.Width; j++ {
            if g.Board[i][j] == true {
                return true
            }
        }
    }
    return false
}

