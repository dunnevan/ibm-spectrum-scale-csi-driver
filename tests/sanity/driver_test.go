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
// 	"os"
// 	"testing"

// 	"github.com/kubernetes-csi/csi-test/pkg/sanity"
// 	"github.ibm.com/FSaaS/explore-csi-scale-driver/pkg/scale"
// 	"github.ibm.com/FSaaS/explore-csi-scale-driver/pkg/scale/settings"
// )

// func TestScaleDriver(t *testing.T) {
// 	var address = os.TempDir() + "/csi.sock"
// 	/*
// 		var parameters = map[string]string{
// 			settings.Type:            "independant",
// 			settings.ClusterID:       "16482346744146153652",
// 			settings.FilesystemName:  "fs1",
// 			settings.FilesetJunction: ".csi/",
// 			settings.VolumePath:      "volume/",
// 			settings.InodeLimit:		"1024",
// 		}
// 	*/

// 	/*
// 		var parameters = map[string]string{
// 			settings.Type:              "directory",
// 			settings.ClusterID:         "16482346744146153652",
// 			settings.FilesystemName:    "fs1",
// 			settings.ParentFilesetName: "test_csi",
// 			settings.VolumePath:        "volumes/",
// 			settings.InodeLimit:		"1024",
// 		}
// 	*/

// 	/*
// 		var parameters = map[string]string{
// 			settings.Type:              "dependant",
// 			settings.ClusterID:         "16482346744146153652",
// 			settings.FilesystemName:    "fs1",
// 			settings.ParentFilesetName: "test_csi",
// 			settings.FilesetJunction:   ".csi/",
// 			settings.VolumePath:        "volume/",
// 		}
// 	*/

// 	var (
// 		fsName     = "fs1"
// 		parameters = map[string]string{
// 			settings.Type:            "independant",
// 			settings.ClusterID:       "16482346744146153652",
// 			settings.FilesystemName:  fsName,
// 			settings.FilesetJunction: ".csi/",
// 			settings.VolumePath:      "volume/",
// 			settings.InodeLimit:      "1024",
// 		}
// 	)

// 	config := &sanity.Config{
// 		Address:              address,
// 		TestVolumeParameters: parameters,
// 	}

// 	driver := scale.NewFakeDriver(fakeConnectorFactory{})

// 	go driver.Run("unix://" + address)

// 	// Now call the test suite
// 	sanity.Test(t, config)
// }
