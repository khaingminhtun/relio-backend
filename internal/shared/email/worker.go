package email

import (
	"context"
	"os"
	"time"

	"github.com/khaingminhtun/relio-backend/internal/infrastructure/redis"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	queue  redis.EmailQueue
	sender Sender
}

func New(
	queue redis.EmailQueue,
	sender Sender,
) *Worker {
	return &Worker{
		queue:  queue,
		sender: sender,
	}
}

func (w *Worker) Start(ctx context.Context) {

	consumerName, err := os.Hostname()
	if err != nil {
		consumerName = "email-worker"
	}

	log.Info().
		Str("consumer", consumerName).
		Msg("email worker started")

	for {
		select {
		case <-ctx.Done():
			log.Info().
				Msg("email worker stopped")

			return

		default:
		}

		messages, err := w.queue.Consume(
			ctx,
			consumerName,
			10,
			5*time.Second,
		)

		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to consume email jobs")

			time.Sleep(time.Second)
			continue
		}

		for _, message := range messages {
			w.process(ctx, message)
		}
	}
}

func (w *Worker) process(
	ctx context.Context,
	message redis.EmailJobMessage,
) {
	job := message.Job

	log.Info().
		Str("message_id", message.MessageID).
		Str("to", job.To).
		Str("template", job.Template).
		Msg("processing email job")

	htmlBody, err := renderTemplate(job)
	if err != nil {
		log.Error().
			Err(err).
			Str("message_id", message.MessageID).
			Msg("failed to render email template")

		return
	}

	if err := w.sender.Send(
		ctx,
		job.To,
		job.Subject,
		htmlBody,
	); err != nil {

		log.Error().
			Err(err).
			Str("message_id", message.MessageID).
			Msg("failed to send email")

		return
	}

	if err := w.queue.Ack(
		ctx,
		message.MessageID,
	); err != nil {

		log.Error().
			Err(err).
			Str("message_id", message.MessageID).
			Msg("failed to acknowledge email job")

		return
	}

	log.Info().
		Str("message_id", message.MessageID).
		Str("to", job.To).
		Msg("email sent successfully")
}
