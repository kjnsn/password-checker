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
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/rs/cors"
	"github.com/willf/bloom"
)

// RequestInput defines what the input looks like to the cloud function.
type RequestInput struct {
	Cleartext string `json:"cleartext"`
}

type Response struct {
	IsCommon bool `json:"is_common"`
}

var (
	filter          *bloom.BloomFilter
	c               = cors.AllowAll()
	initializedChan = make(chan struct{})

	bucket     = flag.String("bucket", "password-check-fn-km-dictionary", "The cloud bucket that contains the password dictionary")
	objectName = flag.String("object", "dictionary.txt", "Name of the dictionary object in the bucket")
)

func init() {
	n := uint(1000000)
	filter = bloom.NewWithEstimates(n, 0.00001)

	go func() {
		// Initialise the bloom filter.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := populateFilter(ctx); err != nil {
			// Cannot function without the filter.
			log.Fatalf("Attempted to populate filter, got error %q", err.Error())
		}

		close(initializedChan)
	}()
}

// CheckPassword returns 200 if the password is okay, or 400
// if the password has been found in the dictionary.
func CheckPassword(w http.ResponseWriter, r *http.Request) {
	// Only wait for 1 second for the system to initialize before failing fast.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Wait for the bloom filter to initialise before processing the request.
	select {
	case <-ctx.Done():
		http.Error(w, ctx.Err().Error(), http.StatusServiceUnavailable)
		return
	case <-initializedChan:
	}

	// Add CORS headers.
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

	isCommon := filter.TestString(input.Cleartext)
	enc := json.NewEncoder(w)
	if err := enc.Encode(&Response{
		IsCommon: isCommon,
	}); err != nil {
		http.Error(w, "Could not marshal response", http.StatusInternalServerError)
	}
}

// PopulateFilter reads a dictionary from a cloud storage bucket
// and populates the bloom filter with the newline delimited words
// in that dictionary.
func populateFilter(ctx context.Context) (err error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer func() { err = client.Close() }()
	rc, err := client.Bucket(*bucket).Object(*objectName).NewReader(ctx)
	if err != nil {
		return err
	}
	defer func() { err = rc.Close() }()

	// Read each word (delimited by newlines) from the object.
	scanner := bufio.NewScanner(rc)
	n := 0
	for scanner.Scan() {
		filter.Add(scanner.Bytes())
		n++
	}
	log.Printf("Added %d words to the filter", n)
	return scanner.Err()
}
