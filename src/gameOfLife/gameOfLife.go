package main

import (
    "fmt"
    "net/http"
    "html/template"
    "game"
    //"time"
)

var width, height = 30, 30
var golTemplate, err = template.ParseFiles("gol.html")
var g = game.Game { make([][]bool, width), make([][]int, width), width, height, 1}
var longestGameGen int64


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
    http.Redirect(w, r, "/game/", http.StatusFound)
}

func runGameEndless() {
    for {
        for game.IsAlive(&g) {
            gen := game.RunGame(&g)
            if gen > longestGameGen {
                longestGameGen = g.Generations
                fmt.Println(longestGameGen)
            }
        }
        game.FillBoard(&g)
    }
}

func main() {
    longestGameGen = 0
    game.InitGame(&g)
    go runGameEndless()    
    http.HandleFunc("/game/", drawGol)
    http.HandleFunc("/new/", newGol)
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
    http.ListenAndServe(":8080", nil)   
}