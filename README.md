# GOAuth - Authorization server
[![Icon](./doc/images/banner.png)](https://github.com/yael-castro)

[![Go Report Card](https://goreportcard.com/badge/github.com/yael-castro/goauth)](https://goreportcard.com/report/github.com/yael-castro/goauth)

Authentication server based on the *Authorization Code Flow* with the extension *Proof Key for Code Exchange (PKCE)* of the protocol *OAuth 2.0*

\- Writed in Go using the standard library

<hr>

###### Architecture style explained
The architectura style used in this project is the most common layered architecture pattern.

```
internal
├── business    (business layer)
├── dependency  (manage dependencies)
├── handler     (presentation layer)
├── model       (objects)
└── repository  (persistence layer)
```
