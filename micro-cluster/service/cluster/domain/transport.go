package domain

import (
	ctx "context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pingcap-inc/tiem/library/client"
	"github.com/pingcap-inc/tiem/library/secondparty/libtiup"
	proto "github.com/pingcap-inc/tiem/micro-cluster/proto"
	db "github.com/pingcap-inc/tiem/micro-metadb/proto"
	"os"
	"strconv"
	"time"
)

type TransportType string

const (
	TransportTypeExport TransportType = "export"
	TransportTypeImport TransportType = "import"
)

const (
	TransportStatusRunning string = "Running"
	TransportStatusFailed  string = "Failed"
	TransportStatusSuccess string = "Success"
)

type ImportInfo struct {
	ClusterId   string
	UserName    string
	Password    string
	FilePath    string
	RecordId    string
	StorageType string
	ConfigPath  string
}

type ExportInfo struct {
	ClusterId   string
	UserName    string
	Password    string
	FileType    string
	RecordId    string
	FilePath    string
	Filter      string
	Sql 		string
	StorageType string
	BucketRegion string
}

type TransportInfo struct {
	ClusterId     string
	RecordId      string
	TransportType string
	Status        string
	FilePath      string
	StartTime     int64
	EndTime       int64
}

/*
	data import toml config for lightning
	https://docs.pingcap.com/zh/tidb/dev/tidb-lightning-configuration
*/
type DataImportConfig struct {
	Lightning    LightningCfg    `toml:"lightning"`
	TikvImporter TikvImporterCfg `toml:"tikv-importer"`
	MyDumper     MyDumperCfg     `toml:"mydumper"`
	Tidb         TidbCfg         `toml:"tidb"`
}

type LightningCfg struct {
	Level             string `toml:"level"`              //lightning log level
	File              string `toml:"file"`               //lightning log path
	CheckRequirements bool   `toml:"check-requirements"` //lightning pre check
}

/*
	tidb-lightning backend
	https://docs.pingcap.com/zh/tidb/stable/tidb-lightning-backends#tidb-lightning-backend
*/
const (
	BackendLocal  string = "local"
	BackendImport string = "importer"
	BackendTidb   string = "tidb"
)

const (
	NfsStorageType string = "nfs"
	S3StorageType  string = "s3"
)

const (
	FileTypeCSV string = "csv"
	FileTypeSQL string = "sql"
)

const (
	DefaultTidbPort       int = 4000
	DefaultTidbStatusPort int = 10080
	DefaultPDClientPort   int = 2379
)

type TikvImporterCfg struct {
	Backend     string `toml:"backend"`       //backend mode: local/normal
	SortedKvDir string `toml:"sorted-kv-dir"` //temp store path
}

type MyDumperCfg struct {
	DataSourceDir string `toml:"data-source-dir"` //import data filepath
}

type TidbCfg struct {
	Host       string `toml:"host"`
	Port       int    `toml:"port"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	StatusPort int    `toml:"status-port"` //table infomation from tidb status port
	PdAddr     string `toml:"pd-addr"`
}

var contextDataTransportKey = "dataTransportInfo"
var defaultTransportDirPrefix = "/tmp/tiem/transport" //todo: move to config

func ExportDataPreCheck(req *proto.DataExportRequest) error {
	if req.GetClusterId() == "" {
		return fmt.Errorf("invalid param clusterId %s", req.GetClusterId())
	}
	if req.GetUserName() == "" {
		return fmt.Errorf("invalid param userName %s", req.GetUserName())
	}
	/*
		if req.GetPassword() == "" {
			return fmt.Errorf("invalid param password %s", req.GetPassword())
		}
	*/

	if FileTypeCSV != req.GetFileType() && FileTypeSQL != req.GetFileType() {
		return fmt.Errorf("invalid param fileType %s", req.GetFileType())
	}
	switch req.GetStorageType() {
	case S3StorageType:
		if req.GetEndpointUrl() == "" {
			return fmt.Errorf("invalid param endpointUrl %s", req.GetEndpointUrl())
		}
		if req.GetBucketUrl() == "" {
			return fmt.Errorf("invalid param bucketUrl %s", req.GetBucketUrl())
		}
		if req.GetAccessKey() == "" {
			return fmt.Errorf("invalid param accessKey %s", req.GetAccessKey())
		}
		if req.GetSecretAccessKey() == "" {
			return fmt.Errorf("invalid param secretAccessKey %s", req.GetSecretAccessKey())
		}
	case NfsStorageType:
		if req.GetFilePath() == "" {
			return fmt.Errorf("invalid param filePath %s", req.GetFilePath())
		}
	default:
		return fmt.Errorf("invalid param storageType %s", req.GetStorageType())
	}

	return nil
}

func ImportDataPreCheck(req *proto.DataImportRequest) error {
	if req.GetClusterId() == "" {
		return fmt.Errorf("invalid param clusterId %s", req.GetClusterId())
	}
	if req.GetUserName() == "" {
		return fmt.Errorf("invalid param userName %s", req.GetUserName())
	}
	/*
		if req.GetPassword() == "" {
			return fmt.Errorf("invalid param password %s", req.GetPassword())
		}
	*/
	if req.GetFilePath() == "" {
		return fmt.Errorf("invalid param filePath %s", req.GetFilePath())
	}
	if S3StorageType != req.GetStorageType() && NfsStorageType != req.GetStorageType() {
		return fmt.Errorf("invalid param storageType %s", req.GetStorageType())
	}

	return nil
}

func ExportData(request *proto.DataExportRequest) (string, error) {
	getLogger().Infof("begin exportdata request %+v", request)
	defer getLogger().Infof("end exportdata")
	//todo: check operator
	operator := parseOperatorFromDTO(request.GetOperator())
	getLogger().Info(operator)
	clusterAggregation, err := ClusterRepo.Load(request.GetClusterId())
	if err != nil {
		getLogger().Errorf("load cluster %s aggregation from metadb failed", request.GetClusterId())
		return "", err
	}

	req := &db.DBCreateTransportRecordRequest{
		Record: &db.TransportRecordDTO{
			ClusterId:     request.GetClusterId(),
			TenantId:      operator.TenantId,
			TransportType: string(TransportTypeExport),
			FilePath:      getDataExportFilePath(request),
			Status:        TransportStatusRunning,
			StartTime:     time.Now().Unix(),
			EndTime:       time.Now().Unix(),
		},
	}
	resp, err := client.DBClient.CreateTransportRecord(ctx.Background(), req)
	if err != nil {
		return "", err
	}

	info := &ExportInfo{
		ClusterId:   request.GetClusterId(),
		UserName:    request.GetUserName(),
		Password:    request.GetPassword(), //todo: need encrypt
		FileType:    request.GetFileType(),
		RecordId:    resp.GetId(),
		FilePath:    getDataExportFilePath(request),
		Filter:      request.GetFilter(),
		Sql:  		 request.GetSql(),
		StorageType: request.GetStorageType(),
		BucketRegion: request.GetBucketRegion(),
	}

	// Start the workflow
	flow, err := CreateFlowWork(request.GetClusterId(), FlowExportData, operator)
	if err != nil {
		return "", err
	}
	flow.AddContext(contextClusterKey, clusterAggregation)
	flow.AddContext(contextDataTransportKey, info)
	flow.Start()

	clusterAggregation.CurrentWorkFlow = flow.FlowWork
	ClusterRepo.Persist(clusterAggregation)
	return info.RecordId, nil
}

func ImportData(request *proto.DataImportRequest) (string, error) {
	getLogger().Infof("begin importdata request %+v", request)
	defer getLogger().Infof("end importdata")
	//todo: check operator
	operator := parseOperatorFromDTO(request.GetOperator())
	getLogger().Info(operator)

	clusterAggregation, err := ClusterRepo.Load(request.GetClusterId())
	if err != nil {
		getLogger().Errorf("load cluster %s aggregation from metadb failed", request.GetClusterId())
		return "", err
	}

	req := &db.DBCreateTransportRecordRequest{
		Record: &db.TransportRecordDTO{
			ClusterId:     request.GetClusterId(),
			TenantId:      operator.TenantId,
			TransportType: string(TransportTypeImport),
			FilePath:      request.GetFilePath(),
			Status:        TransportStatusRunning,
			StartTime:     time.Now().Unix(),
			EndTime:       time.Now().Unix(),
		},
	}
	resp, err := client.DBClient.CreateTransportRecord(ctx.Background(), req)
	if err != nil {
		return "", err
	}
	info := &ImportInfo{
		ClusterId:   request.GetClusterId(),
		UserName:    request.GetUserName(),
		Password:    request.GetPassword(), //todo: need encrypt
		FilePath:    request.GetFilePath(),
		RecordId:    resp.GetId(),
		StorageType: request.GetStorageType(),
		ConfigPath:  getDataImportConfigDir(request.GetClusterId(), TransportTypeImport),
	}

	// Start the workflow
	flow, err := CreateFlowWork(request.GetClusterId(), FlowImportData, operator)
	if err != nil {
		return "", err
	}
	flow.AddContext(contextClusterKey, clusterAggregation)
	flow.AddContext(contextDataTransportKey, info)
	flow.Start()

	clusterAggregation.CurrentWorkFlow = flow.FlowWork
	ClusterRepo.Persist(clusterAggregation)
	return info.RecordId, nil
}

func DescribeDataTransportRecord(ope *proto.OperatorDTO, recordId, clusterId string, page, pageSize int32) ([]*db.TransportRecordDTO, *db.DBPageDTO, error) {
	getLogger().Infof("begin DescribeDataTransportRecord clusterId: %s, recordId: %s, page: %d, pageSize: %d", clusterId, recordId, page, pageSize)
	defer getLogger().Info("end DescribeDataTransportRecord")
	req := &db.DBListTransportRecordRequest{
		Page: &db.DBPageDTO{
			Page:     page,
			PageSize: pageSize,
		},
		ClusterId: clusterId,
		RecordId:  recordId,
	}
	resp, err := client.DBClient.ListTrasnportRecord(ctx.Background(), req)
	if err != nil {
		return nil, nil, err
	}

	return resp.GetRecords(), resp.GetPage(), nil
}

func convertTomlConfig(clusterAggregation *ClusterAggregation, info *ImportInfo) *DataImportConfig {
	getLogger().Info("begin convertTomlConfig")
	defer getLogger().Info("end convertTomlConfig")
	if clusterAggregation == nil || clusterAggregation.CurrentTopologyConfigRecord == nil {
		return nil
	}
	configModel := clusterAggregation.CurrentTopologyConfigRecord.ConfigModel
	if configModel == nil || configModel.TiDBServers == nil || configModel.PDServers == nil {
		return nil
	}
	tidbServer := configModel.TiDBServers[0]
	pdServer := configModel.PDServers[0]

	tidbServerPort := tidbServer.Port
	if tidbServerPort == 0 {
		tidbServerPort = DefaultTidbPort
	}

	tidbStatusPort := tidbServer.StatusPort
	if tidbStatusPort == 0 {
		tidbStatusPort = DefaultTidbStatusPort
	}

	pdClientPort := pdServer.ClientPort
	if pdClientPort == 0 {
		pdClientPort = DefaultPDClientPort
	}

	/*
	 * todo: sorted-kv-dir and data-source-dir in the same disk, may slow down import performance,
	 *  and check-requirements = true can not pass lightning pre-check
	 *  in real environment, config data-source-dir = user nfs storage, sorted-kv-dir = other disk, turn on pre-check
	 */
	config := &DataImportConfig{
		Lightning: LightningCfg{
			Level:             "info",
			File:              fmt.Sprintf("%s/tidb-lightning.log", info.ConfigPath),
			CheckRequirements: false, //todo: TBD
		},
		TikvImporter: TikvImporterCfg{
			Backend:     BackendLocal,
			SortedKvDir: info.ConfigPath, //todo: TBD
		},
		MyDumper: MyDumperCfg{
			DataSourceDir: info.FilePath,
		},
		Tidb: TidbCfg{
			Host:       tidbServer.Host,
			Port:       tidbServerPort,
			User:       info.UserName,
			Password:   info.Password,
			StatusPort: tidbStatusPort,
			PdAddr:     fmt.Sprintf("%s:%d", pdServer.Host, pdClientPort),
		},
	}
	return config
}


func getDataImportConfigDir(clusterId string, transportType TransportType) string {
	return fmt.Sprintf("%s/%s/%s", defaultTransportDirPrefix, clusterId, transportType)
}

func getDataExportFilePath(request *proto.DataExportRequest) string {
	var filePath string
	if S3StorageType == request.GetStorageType() {
		filePath = fmt.Sprintf("%s?access-key=%s&secret-access-key=%s&endpoint=%s&force-path-style=true", request.GetBucketUrl(), request.GetAccessKey(), request.GetSecretAccessKey(), request.GetEndpointUrl())
	} else {
		filePath = request.GetFilePath()
	}
	return filePath
}

func cleanDataTransportDir(filepath string) error {
	getLogger().Infof("clean and re-mkdir data dir: %s", filepath)
	if err := os.RemoveAll(filepath); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func buildDataImportConfig(task *TaskEntity, context *FlowContext) bool {
	getLogger().Info("begin buildDataImportConfig")
	defer getLogger().Info("end buildDataImportConfig")

	clusterAggregation := context.value(contextClusterKey).(*ClusterAggregation)
	info := context.value(contextDataTransportKey).(*ImportInfo)

	if err := cleanDataTransportDir(info.ConfigPath); err != nil {
		getLogger().Errorf("clean import directory failed, %s", err.Error())
		return false
	}

	config := convertTomlConfig(clusterAggregation, info)
	if config == nil {
		getLogger().Errorf("convert toml config failed, cluster: %v", clusterAggregation)
		return false
	}
	filePath := fmt.Sprintf("%s/tidb-lightning.toml", info.ConfigPath)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
	if err != nil {
		getLogger().Errorf("create import toml config failed, %s", err.Error())
		return false
	}

	if err = toml.NewEncoder(file).Encode(config); err != nil {
		getLogger().Errorf("decode data import toml config failed, %s", err.Error())
		return false
	}
	getLogger().Infof("build lightning toml file sucess, %v", config)

	return true
}

func importDataToCluster(task *TaskEntity, context *FlowContext) bool {
	getLogger().Info("begin importDataToCluster")
	defer getLogger().Info("end importDataToCluster")

	info := context.value(contextDataTransportKey).(*ImportInfo)

	//tiup tidb-lightning -config tidb-lightning.toml
	//todo: tiupmgr not return failed err
	resp, err := libtiup.MicroSrvTiupLightning(0,
		[]string{"-config", fmt.Sprintf("%s/tidb-lightning.toml", info.ConfigPath)},
		uint64(task.Id))
	if err != nil {
		getLogger().Errorf("call tiup lightning api failed, %s", err.Error())
		return false
	}
	getLogger().Infof("call tiupmgr tidb-lightning api success, %v", resp)

	return true
}

func updateDataImportRecord(task *TaskEntity, context *FlowContext) bool {
	getLogger().Info("begin updateDataImportRecord")
	defer getLogger().Info("end updateDataImportRecord")
	clusterAggregation := context.value(contextClusterKey).(*ClusterAggregation)
	info := context.value(contextDataTransportKey).(*ImportInfo)
	cluster := clusterAggregation.Cluster

	req := &db.DBUpdateTransportRecordRequest{
		Record: &db.TransportRecordDTO{
			ID:        info.RecordId,
			ClusterId: cluster.Id,
			Status:    TransportStatusSuccess,
			EndTime:   time.Now().Unix(),
		},
	}
	resp, err := client.DBClient.UpdateTransportRecord(ctx.Background(), req)
	if err != nil {
		getLogger().Errorf("update data transport record failed, %s", err.Error())
		return false
	}
	getLogger().Infof("update data transport record success, %v", resp)
	return true
}

func exportDataFromCluster(task *TaskEntity, context *FlowContext) bool {
	getLogger().Info("begin exportDataFromCluster")
	defer getLogger().Info("end exportDataFromCluster")

	clusterAggregation := context.value(contextClusterKey).(*ClusterAggregation)
	info := context.value(contextDataTransportKey).(*ExportInfo)
	configModel := clusterAggregation.CurrentTopologyConfigRecord.ConfigModel
	tidbServer := configModel.TiDBServers[0]
	tidbServerPort := tidbServer.Port
	if tidbServerPort == 0 {
		tidbServerPort = DefaultTidbPort
	}

	if NfsStorageType == info.StorageType {
		if err := cleanDataTransportDir(info.FilePath); err != nil {
			getLogger().Errorf("clean export directory failed, %s", err.Error())
			return false
		}
	}

	//tiup dumpling -u root -P 4000 --host 127.0.0.1 --filetype sql -t 8 -o /tmp/test -r 200000 -F 256MiB --filter "user*"
	//todo: admin root password
	//todo: tiupmgr not return failed err
	cmd := []string{"-u", info.UserName,
		"-p", info.Password,
		"-P", strconv.Itoa(tidbServerPort),
		"--host", tidbServer.Host,
		"--filetype", info.FileType,
		"-t", "8",
		"-o", fmt.Sprintf("%s", info.FilePath),
		"-r", "200000",
		"-F", "256MiB"}
	if info.Filter != "" {
		cmd = append(cmd, "--filter", fmt.Sprintf("\"%s\"", info.Filter))
	}
	if FileTypeCSV == info.FileType && info.Sql != "" {
		cmd = append(cmd, "--sql", fmt.Sprintf("\"%s\"", info.Sql))
	}
	if S3StorageType == info.StorageType && info.BucketRegion != "" {
		cmd = append(cmd, "--s3.region", fmt.Sprintf("\"%s\"", info.BucketRegion))
	}
	getLogger().Infof("call tiupmgr dumpling api, cmd: %v", cmd)
	resp, err := libtiup.MicroSrvTiupDumpling(0, cmd, uint64(task.Id))
	if err != nil {
		getLogger().Errorf("call tiup dumpling api failed, %s", err.Error())
		return false
	}

	getLogger().Infof("call tiupmgr succee, resp: %v", resp)

	return true
}

func updateDataExportRecord(task *TaskEntity, context *FlowContext) bool {
	getLogger().Info("begin updateDataExportRecord")
	defer getLogger().Info("end updateDataExportRecord")
	clusterAggregation := context.value(contextClusterKey).(*ClusterAggregation)
	info := context.value(contextDataTransportKey).(*ExportInfo)
	cluster := clusterAggregation.Cluster

	req := &db.DBUpdateTransportRecordRequest{
		Record: &db.TransportRecordDTO{
			ID:        info.RecordId,
			ClusterId: cluster.Id,
			Status:    TransportStatusSuccess,
			EndTime:   time.Now().Unix(),
		},
	}
	resp, err := client.DBClient.UpdateTransportRecord(ctx.Background(), req)
	if err != nil {
		getLogger().Errorf("update data transport record failed, %s", err.Error())
		return false
	}
	getLogger().Infof("update data transport record success, %v", resp)
	return true
}

/*
func compressExportData(task *TaskEntity, context *FlowContext) bool {
	getLogger().Info("begin compressExportData")
	defer getLogger().Info("end compressExportData")

	info := context.value(contextDataTransportKey).(*ExportInfo)

	dataDir := fmt.Sprintf("%s/data", info.FilePath)
	dataZipDir := fmt.Sprintf("%s/data.zip", info.FilePath)
	if err := zipDir(dataDir, dataZipDir); err != nil {
		getLogger().Errorf("compress export data failed, %s", err.Error())
		return false
	}

	return true
}

func deCompressImportData(task *TaskEntity, context *FlowContext) bool {
	getLogger().Info("begin deCompressImportData")
	defer getLogger().Info("end deCompressImportData")

	info := context.value(contextDataTransportKey).(*ImportInfo)

	dataDir := fmt.Sprintf("%s/data", info.ConfigPath)
	dataZipDir := info.FilePath
	if err := unzipDir(dataZipDir, dataDir); err != nil {
		getLogger().Errorf("deCompress import data failed, %s", err.Error())
		return false
	}

	return true
}
*/

func importDataFailed(task *TaskEntity, context *FlowContext) bool {
	getLogger().Info("begin importDataFailed")
	defer getLogger().Info("end importDataFailed")
	clusterAggregation := context.value(contextClusterKey).(ClusterAggregation)
	info := context.value(contextDataTransportKey).(*ImportInfo)
	cluster := clusterAggregation.Cluster

	if err := updateTransportRecordFailed(info.RecordId, cluster.Id); err != nil {
		return false
	}

	return ClusterFail(task, context)
}

func exportDataFailed(task *TaskEntity, context *FlowContext) bool {
	getLogger().Info("begin exportDataFailed")
	defer getLogger().Info("end exportDataFailed")
	clusterAggregation := context.value(contextClusterKey).(ClusterAggregation)
	info := context.value(contextDataTransportKey).(*ExportInfo)
	cluster := clusterAggregation.Cluster

	if err := updateTransportRecordFailed(info.RecordId, cluster.Id); err != nil {
		return false
	}

	return ClusterFail(task, context)
}

func updateTransportRecordFailed(recordId, clusterId string) error {
	req := &db.DBUpdateTransportRecordRequest{
		Record: &db.TransportRecordDTO{
			ID:        recordId,
			ClusterId: clusterId,
			Status:    TransportStatusFailed,
			EndTime:   time.Now().Unix(),
		},
	}
	resp, err := client.DBClient.UpdateTransportRecord(ctx.Background(), req)
	if err != nil {
		getLogger().Errorf("update data transport record failed, %s", err.Error())
		return err
	}
	getLogger().Infof("update data transport record success, %v", resp)
	return nil
}

/*
func zipDir(dir string, zipFile string) error {
	getLogger().Infof("begin zipDir: dir[%s] to file[%s]", dir, zipFile)
	defer getLogger().Info("end zipDir")
	fz, err := os.Create(zipFile)
	if err != nil {
		return fmt.Errorf("Create zip file failed: %s", err.Error())
	}
	defer fz.Close()

	w := zip.NewWriter(fz)
	defer w.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			relPath := strings.TrimPrefix(path, filepath.Dir(path))
			fDest, err := w.Create(relPath)
			if err != nil {
				return fmt.Errorf("zip Create failed: %s", err.Error())
			}
			fSrc, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("zip Open failed: %s", err.Error())
			}
			defer fSrc.Close()
			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				return fmt.Errorf("zip Copy failed: %s", err.Error())
			}
		}
		return nil
	})
	if err != nil {
		getLogger().Errorf("filepath walk failed, %s", err.Error())
		return err
	}

	return nil
}

func unzipDir(zipFile string, dir string) error {
	getLogger().Infof("begin unzipDir: file[%s] to dir[%s]", zipFile, dir)
	defer getLogger().Info("end unzipDir")
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("Open zip file failed: %s", err.Error())
	}
	defer r.Close()

	for _, f := range r.File {
		func() {
			path := dir + string(filepath.Separator) + f.Name
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				getLogger().Errorf("make filepath failed: %s", err.Error())
				return
			}
			fDest, err := os.Create(path)
			if err != nil {
				getLogger().Errorf("unzip Create failed: %s", err.Error())
				return
			}
			defer fDest.Close()

			fSrc, err := f.Open()
			if err != nil {
				getLogger().Errorf("unzip Open failed: %s", err.Error())
				return
			}
			defer fSrc.Close()

			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				getLogger().Errorf("unzip Copy failed: %s", err.Error())
				return
			}
		}()
	}
	return nil
}
*/
