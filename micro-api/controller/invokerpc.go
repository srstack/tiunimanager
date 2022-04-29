/******************************************************************************
 * Copyright (c)  2021 PingCAP, Inc.                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");            *
 * you may not use this file except in compliance with the License.           *
 * You may obtain a copy of the License at                                    *
 *                                                                            *
 * http://www.apache.org/licenses/LICENSE-2.0                                 *
 *                                                                            *
 * Unless required by applicable law or agreed to in writing, software        *
 * distributed under the License is distributed on an "AS IS" BASIS,          *
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.   *
 * See the License for the specific language governing permissions and        *
 * limitations under the License.                                             *
 ******************************************************************************/

package controller

import (
	"context"
	"github.com/pingcap-inc/tiem/common/errors"
	"github.com/pingcap-inc/tiem/proto/clusterservices"

	"encoding/json"

	"github.com/asim/go-micro/v3/client"
	"github.com/gin-gonic/gin"
	"github.com/pingcap-inc/tiem/library/framework"
)

// InvokeRpcMethod
// @Description: invoke cluster service from api
// @Parameter ctx context generated by gin framework, which contains traceId and operator info
// @Parameter rpcMethod the rpc method that defined in cluster service
// @Parameter requestBody thr request body in the rpc request
// @Parameter opts
func InvokeRpcMethod(
	ctx *gin.Context,
	rpcMethod func(ctx context.Context, in *clusterservices.RpcRequest, opts ...client.CallOption) (*clusterservices.RpcResponse, error),
	response interface{},
	requestBody string,
	opts ...client.CallOption) {

	rpcResponse, err := rpcMethod(framework.NewMicroCtxFromGinCtx(ctx),
		&clusterservices.RpcRequest{
			Request: requestBody,
		},
		opts...,
	)
	var withPage func() Page = nil
	if err == nil && rpcResponse.Page != nil {
		withPage = func() Page {
			return Page{int(rpcResponse.Page.Page),
				int(rpcResponse.Page.PageSize),
				int(rpcResponse.Page.Total)}
		}
	}
	HandleHttpResponse(ctx,
		err,
		func() (errors.EM_ERROR_CODE, string) {
			return errors.EM_ERROR_CODE(rpcResponse.GetCode()), rpcResponse.GetMessage()
		},
		func() (interface{}, error) {
			marshalError := json.Unmarshal([]byte(rpcResponse.Response), response)
			if marshalError != nil {
				return nil, errors.NewError(errors.TIEM_UNMARSHAL_ERROR, marshalError.Error())
			} else {
				return response, nil
			}
		},
		withPage,
	)
}