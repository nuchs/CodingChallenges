package main

import (
	"math/rand"
)

type Balancer struct {
	urls []string
}

func NewBalancer(urls []string) *Balancer {
	return &Balancer{
		urls: urls,
	}
}

func (b *Balancer) selectService() (string, error) {
	service := rand.Intn(len(b.urls))
	return b.urls[service], nil
}
