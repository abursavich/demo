// SPDX-License-Identifier: MIT
//
// Copyright 2022 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

package frontend

import (
	"context"

	bepb "bursavich.dev/grpcprom/testdata/backend"
	fepb "bursavich.dev/grpcprom/testdata/frontend"
)

type Server struct {
	BackendClient bepb.BackendClient
	fepb.UnimplementedFrontendServer
}

func (srv *Server) Query(ctx context.Context, req *fepb.QueryRequest) (*fepb.QueryResponse, error) {
	_, err := srv.BackendClient.Query(ctx, &bepb.QueryRequest{})
	if err != nil {
		return nil, err
	}
	return &fepb.QueryResponse{}, nil
}
