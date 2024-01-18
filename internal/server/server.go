package server

import (
	"bufio"
	"fmt"
	"golang.org/x/exp/slog"
	"net"
	"os"
	"strings"
	"tcp-pow/internal/config"
	"tcp-pow/internal/shared"
	"time"

	"github.com/Code-Hex/go-generics-cache"
)

// ClientRequestType represents the type of client request.
type ClientRequestType int

const (
	// RequestChallenge represents a client request for a new challenge.
	RequestChallenge ClientRequestType = iota
	// SendProof represents a client request to send the proof of work.
	SendProof
	// QuitRequest represents an quit client request type.
	QuitRequest
	// UnknownRequest represents an unknown client request type.
	UnknownRequest
)

var (
	errUnknownRequest = fmt.Errorf("unknown client request")

	quotes = []string{
		"Better by far you should forget and smile. Than that you should remember and be sad.",
		"Can I see another's woe, and not be in sorrow too? Can I see another's grief, and not seek for kind relief?",
		"Do not stand at my grave and cry; I am not there. I did not die.",
		"Don't hate the player. Hate the game.",
		"Even as the stone of the fruit must break, that its heart may stand in the sun, so must you know pain. And could you keep your heart in wonder at the daily miracles of your life, your pain would not seem less wondrous than your joy.",
		"I feel your pain the pain in knowing this has happened to you. The pain in knowing what more tears we have gained. But through all this I feel your pain",
		"If you try to please audiences, uncritically accepting their tastes, it can only mean that you have no respect for them",
		"I feel within me a peace above all earthly dignities, a still and quiet conscience.",
		"In the end… We only regret the chances we didn’t take",
		"LIFE is a mosaic of pleasure and pain - grief is an interval between two moments of joy. Peace is the interlude between two wars. You have no rose without a thorn; the diligent picker will avoid the pricks and gather the flower.",
		"Of all sad words of tongue or pen, the saddest are these, 'It might have been.",
	}
)

type Server struct {
	QuoteService     QuoteService
	PoWService       shared.PoWService
	cache            *cache.Cache[string, string]
	host             string
	port             int64
	hashcashDuration int64
	logger           *slog.Logger
}

func NewServer(cfg *config.Config) *Server {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	powService := shared.NewConcretePoWService(cfg.HashcashZeros, cfg.HashcashChallengeLength, cfg.HashcashChallenge)
	quoteService := NewQuoteService(quotes)

	return &Server{
		QuoteService:     quoteService,
		PoWService:       powService,
		cache:            cache.New[string, string](),
		host:             cfg.ServerHost,
		port:             cfg.ServerPort,
		hashcashDuration: cfg.HashcashDuration,
		logger:           l,
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("listen connection: %w", err)
	}

	defer ln.Close()
	s.logger.Info("Server listening", slog.String("addr", addr))

	for {
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("accept connection: %w", err)
		}
		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	s.logger.Info("Client connected", slog.String("addr", clientAddr))
	reader := bufio.NewReader(conn)

	for {
		requestType, err := s.readRequest(reader)
		if err != nil {
			s.logger.Error("Error reading client request: %v", slog.String("err", err.Error()))
			return
		}

		switch requestType {
		case RequestChallenge:
			if err := s.sendChallenge(conn, clientAddr); err != nil {
				s.logger.Error("Error sending challenge: %v", slog.String("err", err.Error()))
				return
			}

		case SendProof:
			if err := s.handleProof(conn, reader, clientAddr); err != nil {
				s.logger.Error("Error handling proof", slog.String("err", err.Error()))
				return
			}

		case QuitRequest:
			s.logger.Info("Quitting")
			return

		case UnknownRequest:
			s.logger.Error("Unknown client request: %v", slog.String("err", errUnknownRequest.Error()))
			return
		}
	}
}

func (s *Server) readRequest(reader *bufio.Reader) (ClientRequestType, error) {
	req, err := reader.ReadString('\n')
	if err != nil {
		s.logger.Error("Error read connection:", slog.String("err", err.Error()))

		return UnknownRequest, errUnknownRequest
	}

	switch strings.TrimSpace(req) {
	case "REQUEST_CHALLENGE":
		return RequestChallenge, nil
	case "SEND_PROOF":
		return SendProof, nil
	default:
		return UnknownRequest, errUnknownRequest
	}
}

func (s *Server) sendChallenge(conn net.Conn, clientAddr string) error {
	// Check if there's an existing challenge for the client
	if existingChallenge, exists := s.cache.Get(clientAddr); exists {
		if _, err := conn.Write([]byte(existingChallenge + "\n")); err != nil {
			return fmt.Errorf("error writing message: %v", err)
		}

		return nil
	}

	// Generate a new challenge and store it with an expiration time
	challenge := s.PoWService.GenerateChallenge()
	s.cache.Set(clientAddr, challenge, cache.WithExpiration(time.Second*time.Duration(s.hashcashDuration)))

	if _, err := conn.Write([]byte(challenge + "\n")); err != nil {
		return fmt.Errorf("error writing message: %v", err)
	}

	return nil
}

func (s *Server) handleProof(conn net.Conn, reader *bufio.Reader, clientAddr string) error {
	// Check if there's an existing challenge for the client
	existingChallenge, exists := s.cache.Get(clientAddr)
	if !exists {
		return fmt.Errorf("failed to found client's challenge")
	}

	proof, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read proof from client: %v", err)
	}

	proof = strings.TrimSpace(proof)

	if s.PoWService.ValidateProof(proof, existingChallenge) {
		quote := s.QuoteService.GetRandomQuote()

		if _, err := conn.Write([]byte(quote + "\n")); err != nil {
			return fmt.Errorf("error writing message: %v", err)
		}

		return nil
	}

	if _, err := conn.Write([]byte("Proof of Work validation failed. Connection closed.\n")); err != nil {
		return fmt.Errorf("error writing message: %v", err)
	}

	return nil
}
