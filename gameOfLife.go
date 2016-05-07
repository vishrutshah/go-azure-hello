package main

import (
    //"fmt"
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

func runGame(g *game.Game, frameRate int){
    waitTime := 1000 / frameRate
    time.Sleep(time.Duration(waitTime) * time.Millisecond)
    game.UpdateBoard(g)
}

func main() {
    game.InitGame(&g)
    game.PrintBoard(&g)
    
    /*
    fmt.Println("First generation")
    fmt.Print(game.WriteText(&g, "\n"))
    game.UpdateBoard(&g)
    fmt.Println("Second generation")
    fmt.Print(game.WriteText(&g, "\n"))
    */
    
    
    http.HandleFunc("/game", drawGol)
	http.ListenAndServe(":8080", nil)
    
    for true {
        go runGame(&g, 30)
        http.HandleFunc("/game", drawGol)
    }
}