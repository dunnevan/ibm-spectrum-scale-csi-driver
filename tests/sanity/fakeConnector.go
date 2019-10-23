/**
 * Copyright 2019 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sanity

// import (
// 	"hash/fnv"

// 	"github.ibm.com/FSaaS/explore-csi-scale-driver/pkg/scale"
// 	"github.ibm.com/FSaaS/explore-csi-scale-driver/pkg/scale/connectors"
// 	"github.ibm.com/FSaaS/explore-csi-scale-driver/pkg/scale/connectors/api/v2"
// )

// type fakeConnectorFactory struct{}

// type fakeConnector struct {
// 	filesystems   map[string]*api.FileSystem_v2
// 	filesets      map[fsFsetKey]*api.Fileset_v2
// 	filesetIDs    map[int]fsFsetKey
// 	filesetQuotas map[fsFsetKey]*api.Quota_v2
// 	directories   map[string]api.OwnerInfo
// }

// type fsFsetKey struct {
// 	fs   string
// 	fset string
// }

// func hash(s string) uint32 {
// 	h := fnv.New32a()
// 	h.Write([]byte(s))
// 	return h.Sum32()
// }

// func newFakeConnector() *fakeConnector {
// 	return &fakeConnector{
// 		filesystems:   make(map[string]*api.FileSystem_v2),
// 		filesets:      make(map[fsFsetKey]*api.Fileset_v2),
// 		filesetIDs:    make(map[int]fsFsetKey),
// 		filesetQuotas: make(map[fsFsetKey]*api.Quota_v2),
// 		directories:   make(map[string]api.OwnerInfo),
// 	}
// }

// func (fake fakeConnectorFactory) NewConnector(in connectors.HasClusterID) connectors.Connector {
// 	return newFakeConnector()
// }

// func (fake *fakeConnector) GetClusterId() (string, error) {
// 	return "", nil
// }

// func (fake *fakeConnector) GetFilesystemByUUID(fsuuid string) (*api.FileSystem_v2, error) {
// 	//fsuuid == filesystemName in our fake env
// 	fs, ok := fake.filesystems[fsuuid]
// 	if !ok {
// 		return nil, &scale.NotFoundError{}
// 	}

// 	return fs, nil
// }
// func (fake *fakeConnector) GetFilesystemByName(filesystemName string) (*api.FileSystem_v2, error) {
// 	//fsuuid == filesystemName in our fake env
// 	fs, ok := fake.filesystems[filesystemName]
// 	if !ok {
// 		return nil, &scale.NotFoundError{}
// 	}

// 	return fs, nil
// }
// func (fake *fakeConnector) MountFilesystem(filesystemName string, nodeName string) error {
// 	return nil
// }
// func (fake *fakeConnector) UnmountFilesystem(filesystemName string, nodeName string) error {
// 	return nil
// }

// func (fake *fakeConnector) GetFilesetByID(filesystemName string, filesetID int) (*api.Fileset_v2, error) {

// 	if filesetID == 0 {
// 		return &api.Fileset_v2{
// 			FilesetName: "root",
// 		}, nil
// 	}

// 	fsFsetKey := fake.filesetIDs[filesetID]
// 	fset, ok := fake.filesets[fsFsetKey]
// 	if !ok {
// 		return nil, &scale.NotFoundError{}
// 	}

// 	return fset, nil
// }
// func (fake *fakeConnector) GetFilesetByName(filesystemName string, filesetName string) (*api.Fileset_v2, error) {

// 	if filesetName == "root" {
// 		return &api.Fileset_v2{
// 			FilesetName: "root",
// 		}, nil
// 	}

// 	fset, ok := fake.filesets[fsFsetKey{filesystemName, filesetName}]
// 	if !ok {
// 		return nil, &scale.NotFoundError{}
// 	}

// 	return fset, nil
// }

// func (fake *fakeConnector) CreateFileset(filesystemName string, req api.CreateFilesetRequest) error {

// 	fsetName := req.FilesetName
// 	fsetNameHash := int(hash(fsetName))
// 	fsFsetKey := fsFsetKey{filesystemName, fsetName}

// 	fake.filesetIDs[fsetNameHash] = fsFsetKey
// 	fake.filesets[fsFsetKey] = &api.Fileset_v2{
// 		FilesetName: fsetName,
// 		Config: api.FilesetConfig_v2{
// 			Id: fsetNameHash,
// 		},
// 	}
// 	return nil
// }
// func (fake *fakeConnector) DeleteFileset(filesystemName string, filesetName string) error {
// 	delete(fake.filesets, fsFsetKey{filesystemName, filesetName})
// 	return nil
// }

// func (fake *fakeConnector) LinkFileset(filesystemName string, filesetName string, linkpath string) error {
// 	fset := fake.filesets[fsFsetKey{filesystemName, filesetName}]
// 	fset.Config.Path = linkpath
// 	return nil
// }
// func (fake *fakeConnector) UnlinkFileset(filesystemName string, filesetName string) error {
// 	fset := fake.filesets[fsFsetKey{filesystemName, filesetName}]
// 	fset.Config.Path = ""
// 	return nil
// }

// func (fake *fakeConnector) GetFilesetQuota(filesystemName string, filesetName string) (*api.Quota_v2, error) {
// 	quota, ok := fake.filesetQuotas[fsFsetKey{filesystemName, filesetName}]
// 	if !ok {
// 		return nil, &scale.NotFoundError{}
// 	}

// 	return quota, nil
// }
// func (fake *fakeConnector) SetFilesetQuota(filesystemName string, filesetName string, quotaBytes int64) error {
// 	fake.filesetQuotas[fsFsetKey{filesystemName, filesetName}] = &api.Quota_v2{
// 		BlockQuota: int(quotaBytes),
// 	}
// 	return nil
// }
// func (fake *fakeConnector) ListQuotas(filesystem string) (*api.GetQuotaResponse_v2, error) {
// 	return nil, nil
// }

// func (fake *fakeConnector) MakeDirectory(filesystemName string, relPath string, uid int, gid int) error {
// 	fake.directories[relPath] = api.OwnerInfo{UID: uid, GID: gid}
// 	return nil
// }
// func (fake *fakeConnector) DeleteDirectory(filesystemName string, relPath string) error {
// 	delete(fake.directories, relPath)
// 	return nil
// }
// func (fake *fakeConnector) CheckIfFileDirPresent(filesystemName string, relPath string) (bool, error) {
// 	_, ok := fake.directories[relPath]
// 	return ok, nil
// }
