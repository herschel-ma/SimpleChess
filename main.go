package main

import (
	"fmt"

	"github.com/herschel-ma/SimpleChess/ChineseChess"
)

func main() {
	if ok := ChineseChess.NewGame(); !ok {
		fmt.Println("游戏启动失败")
	}
	// res.FileToSlice("./res", "./ChineseChess")
}
