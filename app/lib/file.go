package lib

import "os"

func WriteToTxt(str []byte) {
	f, _ := os.Create("./tmp.txt")
	f.Write(str)
	f.Close()
}
