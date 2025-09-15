package processor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Koshsky/polyschedule-backend/internal/kafka"
	"github.com/rs/zerolog/log"
)

type Producer interface {
	Send(ctx context.Context, key, value []byte) error
}

type Processor struct {
	producer Producer
}

func NewProcessor(producer Producer) *Processor {
	return &Processor{producer: producer}
}

// Request and Response — simple stub message types
type Request struct {
	ChatID int64  `json:"chat_id"`
	Query  string `json:"query"`
}

type Response struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
	TS     int64  `json:"ts"`
}

func (p *Processor) Process(ctx context.Context, key, value []byte) error {
	var req Request
	if err := json.Unmarshal(value, &req); err != nil {
		log.Error().Err(err).Msg("invalid request payload")
		return nil // do not fail the loop, just log
	}

	// Stub: build a simple response
	resp := Response{
		ChatID: req.ChatID,
		Text:   "Schedule (stub) for query: " + req.Query,
		TS:     time.Now().Unix(),
	}
	data, _ := json.Marshal(resp)

	if err := p.producer.Send(ctx, key, data); err != nil {
		log.Error().Err(err).Msg("failed to produce response")
		return err
	}
	log.Info().Int64("chat_id", req.ChatID).Msg("response produced")
	return nil
}

var _ kafka.MessageHandler = (*Processor)(nil).Process
