package ubersmith

import (
  "net/url"
  "fmt"
)

type Client struct {
  URL       *url.URL
  Username  string
  Password  string
}

func NewClient(server, username, password string) (*Client, error) {
  u, err := url.Parse(fmt.Sprintf("%s/api/2.0/", server))
  if err != nil {
    return nil, err
  }

  return &Client {
    URL: u,
    Username: username,
    Password: password,
  }, nil
}
