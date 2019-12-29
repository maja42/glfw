glfw
====

Package glfw experimentally provides a glfw-like API
with desktop (via glfw) and browser (via HTML5 canvas) backends.

It is used for creating a GL context and receiving events.

**Note:** This package was forked from https://github.com/goxjs/glfw to add concurrency support for desktop.
It is intended to be used with https://github.com/maja42/gl.

The original package is currently in development. The API is incomplete and may change.

Installation
------------

```bash
go get -u github.com/maja42/glfw
GOARCH=js go get -u -d github.com/maja42/glfw
```

Directories
-----------

| Path                                                               | Synopsis                                                           |
|--------------------------------------------------------------------|--------------------------------------------------------------------|
| [test/events](https://godoc.org/github.com/goxjs/glfw/test/events) | events hooks every available callback and outputs their arguments. |

License
-------

-	[MIT License](https://opensource.org/licenses/mit-license.php)
