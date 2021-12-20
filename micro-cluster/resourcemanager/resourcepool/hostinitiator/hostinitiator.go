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

package hostinitiator

import (
	"context"

	"github.com/pingcap-inc/tiem/common/structs"
)

type HostInitiator interface {
	VerifyConnect(ctx context.Context, h *structs.HostInfo) (err error)
	CloseSSHConnect()
	VerifyCpuMem(ctx context.Context, h *structs.HostInfo) (err error)
	VerifyDisks(ctx context.Context, h *structs.HostInfo) (err error)
	VerifyFS(ctx context.Context, h *structs.HostInfo) (err error)
	VerifySwap(ctx context.Context, h *structs.HostInfo) (err error)
	VerifyEnv(ctx context.Context, h *structs.HostInfo) (err error)
	VerifyOSEnv(ctx context.Context, h *structs.HostInfo) (err error)
	SetOffSwap(ctx context.Context, h *structs.HostInfo) (err error)

	InstallSoftware(ctx context.Context, h *structs.HostInfo) (err error)
}
