package main

import (
    "fmt"
    "net/http"
    "html/template"
    "io/ioutil"
    "os"
    "github.com/Azure/azure-sdk-for-go/storage"
)

var width, height = 80, 30
var golTemplate, err = template.ParseFiles("../src/gol.html")
var g = Game {Board: make([][]bool, width),
    Neighbors: make([][]int, width),
    Width: width,
    Height: height,
    Generations: 1,
    CurrentLongest: false}
var longestGameGen int64
var current, longest, container = "currentGame", "longestGame", "games"

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

func uploadFile(w http.ResponseWriter, r *http.Request) {
    fmt.Println("hi there")
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
    //LoadFile(longest, &g)
    b,_ := InitStorage()
    _, _, text := LoadFileBlob(longest, container, b)
    FillBoard(&g, text)
    //CreateFile(&g, current)
    CreateFileBlob(current, container, b)
    FillBlob(current, container, text, b)
    
    http.Redirect(w, r, "/game/", http.StatusFound)
}

func resetGame() {
    b,_ := InitStorage()
    if g.Generations >= longestGameGen {
        //copy and paste currentGame to longestGame, overwrites
        //copyPasteFile(longest, current)
        CopyPasteFileBlob(longest, current, container, b)
    }
    g.CurrentLongest = false
    FillRandomBoard(&g)
    //CreateFile(&g, current)
    CreateFileBlob(current, container, b)
    FillGameBlob(&g, current, container, b)
}

func runGameEndless(b *storage.BlobStorageClient) {
    for {
        for IsAlive(&g) {
            RunGame(&g)            
            if g.Generations > longestGameGen {
                longestGameGen = g.Generations
                if g.CurrentLongest == false {
                    g.CurrentLongest = true
                    //copyPasteFile(longest, current)
                    b,_ := InitStorage()
                    CopyPasteFileBlob(longest, current, container, b)
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
    b,_ := InitStorage()
    CreateFileBlob(current, container, b)
    FillGameBlob(&g, current, container, b)
    
    go runGameEndless(b)
    http.HandleFunc("/game/", drawGol)
    http.HandleFunc("/new/", newGol)
    http.HandleFunc("/longest/", longestGame)
    http.HandleFunc("/upload/", uploadFile)
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../src/css"))))
    http.ListenAndServe(":8080", nil)   
}