package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")
		ctx.Header("Access-Control-Allow-Origin", "*")

		log.Println("start")
		w := ctx.Writer
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		r := execTail()

		go func() {
			buf := bufio.NewReader(*r)
			for {
				line, err := buf.ReadBytes('\n')
				if err != nil {
					continue
				}

				w.Write([]byte(fmt.Sprintf("data: %s\n\n", line)))
				w.Flush()

			}
		}()
		//r := execTail()

		//log.Println("start goroutine")
		//buf := bufio.NewScanner(*r)
		//for buf.Scan() {
		//	w.Write([]byte(fmt.Sprintf("data: %s\n\n", buf.Text())))
		//	w.Flush()

		//}

		<-ctx.Done()

	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func execTail() *io.ReadCloser {
	cmd := exec.Command("tail", "-f", "test.txt")
	r, w := io.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w

	go cmd.Run()
	readCloser := ioutil.NopCloser(r)
	return &readCloser
}
