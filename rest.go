package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func keyFromRequest(r *http.Request) string {
	key := strings.TrimLeft(r.URL.Path, "/")
	if cfg.RedisKeyPrefix != "" {
		key = cfg.RedisKeyPrefix + key
	}
	return key
}

func handlerDelete(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := client.Del(r.Context(), keyFromRequest(r)).Err(); err != nil {
			http.Error(w, fmt.Errorf("deleting key: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func handlerGet(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := client.Get(r.Context(), keyFromRequest(r)).Bytes()
		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
			if _, err = w.Write(content); err != nil {
				logrus.WithError(err).Debug("writing HTTP response")
			}

		case errors.Is(err, redis.Nil):
			w.WriteHeader(http.StatusNotFound)

		default:
			http.Error(w, fmt.Errorf("getting key: %w", err).Error(), http.StatusInternalServerError)
		}
	}
}

func handlerPut(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			data   = new(bytes.Buffer)
			err    error
			expire time.Duration
		)

		if rawEx := r.URL.Query().Get("expire"); rawEx != "" {
			if expire, err = time.ParseDuration(rawEx); err != nil {
				http.Error(w, fmt.Errorf("parsing expiry: %w", err).Error(), http.StatusBadRequest)
				return
			}
		}

		if _, err = io.Copy(data, r.Body); err != nil {
			http.Error(w, fmt.Errorf("reading payload: %w", err).Error(), http.StatusBadRequest)
			return
		}

		if err = client.Set(r.Context(), keyFromRequest(r), data.Bytes(), expire).Err(); err != nil {
			http.Error(w, fmt.Errorf("setting key: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
