package ubersmith

import (
  "mime/multipart"
  "path/filepath"
  "encoding/json"
  "io/ioutil"
  "net/http"
  "net/url"
  "bytes"
  "fmt"
  "os"
  "io"
)

type Response struct {
  Status        bool    `json:"status"`
  ErrorCode     int     `json:"error_code"`
  ErrorMessage  string  `json:"error_message"`
  Data          json.RawMessage
}

func (u *Client) doRequest(req *http.Request) ([]byte, error) {
  req.SetBasicAuth(u.Username, u.Password)
  req.Header.Set("User-Agent", "Ubersmith API Client Go/1.0")
  client := &http.Client{}

  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  if 200 != resp.StatusCode {
    return nil, fmt.Errorf("%s", body)
  }
  return body, nil
}

func (c *Client) Request(method string, params url.Values) (*Response, error) { //params map[string][]string
  params.Set("method", method)
  c.URL.RawQuery = params.Encode()

  req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

  var res Response
  err = json.Unmarshal(bytes, &res)
  if err != nil {
		return nil, err
	}

  if !res.Status {
    return nil, fmt.Errorf("Error %d, %s", res.ErrorCode, res.ErrorMessage)
  }

  return &res, nil
}

func (c *Client) PostRequest(method string, params url.Values) (*Response, error) { //params map[string][]string
  m := url.Values{}
  m.Add("method", method)
  c.URL.RawQuery = m.Encode()

  body := bytes.NewBufferString(params.Encode())

  req, err := http.NewRequest("POST", c.URL.String(), body)
	if err != nil {
		return nil, err
	}

  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

  var res Response
  err = json.Unmarshal(bytes, &res)
  if err != nil {
		return nil, err
	}

  if !res.Status {
    return nil, fmt.Errorf("Error %d, %s", res.ErrorCode, res.ErrorMessage)
  }

  return &res, nil
}

func (c *Client) Upload(method string, params url.Values, uploadParam string, uploadPath string) (*Response, error) { //params map[string][]string
  file, err := os.Open(uploadPath)
  if err != nil {
      return nil, err
  }
  defer file.Close()

  body := &bytes.Buffer{}
  writer := multipart.NewWriter(body)
  part, err := writer.CreateFormFile(uploadParam, filepath.Base(uploadPath))
  if err != nil {
      return nil, err
  }

  _, err = io.Copy(part, file)
  if err != nil {
      return nil, err
  }

  err = writer.Close()
  if err != nil {
      return nil, err
  }

  params.Set("method", method)
  c.URL.RawQuery = params.Encode()

  req, err := http.NewRequest("POST", c.URL.String(), body)
	if err != nil {
		return nil, err
	}

  req.Header.Set("Content-Type", writer.FormDataContentType())

	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

  var res Response
  err = json.Unmarshal(bytes, &res)
  if err != nil {
		return nil, err
	}

  if !res.Status {
    return nil, fmt.Errorf("Error %d, %s", res.ErrorCode, res.ErrorMessage)
  }

  return &res, nil
}
