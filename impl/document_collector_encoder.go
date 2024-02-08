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
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/thorgull/yqaas/gen/api"
	"io"
	"net/http"
)

type OnlyCollectEncoder struct {
	Documents []*yqlib.CandidateNode
}

func (o *OnlyCollectEncoder) Encode(_ io.Writer, node *yqlib.CandidateNode) error {
	o.Documents = append(o.Documents, node)
	return nil
}

func (o *OnlyCollectEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (o *OnlyCollectEncoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (o *OnlyCollectEncoder) CanHandleAliases() bool {
	return false
}

func (o *OnlyCollectEncoder) Response() api.ImplResponse {
	switch len(o.Documents) {
	case 0:
		return api.Response(http.StatusNotFound, nil)
	case 1:
		return api.Response(http.StatusOK, o.Documents[0])
	default:
		return api.Response(http.StatusOK, o.Documents)
	}
}

func NewOnlyCollectEncoder() *OnlyCollectEncoder {
	return &OnlyCollectEncoder{Documents: make([]*yqlib.CandidateNode, 0)}
}
