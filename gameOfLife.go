package main

import (
    "fmt"
    "net/http"
    "html/template"
)

var size = 10
var golTemplate, err = template.ParseFiles("gol.html")
var game = Game { make([][]bool, size),make([][]int, size)}

type Game struct {
    Board [][]bool
    Neighbors [][]int
}

func writeText(g *Game) string{
    text := ""
    for i := 0; i < size; i++{
        for j := 0; j < size; j++ {
            g.Board[i][j] = true
            if g.Board[i][j] == true {
                text += "O"
            } else {
                text += " "
            }
        }
        text += "<p>"
    }
    return text
}

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
    text := writeText(&game)
    renderTemplate(w, text)
}

func main() {
    for i := 0; i < size; i++ {
        game.Board[i] = make([]bool, size)
		fmt.Printf("%t\n", game.Board[i])
	}
    
    http.HandleFunc("/game", drawGol)
	http.ListenAndServe(":8080", nil)
}