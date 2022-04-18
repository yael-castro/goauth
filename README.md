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
    <tr>
        <td>PRIVATE_RSA_KEY</td>
        <td>String</td>
        <td></td>
    </tr>
</table>

###### Create your own private RSA key
```dotenv
PRIVATE_RSA_KEY="$(openssl genrsa 1024)"
```

###### Try it!
Using the Testing profile of Dependency Injection as show below you can try the server with the following command
![main](./doc/images/main.png)
```shell
curl -X GET "http://localhost:8080/go-auth/v1/authorization?response_type=code&state=ABC&client_id=mobile&redirect_uri=http://localhost:8080/callback&code_challenge=aAbBcCdDeEfFgGhHiIjJlLmMnNoOpPqQrRsStTuUvVxX&code_challenge_method=PLAIN" -H "Content-Type: application/x-www-form-urlencoded" -H "Authorization: Basic $(printf contacto@yael-castro.com:yael.castro | base64)"
```