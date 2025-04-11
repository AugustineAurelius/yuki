package yuki_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/AugustineAurelius/yuki/converter"
)

type testPut struct {
	Key   string
	Value any
}
type get struct {
	Value json.RawMessage `json:"value"`
}

func YukiTest(host string, port int) error {
	url := "http://" + host + ":" + strconv.Itoa(port)

	slog.Info("created url", "url", url)
	p := testPut{
		Key:   "123",
		Value: 123,
	}
	data, err := json.Marshal(&p)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("bad status code")
	}

	slog.Info("post successfull")

	resp, err = http.Get(url + "/123")
	if err != nil {
		return err
	}

	var g get
	if err = json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return err
	}

	slog.Info("get successfull", "value", string(g.Value))
	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int, b []byte) string {
	if n == 0 {
		n = 10
	}
	for i := range b {
		b[i] = letterBytes[rand.IntN(len(letterBytes))]
	}
	return converter.BytesToString(b)
}

func YukiInfTest(host string, port int) error {
	url := "http://" + host + ":" + strconv.Itoa(port)

	slog.Info("created url", "url", url)
	key := make([]byte, 100)
	val := make([]byte, 100)

	for {
		p := testPut{
			Key:   RandStringBytes(rand.IntN(100), key),
			Value: RandStringBytes(rand.IntN(100), val),
		}
		data, err := json.Marshal(&p)
		if err != nil {
			return err
		}
		resp, err := http.Post(url, "application/json", bytes.NewReader(data))
		if err != nil {
			return err
		}

		if resp.StatusCode != 201 {
			return fmt.Errorf("bad status code")
		}

		slog.Info("post successfull")

		resp, err = http.Get(url + "/" + converter.BytesToString(key))
		if err != nil {
			return err
		}

		var g get
		if err = json.NewDecoder(resp.Body).Decode(&g); err != nil {
			return err
		}

		slog.Info("get successfull", "value", string(g.Value))
	}

	return nil
}
