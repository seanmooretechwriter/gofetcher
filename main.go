package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "time"
    "strings"
    "bufio"
)

func main() {
    start := time.Now()
    ch := make(chan string)
    for _, url := range os.Args[1:] {
        go fetch(url, ch)
    }
    for range os.Args[1:] {
        fmt.Println(<-ch)
    }
    fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
   
func fetch(url string, ch chan<- string) {
    start := time.Now()
    resp, err := http.Get(url)
    if err != nil {
        ch <- fmt.Sprint(err)
        return 
    }
    b, err  := ioutil.ReadAll(resp.Body)
    // todo: convert this to a regular expression
    fName := strings.ReplaceAll(url, "http://", "")
    fName =  strings.ReplaceAll(fName, ".", "")
    fName =  strings.ReplaceAll(fName, "www", "")
    fName = fName[:4]
    fName = fName + ".txt"
    f, e := os.Create(fName)
    if e != nil {
        fmt.Println("Houston, file problem: ", e)
    }
    w := bufio.NewWriter(f)
    w.WriteString(string(b))
    w.Flush()
    resp.Body.Close()
    if err != nil {
        ch <- fmt.Sprintf("while reading %s: %v", url, err)
        return 
    }
    secs := time.Since(start).Seconds()
    ch <- fmt.Sprintf("%.2fs %v bytes  %s", secs, len(b), url)
}

