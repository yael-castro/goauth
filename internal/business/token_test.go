package business

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/yael-castro/goauth/internal/model"
	"reflect"
	"strconv"
	"testing"
	"time"
)

// privateKey is a RSA Private Key used in TestJWTGenerator_GenerateToken
const privateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCaRX0EDvuYvQu1+hPOx36wByGJDk1EfS/PtgUC8nSf+BBYppTS
MrlvLcvg2C64m2sz+I6iiJ9zCZElUT7WaiN7oYrTAFy2elHFq+CvoUXeRPwAgu5c
ypwO6Pbjs3ebYknaCZHAL500DN1D0SNJjni/sviUIy5OZAvOs97NMlTXxQIDAQAB
AoGAMN7ceKGqcDtK7232QmYOKwNHS1yde5nQwJMfyUw3R8nUm9DBikkJV0ABuwug
2VhawI732GNBZ7bbJSA0sMXU0UJnfG37h7AIp7p8bp889Q4bHxHEndmwZNbwMjJw
FkPoJa5O/1K42DkdTEwoRN9p3aQB46FochzKibxbCpcttC0CQQDMqLgJt4wp6kH9
KU4q3vXQToGboxiEOzA82UXklVQsq0rKnHFNj4qB1jckvBrRPK63gNkWm4yKsgfe
EZh13QvbAkEAwPjbk8of8wGtWLghreVVDW39aBsFW/3F28+Ss90UBjobApQFxu+T
Zb68jfawLf2lqA44aa6q8qIXJaV25o/M3wJAeRYt5Rni5P3D0zw4EmdeOsPoLSRf
IgU+8hF/J9IuPkuOcbgD1WbjBRSwBZ0BpOBpYwrp5lVb3secngf9E2cYVwJBALll
XfBbXN6nWdfG7/SWRGSmq7N9YmTDJ3jLsHJFkJt678BGXlaGjeJOofDydMl6y9Dt
+JzwRyTdPcfZdKaGuZkCQGIxjl2LahWayKGj01cE5gV6GBtXqg4KGGEH/UIzC4vC
h874yfVODgDmYJmOKkv8/CixpYludkC5LHCYVD+WIi8=
-----END RSA PRIVATE KEY-----`

// TestJWTGenerator_GenerateToken
// Check the generation of JSON Web Token using a private rsa key to sign them
func TestJWTGenerator_GenerateToken(t *testing.T) {
	generator := JWTGenerator{}

	err := generator.SetPrivateKey([]byte(privateKey))
	if err != nil {
		t.Fatal(err)
	}

	tdt := []struct {
		input       interface{}
		output      model.Token
		expectedErr error
	}{
		{
			input: model.JWT{
				StandardClaims: model.StandardClaims{
					Id:       uuid.New().String(),
					Audience: "unit-tests",
					Issuer:   "go-test",
					Subject:  "Go",
					IssuedAt: time.Now().Unix(),
				},
				Scope: model.Mask{
					"read": 0xAAA,
				},
			},
			output: model.Token{
				Type: "Bearer",
				Scope: model.Mask{
					"read": 0xAAA,
				},
				AccessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1bml0LXRlc3RzIiwianRpIjoiMmQxYWQ2ZGYtMmNhOC00MWUzLWIwMGUtMDIxYmU5YzAzNWNjIiwiaWF0IjoxNjUwMTc0MTQ5LCJpc3MiOiJnby10ZXN0Iiwic3ViIjoiR28iLCJzY3AiOnsicmVhZCI6MjczMH19.b4R7S2luJkBJmziJd-QY1zdwY-HNv3V1gBEQV79o8FIpIDmNtaa7WXs21emeBxOXjv6QVWbyObOA_7VJ5e4QmQkS2zqC7Ez5vWfvcBRHs-46rrh5C5ZBaMD-GnGDtOuXVvWI4h10tNlr_FXObwoSTxgphSEr4quc269cfT4gpug",
			},
		},
	}

	for i, v := range tdt {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			token, err := generator.GenerateToken(v.input)
			if !errors.Is(err, v.expectedErr) {
				t.Fatalf(`unexpected error "%v"`, err)
			}

			if err != nil {
				t.Skip(err)
			}

			parser := &jwt.Parser{}
			gotToken, expectedToken := model.JWT{}, model.JWT{}

			parser.ParseUnverified(token.AccessToken, &gotToken)
			parser.ParseUnverified(v.output.AccessToken, &expectedToken)

			// Ignore JTI
			expectedToken.Id, gotToken.Id = "", ""
			expectedToken.IssuedAt, gotToken.IssuedAt = 0, 0

			if !reflect.DeepEqual(gotToken, expectedToken) {
				t.Fatalf(`
				mismatch expected token "%+v"
				from got token          "%+v"
			`, expectedToken, gotToken)
			}

			v.output.AccessToken = ""
			token.AccessToken = ""

			if !reflect.DeepEqual(token, v.output) {
				t.Fatalf(`mismatch token data "%+v" from "%+v"`, token, v.output)
			}

			t.Logf("%+v", token)
		})
	}
}
