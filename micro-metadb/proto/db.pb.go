// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.3.0
// source: db.proto

package db

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_db_proto protoreflect.FileDescriptor

var file_db_proto_rawDesc = []byte{
	0x0a, 0x08, 0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x64, 0x62, 0x5f, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x64, 0x62, 0x5f, 0x68, 0x6f,
	0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x64, 0x62, 0x5f, 0x63, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x64, 0x62, 0x5f, 0x74,
	0x69, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xad, 0x09, 0x0a, 0x0d, 0x54, 0x69,
	0x43, 0x50, 0x44, 0x42, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x46,
	0x69, 0x6e, 0x64, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x14, 0x2e, 0x44, 0x42, 0x46, 0x69,
	0x6e, 0x64, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x15, 0x2e, 0x44, 0x42, 0x46, 0x69, 0x6e, 0x64, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x0b, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x15, 0x2e, 0x44, 0x42, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x44,
	0x42, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x09, 0x53, 0x61, 0x76, 0x65, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x13, 0x2e, 0x44, 0x42, 0x53, 0x61, 0x76, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x44, 0x42, 0x53, 0x61, 0x76, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x09,
	0x46, 0x69, 0x6e, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x13, 0x2e, 0x44, 0x42, 0x46, 0x69,
	0x6e, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x44, 0x42, 0x46, 0x69, 0x6e, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5a, 0x0a, 0x15, 0x46, 0x69, 0x6e, 0x64, 0x52, 0x6f, 0x6c, 0x65,
	0x73, 0x42, 0x79, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x2e,
	0x44, 0x42, 0x46, 0x69, 0x6e, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x50, 0x65, 0x72,
	0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20,
	0x2e, 0x44, 0x42, 0x46, 0x69, 0x6e, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x50, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x30, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x11, 0x2e, 0x44, 0x42,
	0x41, 0x64, 0x64, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12,
	0x2e, 0x44, 0x42, 0x41, 0x64, 0x64, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x48, 0x6f, 0x73, 0x74,
	0x12, 0x14, 0x2e, 0x44, 0x42, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x44, 0x42, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a,
	0x08, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x13, 0x2e, 0x44, 0x42, 0x4c, 0x69,
	0x73, 0x74, 0x48, 0x6f, 0x73, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x44, 0x42, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x6f, 0x73, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x73, 0x12, 0x16, 0x2e, 0x44, 0x42, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x44, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x44,
	0x42, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0a, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x48, 0x6f,
	0x73, 0x74, 0x73, 0x12, 0x14, 0x2e, 0x44, 0x42, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x48, 0x6f, 0x73,
	0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x44, 0x42, 0x41, 0x6c,
	0x6c, 0x6f, 0x63, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3f, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x17, 0x2e,
	0x44, 0x42, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x44, 0x42, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x3c, 0x0a, 0x0b, 0x46, 0x69, 0x6e, 0x64, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12,
	0x15, 0x2e, 0x44, 0x42, 0x46, 0x69, 0x6e, 0x64, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x44, 0x42, 0x46, 0x69, 0x6e, 0x64, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4b,
	0x0a, 0x10, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x55, 0x50, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x1a, 0x2e, 0x44, 0x42, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x55,
	0x50, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b,
	0x2e, 0x44, 0x42, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x55, 0x50, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x0b, 0x4c,
	0x69, 0x73, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x15, 0x2e, 0x44, 0x42, 0x4c,
	0x69, 0x73, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x44, 0x42, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x16, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x75, 0x70,
	0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x16,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54,
	0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x47, 0x0a, 0x10, 0x46, 0x69, 0x6e, 0x64, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x42,
	0x79, 0x49, 0x44, 0x12, 0x18, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61,
	0x73, 0x6b, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x46, 0x69, 0x6e, 0x64, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x42, 0x79, 0x49, 0x44,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5f, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x54,
	0x69, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x79, 0x42,
	0x69, 0x7a, 0x49, 0x44, 0x12, 0x20, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x69, 0x75, 0x70, 0x54, 0x61,
	0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x79, 0x42, 0x69, 0x7a, 0x49, 0x44, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x69, 0x75, 0x70,
	0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x79, 0x42, 0x69, 0x7a, 0x49,
	0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b,
	0x64, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_db_proto_goTypes = []interface{}{
	(*DBFindTenantRequest)(nil),              // 0: DBFindTenantRequest
	(*DBFindAccountRequest)(nil),             // 1: DBFindAccountRequest
	(*DBSaveTokenRequest)(nil),               // 2: DBSaveTokenRequest
	(*DBFindTokenRequest)(nil),               // 3: DBFindTokenRequest
	(*DBFindRolesByPermissionRequest)(nil),   // 4: DBFindRolesByPermissionRequest
	(*DBAddHostRequest)(nil),                 // 5: DBAddHostRequest
	(*DBRemoveHostRequest)(nil),              // 6: DBRemoveHostRequest
	(*DBListHostsRequest)(nil),               // 7: DBListHostsRequest
	(*DBCheckDetailsRequest)(nil),            // 8: DBCheckDetailsRequest
	(*DBAllocHostsRequest)(nil),              // 9: DBAllocHostsRequest
	(*DBCreateClusterRequest)(nil),           // 10: DBCreateClusterRequest
	(*DBFindClusterRequest)(nil),             // 11: DBFindClusterRequest
	(*DBUpdateTiUPConfigRequest)(nil),        // 12: DBUpdateTiUPConfigRequest
	(*DBListClusterRequest)(nil),             // 13: DBListClusterRequest
	(*CreateTiupTaskRequest)(nil),            // 14: CreateTiupTaskRequest
	(*UpdateTiupTaskRequest)(nil),            // 15: UpdateTiupTaskRequest
	(*FindTiupTaskByIDRequest)(nil),          // 16: FindTiupTaskByIDRequest
	(*GetTiupTaskStatusByBizIDRequest)(nil),  // 17: GetTiupTaskStatusByBizIDRequest
	(*DBFindTenantResponse)(nil),             // 18: DBFindTenantResponse
	(*DBFindAccountResponse)(nil),            // 19: DBFindAccountResponse
	(*DBSaveTokenResponse)(nil),              // 20: DBSaveTokenResponse
	(*DBFindTokenResponse)(nil),              // 21: DBFindTokenResponse
	(*DBFindRolesByPermissionResponse)(nil),  // 22: DBFindRolesByPermissionResponse
	(*DBAddHostResponse)(nil),                // 23: DBAddHostResponse
	(*DBRemoveHostResponse)(nil),             // 24: DBRemoveHostResponse
	(*DBListHostsResponse)(nil),              // 25: DBListHostsResponse
	(*DBCheckDetailsResponse)(nil),           // 26: DBCheckDetailsResponse
	(*DBAllocHostResponse)(nil),              // 27: DBAllocHostResponse
	(*DBCreateClusterResponse)(nil),          // 28: DBCreateClusterResponse
	(*DBFindClusterResponse)(nil),            // 29: DBFindClusterResponse
	(*DBUpdateTiUPConfigResponse)(nil),       // 30: DBUpdateTiUPConfigResponse
	(*DBListClusterResponse)(nil),            // 31: DBListClusterResponse
	(*CreateTiupTaskResponse)(nil),           // 32: CreateTiupTaskResponse
	(*UpdateTiupTaskResponse)(nil),           // 33: UpdateTiupTaskResponse
	(*FindTiupTaskByIDResponse)(nil),         // 34: FindTiupTaskByIDResponse
	(*GetTiupTaskStatusByBizIDResponse)(nil), // 35: GetTiupTaskStatusByBizIDResponse
}
var file_db_proto_depIdxs = []int32{
	0,  // 0: TiCPDBService.FindTenant:input_type -> DBFindTenantRequest
	1,  // 1: TiCPDBService.FindAccount:input_type -> DBFindAccountRequest
	2,  // 2: TiCPDBService.SaveToken:input_type -> DBSaveTokenRequest
	3,  // 3: TiCPDBService.FindToken:input_type -> DBFindTokenRequest
	4,  // 4: TiCPDBService.FindRolesByPermission:input_type -> DBFindRolesByPermissionRequest
	5,  // 5: TiCPDBService.AddHost:input_type -> DBAddHostRequest
	6,  // 6: TiCPDBService.RemoveHost:input_type -> DBRemoveHostRequest
	7,  // 7: TiCPDBService.ListHost:input_type -> DBListHostsRequest
	8,  // 8: TiCPDBService.CheckDetails:input_type -> DBCheckDetailsRequest
	9,  // 9: TiCPDBService.AllocHosts:input_type -> DBAllocHostsRequest
	10, // 10: TiCPDBService.AddCluster:input_type -> DBCreateClusterRequest
	11, // 11: TiCPDBService.FindCluster:input_type -> DBFindClusterRequest
	12, // 12: TiCPDBService.UpdateTiUPConfig:input_type -> DBUpdateTiUPConfigRequest
	13, // 13: TiCPDBService.ListCluster:input_type -> DBListClusterRequest
	14, // 14: TiCPDBService.CreateTiupTask:input_type -> CreateTiupTaskRequest
	15, // 15: TiCPDBService.UpdateTiupTask:input_type -> UpdateTiupTaskRequest
	16, // 16: TiCPDBService.FindTiupTaskByID:input_type -> FindTiupTaskByIDRequest
	17, // 17: TiCPDBService.GetTiupTaskStatusByBizID:input_type -> GetTiupTaskStatusByBizIDRequest
	18, // 18: TiCPDBService.FindTenant:output_type -> DBFindTenantResponse
	19, // 19: TiCPDBService.FindAccount:output_type -> DBFindAccountResponse
	20, // 20: TiCPDBService.SaveToken:output_type -> DBSaveTokenResponse
	21, // 21: TiCPDBService.FindToken:output_type -> DBFindTokenResponse
	22, // 22: TiCPDBService.FindRolesByPermission:output_type -> DBFindRolesByPermissionResponse
	23, // 23: TiCPDBService.AddHost:output_type -> DBAddHostResponse
	24, // 24: TiCPDBService.RemoveHost:output_type -> DBRemoveHostResponse
	25, // 25: TiCPDBService.ListHost:output_type -> DBListHostsResponse
	26, // 26: TiCPDBService.CheckDetails:output_type -> DBCheckDetailsResponse
	27, // 27: TiCPDBService.AllocHosts:output_type -> DBAllocHostResponse
	28, // 28: TiCPDBService.AddCluster:output_type -> DBCreateClusterResponse
	29, // 29: TiCPDBService.FindCluster:output_type -> DBFindClusterResponse
	30, // 30: TiCPDBService.UpdateTiUPConfig:output_type -> DBUpdateTiUPConfigResponse
	31, // 31: TiCPDBService.ListCluster:output_type -> DBListClusterResponse
	32, // 32: TiCPDBService.CreateTiupTask:output_type -> CreateTiupTaskResponse
	33, // 33: TiCPDBService.UpdateTiupTask:output_type -> UpdateTiupTaskResponse
	34, // 34: TiCPDBService.FindTiupTaskByID:output_type -> FindTiupTaskByIDResponse
	35, // 35: TiCPDBService.GetTiupTaskStatusByBizID:output_type -> GetTiupTaskStatusByBizIDResponse
	18, // [18:36] is the sub-list for method output_type
	0,  // [0:18] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_db_proto_init() }
func file_db_proto_init() {
	if File_db_proto != nil {
		return
	}
	file_db_auth_proto_init()
	file_db_host_proto_init()
	file_db_cluster_proto_init()
	file_db_tiup_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_db_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_db_proto_goTypes,
		DependencyIndexes: file_db_proto_depIdxs,
	}.Build()
	File_db_proto = out.File
	file_db_proto_rawDesc = nil
	file_db_proto_goTypes = nil
	file_db_proto_depIdxs = nil
}
