// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package passwordchecker

import (
	"encoding/json"
	"net/http"
	"context"
	"time"

	"github.com/willf/bloom"
	"github.com/rs/cors"
)

// RequestInput defines what the input looks like to the cloud function.
type RequestInput struct {
	Cleartext string `json:"cleartext"`
}

var (
	filter *bloom.BloomFilter
	c = cors.AllowAll()
	initializedChan = make(chan struct{})
)

func init() {
	n := uint(1000)
	filter = bloom.New(20*n, 5)

	go func() {
		// Initialise the bloom filter here.
		close(initializedChan)
	}()
}

// CheckPassword returns 200 if the password is okay, or 400
// if the password has been found in the dictionary.
func CheckPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Wait for the bloom filter to initialise before processing the request.
	select {
	case <-ctx.Done():
		http.Error(w, ctx.Err().Error(), 500)
		return
	case <-initializedChan:
	}

	c.HandlerFunc(w, r)

	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	input := new(RequestInput)
	if err := dec.Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.Cleartext == "" {
		http.Error(w, "Empty cleartext", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getPasswordDictionary(ctx context.Context) ([]byte, error) {
	return []byte{}, nil
}
