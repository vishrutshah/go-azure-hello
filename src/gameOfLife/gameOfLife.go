 package main

 import (
     "fmt"
     "net/http"
     "html/template"
     "io/ioutil"
     "os"
     "github.com/Azure/azure-sdk-for-go/storage"
     "game"
     "cloud"
)

var width, height = 80, 30
var golTemplate, err = template.ParseFiles("../src/gol.html")

var g = game.Game {Board: make([][]bool, width),
    Neighbors: make([][]int, width),
    Width: width,
    Height: height,
    Generations: 1,
    CurrentLongest: false}
     
var longestGameGen int64
var current, longest, custom, container = "currentGame", "longestGame", "customGame", "games"

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
    text := game.WriteText(&g, "O", "·", "<p>")
    renderTemplate(w, text)
}

func newGol(w http.ResponseWriter, r *http.Request) {
    resetGame()
    http.Redirect(w, r, "/game/", http.StatusFound)
}

func redirect(w http.ResponseWriter, r *http.Request){
    http.Redirect(w, r, "/game/", http.StatusFound)
}

func longestGame(w http.ResponseWriter, r *http.Request) {
    //LoadFile(longest, &g)
    b,_ := cloud.InitStorage()
    _, _, text := cloud.LoadFileBlob(longest, container, b)
    cloud.FillBoard(&g, text)
    //CreateFile(&g, current)
    cloud.CreateFileBlob(current, container, b)
    cloud.FillBlob(current, container, text, b)
    
    http.Redirect(w, r, "/game/", http.StatusFound)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
    fmt.Println("hi there")
    file, _, err := r.FormFile("file")
    if err != nil {
        return
    }
    text, err := ioutil.ReadAll(file)
    b,_ := cloud.InitStorage()
    cloud.CreateFileBlob(custom, container, b)
    cloud.FillBlob(custom, container, &text, b)
    cloud.FillBoard(&g, &text)
    
    http.Redirect(w, r, "/game/", http.StatusFound)    
}

func resetGame() {
    b,_ := cloud.InitStorage()
    if g.Generations >= longestGameGen {
        //copy and paste currentGame to longestGame, overwrites
        //copyPasteFile(longest, current)
        cloud.CopyPasteFileBlob(longest, current, container, b)
    }
    g.CurrentLongest = false
    game.FillRandomBoard(&g)
    //CreateFile(&g, current)
    cloud.CreateFileBlob(current, container, b)
    cloud.FillGameBlob(&g, current, container, b)
}

func runGameEndless(b *storage.BlobStorageClient) {
    for {
        for game.IsAlive(&g) {
            game.RunGame(&g)            
            if g.Generations > longestGameGen {
                longestGameGen = g.Generations
                if g.CurrentLongest == false {
                    g.CurrentLongest = true
                    //copyPasteFile(longest, current)
                    b,_ := cloud.InitStorage()
                    cloud.CopyPasteFileBlob(longest, current, container, b)
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
    game.InitGame(&g)
    //fmt.Printf("You just browsed page (if blank you're at the root): \nWidth: %d Game height: %d", width, g.Height)
    //http.HandleFunc("/", handler)
    //http.ListenAndServe(":" + os.Getenv("HTTP_PLATFORM_PORT"), nil)
    
    //CreateFile(&g, current)    
    b,_ := cloud.InitStorage()
    cloud.CreateFileBlob(current, container, b)
    cloud.FillGameBlob(&g, current, container, b)
    
    go runGameEndless(b)
    http.HandleFunc("/", redirect)
    http.HandleFunc("/game/", drawGol)
    http.HandleFunc("/new/", newGol)
    http.HandleFunc("/longest/", longestGame)
    http.HandleFunc("/upload/", uploadFile)
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../src/css"))))
    
    port := os.Getenv("HTTP_PLATFORM_PORT")
    if port == "" {
        port = "8080"
    }    
 }
 
 /*
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "You just browsed page (if blank you're at the root): %s\nWidth: %dGame height: %d", 
        r.URL.Path[1:],
        width,
        g.Height)
}
*/