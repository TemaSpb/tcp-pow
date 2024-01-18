package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"
	"strings"
	"time"
)

const (
	timestampFormat = "2006-01-02T15:04:05Z07:00"
	charset         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	errMaxIterationsExceeded = errors.New("max iterations exceeded")
)

// PoWService provides Proof of Work-related functionality for the Word of Wisdom server.
type PoWService interface {
	GenerateChallenge() string
	ValidateProof(proof, challenge string) bool
	SolveChallenge(challenge string, maxIterations int) (string, error)
}

// ConcretePoWService is the concrete implementation of PoWService.
// Challenge field may be used to specify a constant challenge for validation/testing.
type ConcretePoWService struct {
	ZerosCount      int
	ChallengeLength int
	Challenge       string
}

// NewConcretePoWService creates a new instance of ConcretePoWService.
func NewConcretePoWService(zerosCount, challengeLength int, challenge string) *ConcretePoWService {
	cl := challengeLength
	if challenge != "" {
		cl = len(challenge)
	}

	return &ConcretePoWService{
		ZerosCount:      zerosCount,
		ChallengeLength: cl,
		Challenge:       challenge,
	}
}

// GenerateChallenge generates a random challenge for Proof of Work.
func (ps *ConcretePoWService) GenerateChallenge() string {
	if ps.Challenge != "" {
		return ps.Challenge
	}

	return randomString(ps.ChallengeLength)
}

// ValidateProof checks if the provided proof is valid for the given challenge.
func (ps *ConcretePoWService) ValidateProof(proof, challenge string) bool {
	if !strings.HasPrefix(proof, challenge) {
		return false
	}

	hash := sha256.Sum256([]byte(proof))
	hashHex := hex.EncodeToString(hash[:])
	prefix := strings.Repeat("0", ps.ZerosCount)
	return strings.HasPrefix(hashHex, prefix)
}

// SolveChallenge generates a proof for a given challenge using the Hashcash algorithm.
func (ps *ConcretePoWService) SolveChallenge(challenge string, maxIterations int) (string, error) {
	return findHashcashProof(challenge, ps.ZerosCount, maxIterations)
}

// findHashcashProof implements the Hashcash algorithm to find a valid proof for the challenge.
func findHashcashProof(challenge string, zerosCount, maxIterations int) (string, error) {
	for i := 1; i <= maxIterations || maxIterations <= 0; i++ {
		timestamp := time.Now().UTC().Format(timestampFormat)
		proof := challenge + timestamp
		hash := sha256.Sum256([]byte(proof))
		hashHex := hex.EncodeToString(hash[:])
		prefix := strings.Repeat("0", zerosCount)
		if strings.HasPrefix(hashHex, prefix) {
			return proof, nil
		}
	}

	return "", errMaxIterationsExceeded
}

// randomString generates a random string of the given length.
func randomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
