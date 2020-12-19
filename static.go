package fibers;

import (
    "github.com/gofiber/fiber/v2"
    "github.com/GeertJohan/go.rice"
    "path/filepath"
    "os"
    "crypto/md5"
    "log"
    "io"
    "fmt"
    "strings"
    "net/http"
)

var hash2file map[string]string;
var file2hash map[string]string;

var cacheControlStr = "public, max-age=1000"

func StaticMiddleware(fs *rice.Box, mount string) func(c *fiber.Ctx) error  {


    if !strings.HasSuffix(mount, "/") {
        mount += "/"
    }
    if !strings.HasPrefix(mount, "/") {
        mount = "/" + mount
    }

    log.Println("hashing static files");

    hash2file   = make(map[string]string);
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

        h := md5.New()
        if _, err := io.Copy(h, f); err != nil {
            log.Fatal(err)
        }


        hashed_path := fmt.Sprintf("%s%s-%x%s", mount, boxed_path, h.Sum(nil), filepath.Ext(boxed_path));
        direct_path := mount + boxed_path;

        log.Println(direct_path + " => " + hashed_path );

        hash2file[hashed_path] = boxed_path;
        file2hash[direct_path] = hashed_path;

        return nil
    })
    if err != nil {
        panic(err)
    }


    httpfs := fs.HTTPBox();

    return func(c *fiber.Ctx) error {

		// We only serve static assets on GET or HEAD methods
		method := c.Method()
		if method != fiber.MethodGet && method != fiber.MethodHead {
			return c.Next()
		}

        if !strings.HasPrefix(c.Path(), mount) {
            return c.Next();
        }


        // map hashed files
        path := c.Path();
        if fspath, ok := hash2file[path]; ok {
            path = fspath;
        }

        // Strip prefix
        path = strings.TrimPrefix(path, mount)
        if !strings.HasPrefix(path, "/") {
            path = "/" + path
        }


        var file http.File
        var stat os.FileInfo

		file, err = httpfs.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				return c.Status(fiber.StatusNotFound).Next()
			}
			return err
		}

		if stat, err = file.Stat(); err != nil {
			return err;
		}

		if stat.IsDir() {
			return fiber.ErrForbidden
		}

		modTime := stat.ModTime()
		contentLength := int(stat.Size())

        // Set Content Type header
		c.Type(getFileExtension(stat.Name()))

		// Set Last Modified header
		if !modTime.IsZero() {
			c.Set(fiber.HeaderLastModified, modTime.UTC().Format(http.TimeFormat))
		}

		if method == fiber.MethodGet {
            c.Set(fiber.HeaderCacheControl, cacheControlStr)
			c.Response().SetBodyStream(file, contentLength)
			return nil
		}
		if method == fiber.MethodHead {
			c.Request().ResetBody()
			// Fasthttp should skipbody by default if HEAD?
			c.Response().SkipBody = true
			c.Response().Header.SetContentLength(contentLength)
			if err := file.Close(); err != nil {
				return err
			}
			return nil
		}



        return c.Next();
    }
}

func Static(path string) string {
    return file2hash[path];
}

func getFileExtension(path string) string {
	n := strings.LastIndexByte(path, '.')
	if n < 0 {
		return ""
	}
	return path[n:]
}
