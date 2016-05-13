package main

import (
    "github.com/Azure/azure-sdk-for-go/storage"
    //"game"
    "encoding/json"
    "fmt"
)

var fileBoundary int64 = 512
var cont string = "games"

func InitStorage() (*storage.BlobStorageClient, error) {
    // get accountName, accountKey
    accountName := "hellosto"
    accountKey := "SlZ2qIXn+rcRmFtE5UkUYN8P/mAYMKo48wPNugPF2o5hWnOMWSR+VRP8qHhOO/7EJptBCQoLAObgj3gcPSZQhA=="
    
    client, err := storage.NewBasicClient(accountName, accountKey)
    if err != nil {
        return nil, err
    }
    
    blobStoClient := storage.Client.GetBlobService(client)
    
    return &blobStoClient, nil
}

func CopyPasteFileBlob(destiny, source string) {
    
}

func CreateFileBlob(g *Game, fileName string, b *storage.BlobStorageClient) error {
    //add a name to the container
    //cont := "games"
    _, err := storage.BlobStorageClient.CreateContainerIfNotExists(*b, cont, storage.ContainerAccessTypeBlob)
    if err != nil {
        return err
    }
    fmt.Println("Created container")
    
    fileExists, err := storage.BlobStorageClient.BlobExists(*b, cont, fileName)
    if err != nil {
        return err
    }
    fmt.Println("Asked if blob exists")
    
    text, err := json.MarshalIndent(g, "", "    ")
    if err != nil {
        return err
    }
    fmt.Println("Marshalled game")
    
    var fileSize int64
    
    fileSize = int64(len(text) + (g.Width * g.Height))
    fileSize /= fileBoundary
    fileSize++
    fileSize *= fileBoundary
    fmt.Printf("Filesize: %d\n", fileSize)
    
    
    //is the block blob the best option?
    /*
    if fileExists == false {        
        err := b.storage.CreateBlockBlob(cont, filename)
        if err != nil {
            return err
        }
    }
    */
    
    extraHeaders := make(map[string]string)
    extraHeaders["Content-Type"] = "text/plain"
    
    if fileExists == false {
        err := storage.BlobStorageClient.PutPageBlob(*b, cont, fileName, fileSize, nil)
        if err != nil {
            return err
        }
    }
    fmt.Println("Created blob page")
    
    /*
    //get blockid
    blockID := "yada-yada-yada"
    
    err = b.storage.PutBlock(cont, filename, "", text)
    if err != nil {
        return err
    }
    */
    
    /*
    prop, err := storage.BlobStorageClient.GetBlobProperties(*b, cont, fileName)
    if err != nil {
        return err
    }
    fmt.Printf("Content type: %s\n", prop.ContentType)
    */
    
    err2 := storage.BlobStorageClient.PutPage(*b, cont, fileName, 0, fileSize - 1, storage.PageWriteTypeUpdate, text, nil )
    if err2 != nil {
        return err2
    }
    fmt.Println("Filled the blob")
    
    
    
    
    
    
    return nil
}

/*
func loadFileBlob(g *game.Game, filename string, b cloudManage.BlobStorageClient) error{
    cont := "games"
    fileExists, err := b.storage.BlobExists(cont, filename)
    if err != nil {
        return err
    }
    
    reader, err := b.storage.GetBlob(cont, filename)
    if err != nil {
        return err
    }
    
    text := make([]byte)    
    _, err := reader.ReadFull(reader, text)
    if err != nil {
        return err
    }
    
    reader.Close()
    json.Unmarshal(text, g)
}
*/