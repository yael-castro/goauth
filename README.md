# GOAuth - Authorization server
[![Icon](./doc/images/banner.png)](https://github.com/yael-castro)

[![Go Report Card](https://goreportcard.com/badge/github.com/yael-castro/goauth)](https://goreportcard.com/report/github.com/yael-castro/goauth)

Authorization server based on the [Authorization Code Flow](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1) with the extension [Proof Key for Code Exchange (PKCE)](https://datatracker.ietf.org/doc/html/rfc7636) of the protocol [OAuth 2.0](https://datatracker.ietf.org/doc/html/rfc6749)


###### Optional features excluded
- Refresh tokens
- Redirect URL in the authorization response

<hr>

###### Architecture style explained
The architecture style used in this project is the most common layered architecture pattern
with a little changes.

```
internal
├── business    (business logic layer)
├── dependency  (manage dependencies)
├── handler     (presentation layer)
├── model       (data transfer objects, business objects, errors and enums)
└── repository  (persistence layer)
```
      
###### Required environment variables
<table>
    <tr>
        <th>Variable</th>
        <th>Required value</th>
        <th>Default value</th>
    </tr>
    <tr>
        <td>PORT</td>
        <td>Integer</td>
        <td>8080</td>
    </tr>
    <tr>
        <td>REDIS_HOST</td>
        <td>String</td>
        <td></td>
    </tr>
    <tr>
        <td>REDIS_PORT</td>
        <td>Integer</td>
        <td></td>
    </tr>
    <tr>
        <td>REDIS_USER</td>
        <td>String</td>
        <td></td>
    </tr>
    <tr>
        <td>REDIS_PASSWORD</td>
        <td>String</td>
        <td></td>
    </tr>
    <tr>
        <td>REDIS_DATABASE</td>
        <td>Integer</td>
        <td></td>
    </tr>
</table>