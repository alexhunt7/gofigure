package master

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

// RunAll is a helper function for running a function on many clients in parallel.
func RunAll(clients map[string]*Client, f func(host string, client *Client) error) error {
	type result struct {
		host string
		err  error
	}

	results := make(chan *result, len(clients))

	wrappedF := func(host string, client *Client) {
		err := f(host, client)
		results <- &result{host: host, err: err}
	}

	for host, client := range clients {
		go wrappedF(host, client)
	}

	var merr error
	for range clients {
		res := <-results
		if res.err != nil {
			merr = multierror.Append(merr, fmt.Errorf("%s: %v", res.host, res.err))
		}
	}

	return merr
}
