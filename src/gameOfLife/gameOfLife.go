package main

import (
    "fmt"
    "net/http"
    "html/template"
    "game"
    "io/ioutil"
    "os"
    //"time"
)

var width, height = 5, 5
var golTemplate, err = template.ParseFiles("../src/gol.html")
var g = game.Game {Board: make([][]bool, width),
    Neighbors: make([][]int, width),
    Width: width,
    Height: height,
    Generations: 1}
var longestGameGen int64
var current, longest = "currentGame.txt", "longestGame.txt"

func serverError(w *http.ResponseWriter, err error){
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
        return
	}
}

func renderTemplate(w http.ResponseWriter, text string) {
	err := golTemplate.Execute(w, template.HTML(text))
	serverError(&w, err)
}

func drawGol(w http.ResponseWriter, r *http.Request) {
    text := game.WriteText(&g, "O", "_", "<p>")
    renderTemplate(w, text)
}

func newGol(w http.ResponseWriter, r *http.Request) {
    game.FillBoard(&g)
    game.CreateFile(&g, current)
    http.Redirect(w, r, "/game/", http.StatusFound)
}

func runGameEndless() {
    for {
        for game.IsAlive(&g) {
            fmt.Printf("Gen actual: %d\nGen del juego mas largo: %d\n", g.Generations, longestGameGen)
            game.RunGame(&g)
            
            if g.Generations > longestGameGen {
                longestGameGen = g.Generations
                //fmt.Println(longestGameGen)
            }
        }
        if g.Generations == longestGameGen {
            //copy and paste currentGame to longestGame, overwrites
            copyPasteFile(longest, current)
        }
        game.FillBoard(&g)
        game.CreateFile(&g, current)
    }
}

func copyPasteFile(destiny, source string) error {
    text, _ := ioutil.ReadFile(source)
    
    _, err := os.Stat(destiny)
    if os.IsExist(err) {
        err2 := os.Remove(destiny)
        if err2 != nil {
            return err2
        }
    }    
    
    return ioutil.WriteFile(destiny, text, 0600)
}

func main() {
    longestGameGen = 0
    game.InitGame(&g)
    //game.LoadFile("game.txt", &g)
    //game.PrintBoard(&g)
    game.CreateFile(&g, current)
    go runGameEndless()    
    http.HandleFunc("/game/", drawGol)
    http.HandleFunc("/new/", newGol)
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../src/css"))))
    http.ListenAndServe(":8080", nil)   
}