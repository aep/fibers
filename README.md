middleware for go fiber
-----------------------



# Static

cache busting with content-hashed filenames.
i.e. /static/foo.png is loaded as /static/foo-123123123.png

usage:

```go
app.Use(fibers.StaticMiddleware(rice.MustFindBox("static"), "/static"))
```

and in views:


```
<img src="{{ call .static "/static/foo.png" }}" >
```

pass .static to c.Render, since views.AddFunc seems broken?
https://github.com/gofiber/template/issues/77


