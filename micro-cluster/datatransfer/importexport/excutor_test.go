/******************************************************************************
 * Copyright (c)  2021 PingCAP, Inc.                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");            *
 * you may not use this file except in compliance with the License.           *
 * You may obtain a copy of the License at                                    *
 *                                                                            *
 * http://www.apache.org/licenses/LICENSE-2.0                                 *
 *                                                                            *
 *  Unless required by applicable law or agreed to in writing, software       *
 *  distributed under the License is distributed on an "AS IS" BASIS,         *
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  *
 *  See the License for the specific language governing permissions and       *
 *  limitations under the License.                                            *
 ******************************************************************************/

package importexport

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pingcap-inc/tiem/library/secondparty"
	"github.com/pingcap-inc/tiem/micro-cluster/cluster/management/handler"
	"github.com/pingcap-inc/tiem/models"
	workflowModel "github.com/pingcap-inc/tiem/models/workflow"
	"github.com/pingcap-inc/tiem/test/mockmodels/mockimportexport"
	mock_secondparty_v2 "github.com/pingcap-inc/tiem/test/mocksecondparty_v2"
	"github.com/pingcap-inc/tiem/workflow"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExecutor_buildDataImportConfig(t *testing.T) {
	os.MkdirAll("./testdata", 0755)
	defer os.RemoveAll("./testdata")

	flowContext := workflow.NewFlowContext(context.TODO())
	flowContext.SetData(contextClusterMetaKey, &handler.ClusterMeta{})
	flowContext.SetData(contextDataTransportRecordKey, &ImportInfo{ConfigPath: "./testdata"})
	err := buildDataImportConfig(&workflowModel.WorkFlowNode{}, flowContext)
	assert.Nil(t, err)
}

func TestExecutor_importDataToCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTiupManager := mock_secondparty_v2.NewMockSecondPartyService(ctrl)
	mockTiupManager.EXPECT().Lightning(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
	secondparty.Manager = mockTiupManager

	flowContext := workflow.NewFlowContext(context.TODO())
	flowContext.SetData(contextDataTransportRecordKey, &ImportInfo{ConfigPath: "./testdata"})
	err := importDataToCluster(&workflowModel.WorkFlowNode{}, flowContext)
	assert.Nil(t, err)
}

func TestExecutor_updateDataImportRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImportExportRW := mockimportexport.NewMockReaderWriter(ctrl)
	mockImportExportRW.EXPECT().UpdateDataTransportRecord(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	models.SetImportExportReaderWriter(mockImportExportRW)

	flowContext := workflow.NewFlowContext(context.TODO())
	flowContext.SetData(contextDataTransportRecordKey, &ImportInfo{RecordId: "record-xxx"})
	err := updateDataImportRecord(&workflowModel.WorkFlowNode{}, flowContext)
	assert.Nil(t, err)
}

func TestExecutor_updateDataExportRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImportExportRW := mockimportexport.NewMockReaderWriter(ctrl)
	mockImportExportRW.EXPECT().UpdateDataTransportRecord(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	models.SetImportExportReaderWriter(mockImportExportRW)

	flowContext := workflow.NewFlowContext(context.TODO())
	flowContext.SetData(contextDataTransportRecordKey, &ExportInfo{RecordId: "record-xxx"})
	err := updateDataExportRecord(&workflowModel.WorkFlowNode{}, flowContext)
	assert.Nil(t, err)
}

func TestExecutor_importDataFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImportExportRW := mockimportexport.NewMockReaderWriter(ctrl)
	mockImportExportRW.EXPECT().UpdateDataTransportRecord(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	models.SetImportExportReaderWriter(mockImportExportRW)

	flowContext := workflow.NewFlowContext(context.TODO())
	flowContext.SetData(contextDataTransportRecordKey, &ImportInfo{RecordId: "record-xxx"})
	err := importDataFailed(&workflowModel.WorkFlowNode{}, flowContext)
	assert.Nil(t, err)
}

func TestExecutor_exportDataFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImportExportRW := mockimportexport.NewMockReaderWriter(ctrl)
	mockImportExportRW.EXPECT().UpdateDataTransportRecord(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	models.SetImportExportReaderWriter(mockImportExportRW)

	flowContext := workflow.NewFlowContext(context.TODO())
	flowContext.SetData(contextDataTransportRecordKey, &ExportInfo{RecordId: "record-xxx"})
	err := exportDataFailed(&workflowModel.WorkFlowNode{}, flowContext)
	assert.Nil(t, err)
}

func TestExecutor_exportDataFromCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTiupManager := mock_secondparty_v2.NewMockSecondPartyService(ctrl)
	mockTiupManager.EXPECT().Dumpling(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()
	secondparty.Manager = mockTiupManager

	os.MkdirAll("./testdata", 0755)
	defer os.RemoveAll("./testdata")

	flowContext := workflow.NewFlowContext(context.TODO())
	flowContext.SetData(contextDataTransportRecordKey, &ExportInfo{
		StorageType: "nfs",
		FileType:    "csv",
		FilePath:    "./testdata",
		Filter:      "*.db",
	})
	err := exportDataFromCluster(&workflowModel.WorkFlowNode{}, flowContext)
	assert.Nil(t, err)
}
