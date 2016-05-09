package main

import (
    "fmt"
    "net/http"
    "html/template"
    "game"
    "time"
)

var width, height = 10, 10
var golTemplate, err = template.ParseFiles("gol.html")
var g = game.Game { make([][]bool, width), make([][]int, width), width, height}


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
    text := game.WriteText(&g, "<p>")
    renderTemplate(w, text)
}

func runGame(g *game.Game){
    time.Sleep(1000 * time.Millisecond)
    game.UpdateBoard(g)
    fmt.Print(game.WriteText(g, "\n"))
    fmt.Println()
}

func main() {
    game.InitGame(&g)
    
    for game.IsAlive(&g) {
        runGame(&g)
    }
}