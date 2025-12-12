package main_test

import (
	"errors"
	"reflect"
	"testing"

	lb "github.com/nuchs/cc/lb/cmd/lb"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		want lb.Config
	}{
		{
			desc: "Set port",
			args: []string{"-p", "22", "http://localhost"},
			want: lb.Config{Port: 22, Urls: []string{"http://localhost"}},
		},
		{
			desc: "Set urls",
			args: []string{"-p", "3", "http://a:1", "http://b:2"},
			want: lb.Config{Port: 3, Urls: []string{"http://a:1", "http://b:2"}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := lb.NewConfig(tC.args)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Bad config - got: %+v, want: %+v", got, tC.want)
			}
		})
	}
}

func TestBadConfig(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		want error
	}{
		{
			desc: "No urls",
			args: []string{"-p", "1"},
			want: lb.ErrNoUrls,
		},
		{
			desc: "No port",
			args: []string{"http:a:8"},
			want: lb.ErrInvalidPort,
		},
		{
			desc: "missing protocol",
			args: []string{"-p", "1", "bacon"},
			want: lb.ErrBadProtocol,
		},
		{
			desc: "Unsupported protocol",
			args: []string{"-p", "1", "ftp://:8080"},
			want: lb.ErrBadProtocol,
		},
		{
			desc: "no host",
			args: []string{"-p", "1", "http://:80"},
			want: lb.ErrMissingHost,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := lb.NewConfig(tC.args)
			if err == nil {
				t.Fatalf("Expected error, got nil")
			}
			if !errors.Is(err, tC.want) {
				t.Fatalf("Wrong error - got %s, want %s", err, tC.want)
			}
		})
	}
}
