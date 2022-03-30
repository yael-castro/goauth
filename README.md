# GOAuth - Authorization server
[![Icon](./doc/images/banner.png)](https://github.com/yael-castro)

[![Go Report Card](https://goreportcard.com/badge/github.com/yael-castro/goauth)](https://goreportcard.com/report/github.com/yael-castro/goauth)

Authorization server based on the [Authorization Code Flow](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1) with the extension [Proof Key for Code Exchange (PKCE)](https://datatracker.ietf.org/doc/html/rfc7636) of the protocol [OAuth 2.0](https://datatracker.ietf.org/doc/html/rfc6749)

<hr>

###### Architecture style explained
The architectura style used in this project is the most common layered architecture pattern.

```
internal
├── business    (business logic layer)
├── dependency  (manage dependencies)
├── handler     (presentation layer)
├── model       (data transfer objects, business objects, errors and enums)
└── repository  (persistence layer)
```
