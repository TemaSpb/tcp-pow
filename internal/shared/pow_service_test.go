package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateChallenge(t *testing.T) {
	// Generating a challenge should return the known challenge
	service := &ConcretePoWService{
		ChallengeLength: 1,
		Challenge:       "known_challenge_for_testing",
	}

	generatedChallenge := service.GenerateChallenge()
	assert.Equal(t, "known_challenge_for_testing", generatedChallenge, "Generated challenge does not match expected value")

	// Generating a challenge should return a challenge of specific length
	challengeLength := 3
	service = &ConcretePoWService{
		ChallengeLength: challengeLength,
	}

	generatedChallenge = service.GenerateChallenge()
	assert.Equal(t, challengeLength, len(generatedChallenge), "Generated challenge does not match expected length")
}

func TestValidateProof(t *testing.T) {
	service := &ConcretePoWService{
		Challenge:  "t",
		ZerosCount: 2,
	}
	challenge := "test"

	// Validate the valid proof
	validProof := "test2024-01-17T11:49:38Z"
	assert.True(t, service.ValidateProof(validProof, challenge), "Validation of a valid proof succeeded")

	// Validate the invalid proof
	invalidProof := "test2024-01-18T11:49:38Z"

	assert.False(t, service.ValidateProof(invalidProof, challenge), "Validation of an invalid proof failed")

	// Validate the valid proof (with a different prefix)
	invalidProof = "t2024-01-17T14:29:01Z"

	assert.False(t, service.ValidateProof(invalidProof, challenge), "Validation of an invalid proof failed")
}

// ToDo: Fix random seed for faster execution
func TestSolveChallenge(t *testing.T) {
	service := &ConcretePoWService{
		Challenge:  "test",
		ZerosCount: 1,
	}

	// Use the known challenge and prefix for testing
	challenge := service.GenerateChallenge()

	// Solve the challenge
	proof, err := service.SolveChallenge(challenge, -1)

	// Validate the generated proof
	assert.NoError(t, err)
	assert.True(t, service.ValidateProof(proof, challenge), "Validation of the generated proof failed")
}
