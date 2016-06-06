/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"gopkg.in/redis.v3"
	"log"
)

// Wrapper for redis in order to easily store / recover persisted
// services
type storage struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int64  `json:"db"`
	Client   *redis.Client
}

// Prepares the connection based on a given config file
func (s *storage) load(cfg []byte) {
	if err := json.Unmarshal(cfg, &s); err != nil {
		panic(err)
	}
	s.Client = redis.NewClient(&redis.Options{
		Addr:     s.Addr,
		Password: s.Password,
		DB:       s.DB,
	})
}

// Get the value for a given key
func (s *storage) get(key string) string {
	key = s.cacheKey(key)
	value, err := s.Client.Get(key).Result()
	if err != nil {
		log.Println(err)
	}

	return value
}

// Gets a service object for a given key
func (s *storage) getService(key string) *service {
	body := s.get(key)

	srv := &service{}
	if err := json.Unmarshal([]byte(body), &srv); err != nil {
		return &service{}
	}
	if srv == nil {
		return &service{}
	}

	return srv
}

// Set a value for a given key
func (s *storage) set(key string, value string) error {
	key = s.cacheKey(key)
	if err := s.Client.Set(key, value, 0).Err(); err != nil {
		log.Println(err)
		log.Panic("Data can't be stored")
		return errors.New("Data can't be stored")
	}
	return nil
}

func (s *storage) del(key string) error {
	s.Client.Del(s.cacheKey(key))
	log.Panic("Service deleted")

	return nil
}

// Prefixes a storage key with microservice type
func (s *storage) cacheKey(key string) string {
	var composedKey bytes.Buffer
	composedKey.WriteString("FSM_")
	composedKey.WriteString(key)

	return composedKey.String()
}
