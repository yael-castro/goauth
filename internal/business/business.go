package business

import (
	"github.com/google/uuid"
	"github.com/yael-castro/go-auth/internal/model"
	"math/rand"
	"strconv"
	"time"
)

// CodeGenerator generator of authorization codes
type CodeGenerator interface {
	// GenerateCode generates a random model.AuthorizationCode
	GenerateCode() model.AuthorizationCode
}

// CodeGeneratorFunc functional interface for CodeGenerator
type CodeGeneratorFunc func() model.AuthorizationCode

// GenerateCode executes f() to generate a random code
func (f CodeGeneratorFunc) GenerateCode() model.AuthorizationCode {
	return f()
}

// GenerateUUID creates a new model.AuthorizationCode that is basically an UUID
func GenerateUUID() model.AuthorizationCode {
	return model.AuthorizationCode(uuid.New().String())
}

// GenerateRandomCode closure function that generates a random number using as seed the unix time when the server starts
var GenerateRandomCode = generateRandomCode()

// generateRandomCode builds a CodeGeneratorFunc to create random numbers as model.AuthorizationCode
func generateRandomCode() CodeGeneratorFunc {
	random := rand.New(rand.NewSource(time.Now().Unix()))

	return func() model.AuthorizationCode {
		return model.AuthorizationCode(strconv.Itoa(random.Int()))
	}
}
