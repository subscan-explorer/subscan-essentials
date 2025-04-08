/********************************************************************************
   This file is part of go-web3.
   go-web3 is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-web3 is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-web3.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

/**
 * @file http-provider.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package providers

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/itering/subscan/pkg/go-web3/providers/util"
)

type HTTPProvider struct {
	address string
	timeout int32
	secure  bool
	client  *http.Client
}

func NewHTTPProvider(address string, timeout int32, secure bool) *HTTPProvider {
	return NewHTTPProviderWithClient(address, timeout, secure, &http.Client{
		Timeout: time.Second * time.Duration(timeout),
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   time.Second * 10,
				KeepAlive: time.Second * 30,
			}).DialContext,
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 50,
			IdleConnTimeout:     time.Second * 90,
			TLSHandshakeTimeout: time.Second * 10,
			ForceAttemptHTTP2:   true,
		},
	})
}

func NewHTTPProviderWithClient(address string, timeout int32, secure bool, client *http.Client) *HTTPProvider {
	provider := new(HTTPProvider)
	provider.address = address
	provider.timeout = timeout
	provider.secure = secure
	provider.client = client

	return provider
}

func (provider HTTPProvider) SendRequest(ctx context.Context, v interface{}, method string, params interface{}) (err error) {
	bodyString := util.JSONRPCObject{Version: "2.0", Method: method, Params: params, ID: rand.Intn(1000000) + len(method)}
	body := strings.NewReader(bodyString.AsJsonString())
	var (
		req  *http.Request
		rsp  *http.Response
		data []byte
	)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, provider.address, body); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")
	if rsp, err = provider.client.Do(req); err != nil {
		return
	}
	defer rsp.Body.Close()
	if data, err = io.ReadAll(rsp.Body); err != nil {
		return
	}
	if rsp.StatusCode != 200 {
		return errors.New(rsp.Status)
	}
	if err = json.Unmarshal(data, v); err != nil {
		return
	}
	return
}

func (provider HTTPProvider) Close() error { return nil }
