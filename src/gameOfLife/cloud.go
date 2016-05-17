package main

import (
    "github.com/Azure/azure-sdk-for-go/storage"
    "encoding/json"
    "io/ioutil"
    "game"
)

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

func FillGameBlob(g *game.Game, fileName, cont string, b *storage.BlobStorageClient) error{    
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

func FillBoard(g *game.Game, text *[]byte){
    json.Unmarshal(*text, g)
}