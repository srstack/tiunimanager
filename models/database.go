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

package models

import (
	"context"

	"github.com/pingcap-inc/tiem/common/constants"
	"github.com/pingcap-inc/tiem/library/common"
	"github.com/pingcap-inc/tiem/library/framework"
	"github.com/pingcap-inc/tiem/models/cluster/backuprestore"
	"github.com/pingcap-inc/tiem/models/cluster/changefeed"
	"github.com/pingcap-inc/tiem/models/cluster/management"
	"github.com/pingcap-inc/tiem/models/cluster/parameter"
	"github.com/pingcap-inc/tiem/models/cluster/upgrade"
	"github.com/pingcap-inc/tiem/models/datatransfer/importexport"
	"github.com/pingcap-inc/tiem/models/parametergroup"
	"github.com/pingcap-inc/tiem/models/platform/config"
	"github.com/pingcap-inc/tiem/models/workflow"
	"github.com/pingcap-inc/tiem/models/workflow/secondparty"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

var defaultDb *database

type database struct {
	base                             *gorm.DB
	workFlowReaderWriter             workflow.ReaderWriter
	importExportReaderWriter         importexport.ReaderWriter
	brReaderWriter                   backuprestore.ReaderWriter
	changeFeedReaderWriter           changefeed.ReaderWriter
	upgradeReadWriter        upgrade.ReaderWriter
	clusterReaderWriter              management.ReaderWriter
	parameterGroupReaderWriter       parametergroup.ReaderWriter
	clusterParameterReaderWriter     parameter.ReaderWriter
	configReaderWriter               config.ReaderWriter
	secondPartyOperationReaderWriter secondparty.ReaderWriter
}

func Open(fw *framework.BaseFramework, reentry bool) error {
	dbFile := fw.GetDataDir() + common.DBDirPrefix + common.SqliteFileName
	logins := framework.LogForkFile(common.LogFileSystem)
	// todo tidb?
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})

	if err != nil || db.Error != nil {
		logins.Fatalf("open database failed, filepath: %s database error: %s, meta database error: %v", dbFile, err, db.Error)
		return err
	} else {
		logins.Infof("open database succeed, filepath: %s", dbFile)
	}

	defaultDb = &database{
		base: db,
	}

	defaultDb.initReaderWriters()
	if !reentry {
		defaultDb.initTables()
		defaultDb.initSystemData()
	}

	return nil
}

func (p *database) initTables() {
	p.addTable(new(changefeed.ChangeFeedTask))
	p.addTable(new(workflow.WorkFlow))
	p.addTable(new(workflow.WorkFlowNode))
	p.addTable(new(upgrade.ProductUpgradePath))
	p.addTable(new(management.Cluster))
	p.addTable(new(management.ClusterInstance))
	p.addTable(new(management.ClusterRelation))
	p.addTable(new(management.ClusterTopologySnapshot))
	p.addTable(new(importexport.DataTransportRecord))
	p.addTable(new(backuprestore.BackupRecord))
	p.addTable(new(backuprestore.BackupStrategy))
	p.addTable(new(config.SystemConfig))
	p.addTable(new(secondparty.SecondPartyOperation))
	p.addTable(new(parametergroup.Parameter))
	p.addTable(new(parametergroup.ParameterGroup))
	p.addTable(new(parametergroup.ParameterGroupMapping))
	p.addTable(new(parameter.ClusterParameterMapping))

	// other tables
}

func (p *database) initReaderWriters() {
	defaultDb.changeFeedReaderWriter = changefeed.NewGormChangeFeedReadWrite(defaultDb.base)
	defaultDb.workFlowReaderWriter = workflow.NewFlowReadWrite(defaultDb.base)
	defaultDb.importExportReaderWriter = importexport.NewImportExportReadWrite(defaultDb.base)
	defaultDb.brReaderWriter = backuprestore.NewBRReadWrite(defaultDb.base)
	defaultDb.upgradeReadWriter = upgrade.NewGormProductUpgradePath(defaultDb.base)
	defaultDb.parameterGroupReaderWriter = parametergroup.NewParameterGroupReadWrite(defaultDb.base)
	defaultDb.clusterParameterReaderWriter = parameter.NewClusterParameterReadWrite(defaultDb.base)
	defaultDb.configReaderWriter = config.NewConfigReadWrite(defaultDb.base)
	defaultDb.secondPartyOperationReaderWriter = secondparty.NewGormSecondPartyOperationReadWrite(defaultDb.base)
}

func (p *database) initSystemData() {
	defaultDb.configReaderWriter.CreateConfig(context.TODO(), &config.SystemConfig{ConfigKey: constants.ConfigKeyBackupStorageType, ConfigValue: string(constants.StorageTypeS3)})
	defaultDb.configReaderWriter.CreateConfig(context.TODO(), &config.SystemConfig{ConfigKey: constants.ConfigKeyBackupStoragePath, ConfigValue: constants.DefaultBackupStoragePath})
	defaultDb.configReaderWriter.CreateConfig(context.TODO(), &config.SystemConfig{ConfigKey: constants.ConfigKeyBackupS3AccessKey, ConfigValue: constants.DefaultBackupS3AccessKey})
	defaultDb.configReaderWriter.CreateConfig(context.TODO(), &config.SystemConfig{ConfigKey: constants.ConfigKeyBackupS3SecretAccessKey, ConfigValue: constants.DefaultBackupS3SecretAccessKey})
	defaultDb.configReaderWriter.CreateConfig(context.TODO(), &config.SystemConfig{ConfigKey: constants.ConfigKeyBackupS3Endpoint, ConfigValue: constants.DefaultBackupS3Endpoint})
	defaultDb.configReaderWriter.CreateConfig(context.TODO(), &config.SystemConfig{ConfigKey: constants.ConfigKeyExportShareStoragePath, ConfigValue: constants.DefaultExportPath})
	defaultDb.configReaderWriter.CreateConfig(context.TODO(), &config.SystemConfig{ConfigKey: constants.ConfigKeyImportShareStoragePath, ConfigValue: constants.DefaultImportPath})
}

func (p *database) addTable(gormModel interface{}) error {
	log := framework.LogForkFile(common.LogFileSystem)
	if !p.base.Migrator().HasTable(gormModel) {
		err := p.base.Migrator().CreateTable(gormModel)
		if err != nil {
			log.Errorf("create table failed, error : %v.", err)
			return err
		}
	}

	return nil
}

func GetChangeFeedReaderWriter() changefeed.ReaderWriter {
	return defaultDb.changeFeedReaderWriter
}

func SetChangeFeedReaderWriter(rw changefeed.ReaderWriter) {
	defaultDb.changeFeedReaderWriter = rw
}

func GetWorkFlowReaderWriter() workflow.ReaderWriter {
	return defaultDb.workFlowReaderWriter
}

func SetWorkFlowReaderWriter(rw workflow.ReaderWriter) {
	defaultDb.workFlowReaderWriter = rw
}

func GetImportExportReaderWriter() importexport.ReaderWriter {
	return defaultDb.importExportReaderWriter
}

func GetBRReaderWriter() backuprestore.ReaderWriter {
	return defaultDb.brReaderWriter
}

func GetUpgradeReaderWriter() upgrade.ReaderWriter {
	return defaultDb.upgradeReadWriter
}

func SetUpgradeReaderWriter(rw upgrade.ReaderWriter) {
	defaultDb.upgradeReadWriter = rw
}

func GetClusterReaderWriter() management.ReaderWriter {
	return defaultDb.clusterReaderWriter
}

func SetClusterReaderWriter(rw management.ReaderWriter) {
	defaultDb.clusterReaderWriter = rw
}

func GetConfigReaderWriter() config.ReaderWriter {
	return defaultDb.configReaderWriter
}

func SetConfigReaderWriter(rw config.ReaderWriter) {
	defaultDb.configReaderWriter = rw
}

func GetSecondPartyOperationReaderWriter() secondparty.ReaderWriter {
	return defaultDb.secondPartyOperationReaderWriter
}

func SetSecondPartyOperationReaderWriter(rw secondparty.ReaderWriter) {
	defaultDb.secondPartyOperationReaderWriter = rw
}

func GetParameterGroupReaderWriter() parametergroup.ReaderWriter {
	return defaultDb.parameterGroupReaderWriter
}

func SetParameterGroupReaderWriter(rw parametergroup.ReaderWriter) {
	defaultDb.parameterGroupReaderWriter = rw
}

func GetClusterParameterReaderWriter() parameter.ReaderWriter {
	return defaultDb.clusterParameterReaderWriter
}

func SetClusterParameterReaderWriter(rw parameter.ReaderWriter) {
	defaultDb.clusterParameterReaderWriter = rw
}

func MockDB() {
	defaultDb = &database{}
}
