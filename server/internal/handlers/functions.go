package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func wrap(c Context, fn func(c Context, r *http.Request) (any, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if response, err := fn(c, r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			var b []byte
			if b, err = json.Marshal(response); err != nil {
				b = []byte(fmt.Sprintf("%v", response))
			}
			if _, err = w.Write(b); err != nil {
				slog.Error(err.Error())
			}

		}
	}
}

func text(str string) (any, error) {
	return struct {
		Text string `json:"text"`
	}{Text: str}, nil
}

func paramInt(r *http.Request, name string) (int, error) {
	queryVal := r.URL.Query().Get(name)
	if queryVal == "" {
		return 0, nil
	}
	if val, err := strconv.Atoi(queryVal); err != nil {
		return 0, err
	} else {
		return val, nil
	}
}

func paramIntSlice(r *http.Request, name string) ([]int, error) {
	queryVal := r.URL.Query().Get(name)
	if queryVal == "" {
		return nil, nil
	}
	strValues := strings.Split(queryVal, ",")
	result := make([]int, len(strValues))
	for i, v := range strValues {
		if val, err := strconv.Atoi(v); err != nil {
			return nil, err
		} else {
			result[i] = val
		}
	}

	return result, nil
}

func headerAppID(r *http.Request) (int, error) {
	valStr := r.Header.Get("X-App-Id")
	if valStr == "" {
		return 0, errors.New("missing X-App-Id")
	}
	id, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, errors.New("X-App-Id must be integer")
	}
	if id == 0 {
		return 0, errors.New("X-App-Id cannot be 0")
	}

	return id, nil
}
