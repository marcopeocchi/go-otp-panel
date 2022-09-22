package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func splash(port string) {
	fmt.Println("                                                                          ")
	fmt.Println(" ██████╗ ████████╗██████╗       ██████╗  █████╗ ███╗   ██╗███████╗██╗     ")
	fmt.Println("██╔═══██╗╚══██╔══╝██╔══██╗      ██╔══██╗██╔══██╗████╗  ██║██╔════╝██║     ")
	fmt.Println("██║   ██║   ██║   ██████╔╝█████╗██████╔╝███████║██╔██╗ ██║█████╗  ██║     ")
	fmt.Println("██║   ██║   ██║   ██╔═══╝ ╚════╝██╔═══╝ ██╔══██║██║╚██╗██║██╔══╝  ██║     ")
	fmt.Println("╚██████╔╝   ██║   ██║           ██║     ██║  ██║██║ ╚████║███████╗███████╗")
	fmt.Println(" ╚═════╝    ╚═╝   ╚═╝           ╚═╝     ╚═╝  ╚═╝╚═╝  ╚═══╝╚══════╝╚══════╝")
	fmt.Println("                                                                          ")
	fmt.Printf("> Ready | Listening on port:%s                                   \n\n", port)
}

func staticHandler(route string, index bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		path := filepath.Clean(ctx.Request.URL.Path)
		if index {
			if path == "/" {
				path = "index.html"
			}
		}
		path = strings.TrimPrefix(path, "/")

		fmt.Println(path)

		file, err := vueFS.Open(path)
		if err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(path))
		ctx.Writer.Header().Set("Content-Type", contentType)

		if strings.HasPrefix(path, "assets/") {
			ctx.Writer.Header().Set("Cache-Control", "public, max-age=2592000")
		}

		stat, err := file.Stat()
		if err == nil && stat.Size() > 0 {
			ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
		}

		io.Copy(ctx.Writer, file)
	}
}
