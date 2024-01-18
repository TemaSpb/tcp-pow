// Package client - implements TCP-client.
package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"tcp-pow/internal/config"
	"tcp-pow/internal/shared"
	"time"

	"golang.org/x/exp/slog"
)

const (
	// MessageTypeQuit - quit message type.
	MessageTypeQuit = "QUIT"
	// MessageTypeRequestChallenge - request challenge message type.
	MessageTypeRequestChallenge = "REQUEST_CHALLENGE"
	// MessageTypeSendProof - sending proof message type.
	MessageTypeSendProof = "SEND_PROOF"
)

// Client - client for TCP server.
type Client struct {
	PoWService shared.PoWService

	host string
	port int64

	maxIterations int
	logger        *slog.Logger
}

// NewClient - constructor for Client.
func NewClient(cfg *config.Config) *Client {
	powService := shared.NewConcretePoWService(cfg.HashcashZeros, cfg.HashcashChallengeLength, cfg.HashcashChallenge)

	return &Client{
		PoWService: powService,
		host:       cfg.ServerHost,
		port:       cfg.ServerPort,
		logger:     slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (c *Client) Run() error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.logger.Info("Connected to", slog.String("address", addr))

	defer conn.Close()

	for {
		message, err := c.handleConnection(conn)
		if err != nil {
			return err
		}

		c.logger.Info("Quote result:", slog.String("quote", message))
		time.Sleep(5 * time.Second)
	}
}

func (c *Client) handleConnection(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	// Requesting challenge.
	_, err := conn.Write([]byte(MessageTypeRequestChallenge + "\n"))
	if err != nil {
		c.logger.Error("Write challenge request:", slog.String("err", err.Error()))

		return "", fmt.Errorf("send request: %w", err)
	}

	// Reading and parsing response.
	challenge, err := reader.ReadString('\n')
	if err != nil {
		c.logger.Error("Read challenge:", slog.String("err", err.Error()))

		return "", fmt.Errorf("read msg: %w", err)
	}

	solvedChallenge, err := c.PoWService.SolveChallenge(strings.TrimSpace(challenge), c.maxIterations)
	if err != nil {
		c.logger.Error("Compute hashcash:", slog.String("err", err.Error()))

		return "", fmt.Errorf("compute hashcash: %w", err)
	}

	_, err = conn.Write([]byte(MessageTypeSendProof + "\n"))
	if err != nil {
		c.logger.Error("Write sending proof request:", slog.String("err", err.Error()))

		return "", fmt.Errorf("send request: %w", err)
	}

	_, err = conn.Write([]byte(solvedChallenge + "\n"))
	if err != nil {
		c.logger.Error("Write solved challenge:", slog.String("err", err.Error()))

		return "", fmt.Errorf("send request: %w", err)
	}

	// Get result quote.
	quote, err := reader.ReadString('\n')
	if err != nil {
		c.logger.Error("Read msg:", slog.String("err", err.Error()))

		return "", fmt.Errorf("read msg: %w", err)
	}

	return strings.TrimSpace(quote), nil
}
