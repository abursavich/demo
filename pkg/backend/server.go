// SPDX-License-Identifier: MIT
//
// Copyright 2022 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

package backend

import (
	"context"

	bepb "bursavich.dev/grpcprom/testdata/backend"
)

type Server struct {
	bepb.UnimplementedBackendServer
}

func (*Server) Query(context.Context, *bepb.QueryRequest) (*bepb.QueryResponse, error) {
	return &bepb.QueryResponse{}, nil
}
