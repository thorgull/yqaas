/*
Package impl of YQ As A Service
Copyright (C) 2024 Thorgull

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/thorgull/yqaas/gen/api"
	"net/http"
	"os/exec"
)

type DefaultAPIService struct {
}

// NewDefaultAPIService creates a default api service
func NewDefaultAPIService() api.DefaultAPIServicer {
	return &DefaultAPIService{}
}

// EvaluatePost - Evaluate expression
func (s *DefaultAPIService) EvaluatePost(_ context.Context, evaluatePostRequest api.EvaluatePostRequest) (api.ImplResponse, error) {

	bs, err := json.Marshal(evaluatePostRequest.Data)

	if err != nil {
		return api.Response(http.StatusInternalServerError, nil), fmt.Errorf("can not serialize data %w", err)
	}

	var collector = NewOnlyCollectEncoder()
	_, err = yqlib.NewStringEvaluator().Evaluate(evaluatePostRequest.Expression, string(bs), collector, yqlib.NewJSONDecoder())
	if err != nil {
		return api.Response(http.StatusInternalServerError, nil), fmt.Errorf("can not evaluate expression %w", err)
	}

	return collector.Response(), nil

}

func (s *DefaultAPIService) JqEvaluatePost(ctx context.Context, evaluatePostRequest api.EvaluatePostRequest) (api.ImplResponse, error) {

	bs, err := json.Marshal(evaluatePostRequest.Data)

	if err != nil {
		return api.Response(http.StatusInternalServerError, nil), fmt.Errorf("can not serialize data %w", err)
	}

	cmd := exec.CommandContext(context.Background(), "jq", evaluatePostRequest.Expression)
	var output = bytes.NewBuffer(make([]byte, 0))
	cmd.Stdin = bytes.NewReader(bs)
	cmd.Stdout = output

	if cmd.Start() != nil {
		return api.Response(http.StatusInternalServerError, nil), fmt.Errorf("error while starting jq : %w", err)
	}

	if cmd.Wait() != nil {
		return api.Response(http.StatusInternalServerError, nil), fmt.Errorf("jq returned error : %w", err)
	}

	var result any
	if json.Unmarshal(output.Bytes(), &result) != nil {
		return api.Response(http.StatusInternalServerError, nil), fmt.Errorf("can not serialize jq output : %w", err)
	}

	return api.Response(200, result), nil

}
