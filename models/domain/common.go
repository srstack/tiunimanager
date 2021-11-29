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
 *                                                                            *
 ******************************************************************************/

package domain

import (
	"time"

	"github.com/mozillazg/go-pinyin"
	"github.com/pingcap-inc/tiem/library/util/uuidutil"
	"gorm.io/gorm"
)

const (
	TABLE_NAME_CLUSTER            = "clusters"
	TABLE_NAME_DEMAND_RECORD      = "demand_records"
	TABLE_NAME_ACCOUNT            = "accounts"
	TABLE_NAME_TENANT             = "tenants"
	TABLE_NAME_ROLE               = "roles"
	TABLE_NAME_ROLE_BINDING       = "role_bindings"
	TABLE_NAME_PERMISSION         = "permissions"
	TABLE_NAME_PERMISSION_BINDING = "permission_bindings"
	TABLE_NAME_TOKEN              = "tokens"
	TABLE_NAME_TASK               = "tasks"
	TABLE_NAME_HOST               = "hosts"
	TABLE_NAME_DISK               = "disks"
	TABLE_NAME_USED_COMPUTE       = "used_computes"
	TABLE_NAME_USED_PORT          = "used_ports"
	TABLE_NAME_USED_DISK          = "used_disks"
	TABLE_NAME_LABEL              = "labels"
	TABLE_NAME_TIUP_CONFIG        = "tiup_configs"
	TABLE_NAME_TIUP_TASK          = "tiup_tasks"
	TABLE_NAME_FLOW               = "flows"
	TABLE_NAME_PARAMETERS_RECORD  = "parameters_records"
	TABLE_NAME_BACKUP_RECORD      = "backup_records"
	TABLE_NAME_BACKUP_STRATEGY    = "backup_strategies"
	TABLE_NAME_TRANSPORT_RECORD   = "transport_records"
	TABLE_NAME_RECOVER_RECORD     = "recover_records"
	TABLE_NAME_COMPONENT_INSTANCE = "component_instances"
	TABLE_NAME_CHANGE_FEED_TASKS  = "change_feed_tasks"
)

type Entity struct {
	ID        string    `gorm:"primaryKey;"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Code     string `gorm:"uniqueIndex;default:null;not null;<-:create"`
	TenantId string `gorm:"default:null;not null;<-:create"`
	Status   int8   `gorm:"type:SMALLINT;default:0"`
}

func (e *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuidutil.GenerateID()
	if e.Code == "" {
		e.Code = e.ID
	}
	e.Status = 0

	return nil
}

type Record struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	TenantId string `gorm:"default:null;not null;<-:create"`
}

type Data struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	BizId  string `gorm:"default:null;<-:create"`
	Status int8   `gorm:"type:SMALLINT;default:0"`
}

var split = []byte("_")

func generateEntityCode(name string) string {
	a := pinyin.NewArgs()

	a.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}

	code := pinyin.Pinyin(name, a)

	bytes := make([]byte, 0, len(name)*4)

	// split code with "_" before and after chinese word
	previousSplitFlag := false
	for _, v := range code {
		currentWord := v[0]

		if previousSplitFlag || len(currentWord) > 1 {
			bytes = append(bytes, split...)
		}

		bytes = append(bytes, []byte(currentWord)...)
		previousSplitFlag = len(currentWord) > 1
	}
	return string(bytes)
}
