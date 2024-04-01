package debug

import "fmt"

func Trace(s string) string {
    fmt.Println("START ", s) 
    return s
}

func Un(s string) { 
    fmt.Println("END ", s) 
}
