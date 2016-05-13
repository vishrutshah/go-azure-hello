package main

import (
    //"fmt"
    "net/http"
    "html/template"
    //"game"
    "io/ioutil"
    "os"
    //"github.com/Azure/azure-sdk-for-go/storage" 
    //"time"
)

var width, height = 100, 50
var golTemplate, err = template.ParseFiles("../src/gol.html")
var g = Game {Board: make([][]bool, width),
    Neighbors: make([][]int, width),
    Width: width,
    Height: height,
    Generations: 1,
    CurrentLongest: false}
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
    text := WriteText(&g, "O", "·", "<p>")
    renderTemplate(w, text)
}

func newGol(w http.ResponseWriter, r *http.Request) {
    resetGame()
    http.Redirect(w, r, "/game/", http.StatusFound)
}

func longestGame(w http.ResponseWriter, r *http.Request) {
    LoadFile(longest, &g)
    CreateFile(&g, current)
    http.Redirect(w, r, "/game/", http.StatusFound)
}

func resetGame() {
    if g.Generations >= longestGameGen {
        //copy and paste currentGame to longestGame, overwrites
        copyPasteFile(longest, current)
    }
    g.CurrentLongest = false
    FillBoard(&g)
    CreateFile(&g, current)
}

func runGameEndless() {
    for {
        for IsAlive(&g) {
            RunGame(&g)            
            if g.Generations > longestGameGen {
                longestGameGen = g.Generations
                if g.CurrentLongest == false {
                    g.CurrentLongest = true
                    copyPasteFile(longest, current)
                }
            }
        }
        resetGame()
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
    longestGameGen = 1
    InitGame(&g)
    //CreateFile(&g, current)
    
    c,_ := InitStorage()    
    CreateFileBlob(&g, current, c)
    
    
    //go runGameEndless()
    http.HandleFunc("/game/", drawGol)
    http.HandleFunc("/new/", newGol)
    http.HandleFunc("/longest/", longestGame)
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../src/css"))))
    http.ListenAndServe(":8080", nil)   
}