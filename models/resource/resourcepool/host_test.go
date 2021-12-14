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

package resourcepool

import (
	"os"
	"testing"

	"github.com/pingcap-inc/tiem/common/constants"
	"github.com/pingcap-inc/tiem/library/common"
	"github.com/pingcap-inc/tiem/library/framework"
	"github.com/pingcap-inc/tiem/library/util/uuidutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_BuildDefaultTraits(t *testing.T) {
	type want struct {
		errcode common.TIEM_ERROR_CODE
		Traits  int64
	}
	tests := []struct {
		name string
		host Host
		want want
	}{
		{"test1", Host{ClusterType: string(constants.EMProductNameTiDB), Purpose: string(constants.PurposeCompute), DiskType: string(constants.NVMeSSD)}, want{common.TIEM_SUCCESS, 73}},
		{"test2", Host{ClusterType: string(constants.EMProductNameTiDB), Purpose: string(constants.PurposeCompute) + "," + string(constants.PurposeSchedule), DiskType: string(constants.SSD)}, want{common.TIEM_SUCCESS, 169}},
		{"test3", Host{ClusterType: string(constants.EMProductNameTiDB), Purpose: "General", DiskType: string(constants.NVMeSSD)}, want{common.TIEM_RESOURCE_TRAIT_NOT_FOUND, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.host.BuildDefaultTraits()
			if err == nil {
				assert.Equal(t, tt.want.Traits, tt.host.Traits)
			} else {
				te, ok := err.(framework.TiEMError)
				assert.Equal(t, true, ok)
				assert.True(t, tt.want.errcode.Equal(int32(te.GetCode())))
			}
		})
	}
}

func TestHost_IsLoadLess(t *testing.T) {
	type want struct {
		IsLoadless bool
	}
	tests := []struct {
		name string
		host Host
		want want
	}{
		{"normal", Host{CpuCores: 4, Memory: 8, FreeCpuCores: 4, FreeMemory: 8, Disks: []Disk{{Status: string(constants.DiskAvailable)}}}, want{true}},
		{"diskused", Host{CpuCores: 16, Memory: 64, FreeCpuCores: 16, FreeMemory: 64, Disks: []Disk{{Status: string(constants.DiskInUsed)}}}, want{false}},
		{"diskused2", Host{CpuCores: 16, Memory: 64, FreeCpuCores: 16, FreeMemory: 64, Disks: []Disk{{Status: string(constants.DiskExhaust)}}}, want{false}},
		{"diskused3", Host{CpuCores: 16, Memory: 64, FreeCpuCores: 16, FreeMemory: 64, Disks: []Disk{{Status: string(constants.DiskExhaust)}, {Status: string(constants.DiskInUsed)}}}, want{false}},
		{"cpuused", Host{CpuCores: 16, Memory: 64, FreeCpuCores: 12, FreeMemory: 64, Disks: []Disk{{Status: string(constants.DiskAvailable)}}}, want{false}},
		{"memoryused", Host{CpuCores: 16, Memory: 64, FreeCpuCores: 16, FreeMemory: 8}, want{false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsLoadless := tt.host.IsLoadless(); IsLoadless != tt.want.IsLoadless {
				t.Errorf("IsLoadless() = %v, want %v", IsLoadless, tt.want)
			}
		})
	}
}

func TestHost_IsInused(t *testing.T) {
	type want struct {
		IsLoadless bool
	}
	tests := []struct {
		name string
		host Host
		want want
	}{
		{"normal", Host{CpuCores: 4, Memory: 8, FreeCpuCores: 4, FreeMemory: 8, Stat: string(constants.HostLoadLoadLess), Disks: []Disk{{Status: string(constants.DiskAvailable)}}}, want{false}},
		{"inused", Host{CpuCores: 16, Memory: 64, FreeCpuCores: 16, FreeMemory: 64, Stat: string(constants.HostLoadInUsed), Disks: []Disk{{Status: string(constants.DiskInUsed)}}}, want{true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsInused := tt.host.IsInused(); IsInused != tt.want.IsLoadless {
				t.Errorf("IsInused() = %v, want %v", IsInused, tt.want)
			}
		})
	}
}

func createDB(path string) (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

func genFakeHost(region, zone, rack, hostName, ip string, freeCpuCores, freeMemory int32, clusterType, purpose, diskType string) (h *Host) {
	host := Host{
		HostName:     hostName,
		IP:           ip,
		UserName:     "root",
		Passwd:       "4bc5947d63aab7ad23cda5ca33df952e9678d7920428",
		Status:       string(constants.HostOnline),
		Stat:         string(constants.HostLoadLoadLess),
		Arch:         "X86",
		OS:           "CentOS",
		Kernel:       "5.0.0",
		CpuCores:     freeCpuCores,
		Memory:       freeCpuCores,
		FreeCpuCores: freeCpuCores,
		FreeMemory:   freeMemory,
		Nic:          "1GE",
		Region:       region,
		AZ:           zone,
		Rack:         rack,
		ClusterType:  clusterType,
		Purpose:      purpose,
		DiskType:     diskType,
		Disks: []Disk{
			{Name: "sda", Path: "/", Capacity: 256, Status: string(constants.DiskReserved), Type: diskType},
			{Name: "sdb", Path: "/", Capacity: 256, Status: string(constants.DiskAvailable), Type: diskType},
		},
	}
	host.BuildDefaultTraits()
	return &host
}

func Test_Host_Hooks(t *testing.T) {
	dbPath := "./test_resource_" + uuidutil.ShortId() + ".db"
	db, err := createDB(dbPath)
	assert.Nil(t, err)
	defer func() { _ = os.Remove(dbPath) }()

	host := genFakeHost("Region1", "Region1,Zone1", "Region1,Zone1,Rack1", "TEST_HOST1", "192.168.999.999", 32, 64,
		string(constants.EMProductNameDataMigration), string(constants.PurposeSchedule), string(constants.NVMeSSD))
	db.AutoMigrate(&Host{})
	db.AutoMigrate(&Disk{})
	err = db.Model(&Host{}).Create(host).Error
	assert.Nil(t, err)
	hostId := host.ID
	assert.NotNil(t, hostId)
	for _, v := range host.Disks {
		assert.NotNil(t, v)
	}
	err = db.Delete(&Host{ID: host.ID}).Error
	assert.Nil(t, err)
	var queryHost Host
	err = db.Unscoped().Where("id = ?", host.ID).Find(&queryHost).Error
	assert.Nil(t, err)
	assert.Equal(t, string(constants.HostDeleted), queryHost.Status)
}
