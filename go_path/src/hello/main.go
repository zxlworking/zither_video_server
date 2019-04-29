package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("zxl--->")
	video, err := os.Open("C:\\zxl\\go_path\\src\\video_project\\video_file\\test.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer video.Close()

	http.ServeContent(w, r, "test.mp4", time.Now(), video)
}

func main() {
	//http.HandleFunc("/", ServeHTTP)
	//http.ListenAndServe(":8080", nil)
	temp := 0
	t1 := time.Now()

	fmt.Println(temp)
	t2 := time.Now()
	fmt.Println(t2.Sub(t1))
}
