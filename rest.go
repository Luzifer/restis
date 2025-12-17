package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
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
			http.Error(w, errors.Wrap(err, "deleting key").Error(), http.StatusInternalServerError)
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
			w.Write(content)

		case errors.Is(err, redis.Nil):
			w.WriteHeader(http.StatusNotFound)

		default:
			http.Error(w, errors.Wrap(err, "getting key").Error(), http.StatusInternalServerError)
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
				http.Error(w, errors.Wrap(err, "parsing expiry").Error(), http.StatusBadRequest)
				return
			}
		}

		if _, err = io.Copy(data, r.Body); err != nil {
			http.Error(w, errors.Wrap(err, "reading payload").Error(), http.StatusBadRequest)
			return
		}

		if err = client.Set(r.Context(), keyFromRequest(r), data.Bytes(), expire).Err(); err != nil {
			http.Error(w, errors.Wrap(err, "setting key").Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
