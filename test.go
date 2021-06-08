package main

import (
	"fmt"
)



func main() {
	fmt.Println("ok")

	searchFile := "search"
	//BuildIndexFromLocal(searchFile,"./docs/")
	SearchFromIndex(searchFile,"基因",1,50)

}

