package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Handler struct {
	topics []string
	group  sarama.ConsumerGroup
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	handle Handle
	logger *zap.Logger
}

type Handle func(msg *sarama.ConsumerMessage)

func NewHandler(ctx context.Context, group sarama.ConsumerGroup, topics []string) (
	handler *Handler) {

	handler = &Handler{
		topics: topics,
		group:  group,
		wg:     new(sync.WaitGroup),
	}

	handler.ctx, handler.cancel = context.WithCancel(ctx)

	return handler
}

func HandlerFromConfig(ctx context.Context, vp *viper.Viper, field string) (
	handler *Handler, err error) {
	var (
		config *Config
		scfg   *sarama.Config
		group  sarama.ConsumerGroup
	)

	if config, scfg, err = NewConfigFromViper(vp, field); err != nil {
		return nil, err
	}

	if group, err = sarama.NewConsumerGroup(config.Addrs, config.GroupId, scfg); err != nil {
		return nil, err
	}

	handler = NewHandler(ctx, group, []string{config.Topic})

	// alter handler.Logger later
	return handler, nil
}

func (handler *Handler) WithHandle(handle Handle) *Handler {
	handler.handle = handle
	return handler
}

func (handler *Handler) WithLogger(logger *zap.Logger) *Handler {
	handler.logger = logger
	return handler
}

func (handler *Handler) Ok() (err error) {
	if handler.handle == nil {
		return fmt.Errorf("hanlder is unset")
	}

	if handler.logger == nil {
		return fmt.Errorf("logger is unset")
	}

	return nil
}

func (handler *Handler) Consume() {
	go handler.consume()
}

// maxRetries: 5, wait: 3 * time.Second
func (handler *Handler) consume() {
	var (
		err        error
		n          int
		maxRetries int
		wait       time.Duration
	)
	n, maxRetries = 0, 5
	wait = 3 * time.Second

	handler.logger.Info("Consume start")
	for {
		err = handler.group.Consume(handler.ctx, handler.topics, handler)
		if err != nil {
			if errors.Is(sarama.ErrClosedConsumerGroup, err) {
				handler.logger.Warn(fmt.Sprintf("Consume closed: %v", err))
				return
			}
			handler.logger.Error(fmt.Sprintf("Consume error: %v", err))
		} else {
			handler.logger.Info("Consume end") // occurs when reset offset
		}

		if n++; n > maxRetries {
			handler.logger.Warn(
				fmt.Sprintf("Consume exceeds max retries limit: %d", maxRetries),
			)
			return
		}

		// if err = handler.ctx.Err(); err != nil {
		// 	handler.logger.Error("!!! Handler.Consume ctx.Err(): %v", err)
		// 	return
		// }
		select {
		case <-handler.ctx.Done():
			return
		default:
		}
		time.Sleep(wait)
	}
}

func (handler *Handler) Close() error {
	handler.cancel()
	handler.wg.Wait()

	return handler.group.Close()
}

func (handler *Handler) Setup(sess sarama.ConsumerGroupSession) (err error) {
	handler.logger.Info("Setup run")

	go func() {
		var err error

		for {
			select {
			case err = <-handler.group.Errors():
				if errors.Is(sarama.ErrClosedConsumerGroup, err) {
					handler.logger.Warn("Setup closed")
					return
				}
				handler.logger.Error(fmt.Sprintf("Setup error: %v", err))
			case <-handler.ctx.Done():
				handler.logger.Info("Setup done")
				return
			}
		}
	}()

	return nil
}

func (handler *Handler) Cleanup(sess sarama.ConsumerGroupSession) (err error) {
	// TODO
	handler.logger.Info("Cleanup run")
	return nil
}

func (handler *Handler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {

	handler.wg.Add(1)
	defer handler.wg.Done()

	for {
		select {
		case msg := <-claim.Messages():
			if msg == nil {
				return nil
			}

			handler.handle(msg)
			// sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset, "some-metadata")
			sess.MarkMessage(msg, "consumed")
		case <-handler.ctx.Done():
			handler.logger.Warn("ConsumeClaim canceled")
			return nil
		}
	}
}
