package sour;

import (
    "github.com/gin-gonic/gin"
    "github.com/GeertJohan/go.rice"
    "os"
    "crypto/sha256"
    "log"
    "io"
    "fmt"
    "strings"
)

var hash2file map[string]string;
var file2hash map[string]string;

var cacheControlStr = "public, max-age=1000"


func StaticMount(r*gin.Engine, mount string, fs *rice.Box) {

    if !strings.HasSuffix(mount, "/") {
        mount += "/"
    }
    if !strings.HasPrefix(mount, "/") {
        mount = "/" + mount
    }

    log.Println("hashing static files");

    file2hash   = make(map[string]string);

    err := fs.Walk(".", func(boxed_path string, info os.FileInfo, err error) error {

        if info.IsDir() {
            return nil;
        }

        f, err := fs.Open(boxed_path)
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()

        h := sha256.New()
        if _, err := io.Copy(h, f); err != nil {
            log.Fatal(err)
        }


        hash  := fmt.Sprintf("%x", h.Sum(nil))
        http_path := mount + boxed_path;

        log.Println(http_path + " => " + hash);

        file2hash[http_path] = hash;

        return nil
    })
    if err != nil {
        panic(err)
    }


    r.StaticFS(mount, fs.HTTPBox());

}

func Static(path string) string {
    return path + "?sha256=" + file2hash[path];
}

