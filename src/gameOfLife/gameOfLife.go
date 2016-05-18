 package main

 import (
     "fmt"
     "net/http"
     "html/template"
     //"io/ioutil"
     "os"
     //"github.com/Azure/azure-sdk-for-go/storage"
     //"encoding/json"
     //"math/rand"
     //"time"
)

//GAME

type Game struct {
    Board [][]bool
    Neighbors [][]int
    Width int
    Height int
    Generations int64
    CurrentLongest bool
}
/*
func InitGame(g *Game) {
    for i := 0; i < g.Height; i++ {
        g.Board[i] = make([]bool, g.Width)
        g.Neighbors[i] = make([]int, g.Width)
	}
    FillRandomBoard(g)
}

func FillRandomBoard(g *Game){
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
    time.Sleep(10000 * time.Millisecond)
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
*/
//CLOUD
/*
func InitStorage() (*storage.BlobStorageClient, error) {
    accountName := "hellosto"
    accountKey := "SlZ2qIXn+rcRmFtE5UkUYN8P/mAYMKo48wPNugPF2o5hWnOMWSR+VRP8qHhOO/7EJptBCQoLAObgj3gcPSZQhA=="    
    client, err := storage.NewBasicClient(accountName, accountKey)
    if err != nil {
        return nil, err
    }    
    blobStoClient := storage.Client.GetBlobService(client)
    
    return &blobStoClient, nil      
}

func CopyPasteFileBlob(destiny, source, cont string, b *storage.BlobStorageClient) error{
    err := CreateFileBlob(destiny, cont, b)
    if err != nil {
        return err
    }
    
    _, err, text := LoadFileBlob(source, cont, b)
    if err != nil {
        return err
    }
    
    err1 := storage.BlobStorageClient.AppendBlock(*b, cont, destiny, *text, nil)
    if err1 != nil {
        return err1
    }
    
    return nil
}

func CreateFileBlob(fileName, cont string, b *storage.BlobStorageClient) error {
    _, err := storage.BlobStorageClient.CreateContainerIfNotExists(*b, cont, storage.ContainerAccessTypeBlob)
    if err != nil {
        return err
    }
        
    //as these are append blobs, and just one game is needed, each time a game is stored,
    //the previous blob should be completely deleted
    _, err1 := storage.BlobStorageClient.DeleteBlobIfExists(*b, cont, fileName, nil)
    if err1 != nil {
        return err1
    }
    
    err2 := storage.BlobStorageClient.PutAppendBlob(*b, cont, fileName, nil)       
    if err2 != nil {
        return err2
    }
    
    return nil
}

func FillGameBlob(g *Game, fileName, cont string, b *storage.BlobStorageClient) error{    
    text, err := json.MarshalIndent(g, "", "    ")
    if err != nil {
        return err
    }
    FillBlob(fileName, cont, &text, b)
    
    return nil
}

func FillBlob(fileName, cont string, text *[]byte, b *storage.BlobStorageClient) error {
    err := storage.BlobStorageClient.AppendBlock(*b, cont, fileName, *text, nil)
    if err != nil {
        return err
    }
    
    return nil
}


func LoadFileBlob(filename, cont string, b *storage.BlobStorageClient) (bool, error, *[]byte){
    fileExists, err := storage.BlobStorageClient.BlobExists(*b, cont, filename)
    if err != nil {
        return fileExists, err, nil
    } else if fileExists == false{
        return fileExists, nil, nil
    }
    
    reader, err := storage.BlobStorageClient.GetBlob(*b, cont, filename)
    if err != nil {
        return fileExists, err, nil
    }    
    text, err := ioutil.ReadAll(reader)
    if err != nil {
        return fileExists, err, nil
    }    
    reader.Close()
        
    return fileExists, nil, &text
}

func FillBoard(g *Game, text *[]byte){
    json.Unmarshal(*text, g)
}
*/

//MAIN

var width, height = 80, 30
var golTemplate, err = template.ParseFiles("../src/gol.html")

var g = Game {Board: make([][]bool, width),
    Neighbors: make([][]int, width),
    Width: width,
    Height: height,
    Generations: 1,
    CurrentLongest: false}
     
var longestGameGen int64
var current, longest, custom, container = "currentGame", "longestGame", "customGame", "games"
/*
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

func redirect(w http.ResponseWriter, r *http.Request){
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

func uploadFile(w http.ResponseWriter, r *http.Request) {
    fmt.Println("hi there")
    file, _, err := r.FormFile("file")
    if err != nil {
        return
    }
    text, err := ioutil.ReadAll(file)
    b,_ := InitStorage()
    CreateFileBlob(custom, container, b)
    FillBlob(custom, container, &text, b)
    FillBoard(&g, &text)
    
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
*/
 func main() {
    //longestGameGen = 1
    //InitGame(&g)
    //fmt.Printf("You just browsed page (if blank you're at the root): \nWidth: %d Game height: %d", width, g.Height)
    http.HandleFunc("/", handler)
    //http.ListenAndServe(":" + os.Getenv("HTTP_PLATFORM_PORT"), nil)
    
    //CreateFile(&g, current)    
    //b,_ := InitStorage()
    //CreateFileBlob(current, container, b)
    //FillGameBlob(&g, current, container, b)
    
    //go runGameEndless(b)
    //http.HandleFunc("/", redirect)
    //http.HandleFunc("/game/", drawGol)
    //http.HandleFunc("/new/", newGol)
    //http.HandleFunc("/longest/", longestGame)
    //http.HandleFunc("/upload/", uploadFile)
    //http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../src/css"))))
    //http.ListenAndServe(":8080", nil)   
    http.ListenAndServe(":" + os.Getenv("HTTP_PLATFORM_PORT"), nil)
 }
 
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "You just browsed page (if blank you're at the root): %s\nWidth: %dGame height: %d", 
        r.URL.Path[1:],
        width,
        g.Height)
}