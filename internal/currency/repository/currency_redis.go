package repository

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sefikcan/kanbersky.ca/internal/dto/response/currency"
	"time"
)

type CurrencyRedisRepository interface {
	GetByKey(ctx context.Context, key string) (*currency.CurrencyResponse, error)
	Set(ctx context.Context, key string, seconds int, param any) error
	Delete(ctx context.Context, key string) error
}

type currencyRedisRepository struct {
	redisClient *redis.Client
}

func (c currencyRedisRepository) GetByKey(ctx context.Context, key string) (*currency.CurrencyResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyRedisRepository.GetByKey")
	defer span.Finish()

	currencyByte, err := c.redisClient.Get(spanContext, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err,"currencyRedisRepository.GetByKey.RedisClient.Get")
	}

	currency := &currency.CurrencyResponse{}
	if err = json.Unmarshal(currencyByte, currency); err != nil {
		return nil, errors.Wrap(err, "currencyRedisRepository.GetByKey.Json.Unmarshal")
	}

	return currency, nil
}

func (c currencyRedisRepository) Set(ctx context.Context, key string, seconds int, param any) error {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyRedisRepository.Set")
	defer span.Finish()

	currencyByte, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, "currencyRedisRepository.Set.Json.Marshal")
	}
	if err = c.redisClient.Set(spanContext, key, currencyByte, time.Second * time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "currencyRedisRepository.Set.RedisClient.Set")
	}

	return nil
}

func (c currencyRedisRepository) Delete(ctx context.Context, key string) error {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyRedisRepository.Delete")
	defer span.Finish()

	if err := c.redisClient.Del(spanContext, key).Err(); err != nil {
		return errors.Wrap(err, "currencyRedisRepository.Delete.RedisClient.Del")
	}

	return nil
}

func NewCurrencyRedisRepository(redisClient *redis.Client) CurrencyRedisRepository {
	return &currencyRedisRepository{
		redisClient: redisClient,
	}
}
