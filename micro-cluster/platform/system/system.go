/******************************************************************************
 * Copyright (c)  2022 PingCAP, Inc.                                          *
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

/*******************************************************************************
 * @File: system.go
 * @Description:
 * @Author: zhangpeijin@pingcap.com
 * @Version: 1.0.0
 * @Date: 2022/2/16
*******************************************************************************/

package system

import (
	"context"
	"github.com/pingcap-inc/tiem/common/constants"
	"github.com/pingcap-inc/tiem/common/errors"
	"github.com/pingcap-inc/tiem/library/framework"
	"github.com/pingcap-inc/tiem/models"
	"github.com/pingcap-inc/tiem/models/platform/system"
	"sync"
)

type SystemManager struct {
}

var manager *SystemManager
var once sync.Once

func GetSystemManager() *SystemManager {
	once.Do(func() {
		if manager == nil {
			manager = &SystemManager{}
		}
	})
	return manager
}

func (p *SystemManager) AcceptSystemEvent(ctx context.Context, event constants.SystemEvent) error {
	if len(event) == 0 {
		panic("unknown system event")
	}

	if statusMapAction, ok := actionBindings[event]; ok {
		systemInfo, err := p.GetSystemInfo(context.TODO())
		if err != nil {
			return err
		}
		if actionFunc, statusOK := statusMapAction[systemInfo.State]; statusOK {
			return actionFunc(ctx, event, systemInfo.State)
		}
	} else {
		panic("unknown system event")
	}
	return nil
}

func (p *SystemManager) GetSystemInfo(ctx context.Context) (*system.SystemInfo, error) {
	return models.GetSystemReaderWriter().GetSystemInfo(ctx)
}

func (p *SystemManager) GetVersionInfo(ctx context.Context, versionID string) (*system.VersionInfo, error) {
	var systemInfo *system.SystemInfo
	var versionInfo *system.VersionInfo
	return versionInfo, errors.OfNullable(nil).
		BreakIf(func() error {
			got, err := p.GetSystemInfo(ctx)
			systemInfo = got
			return err
		}).
		BreakIf(func() error {
			got, err := models.GetSystemReaderWriter().GetVersion(ctx, versionID)
			if err != nil {
				versionInfo = got
			}
			return err
		}).
		If(func(err error) {
			framework.LogWithContext(ctx).Errorf("get version info failed, versionID = %s, err = %s", versionID, err.Error())
		}).
		Else(func() {
			framework.LogWithContext(ctx).Infof("get version info succeed, versionID = %s,info = %v", versionID, *versionInfo)
		}).
		Present()
}
