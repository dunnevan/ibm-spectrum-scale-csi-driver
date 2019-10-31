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

package integration

import (
	"os"
	"testing"

	scale "github.com/IBM/ibm-spectrum-scale-csi-driver/csiplugin"
	"github.com/IBM/ibm-spectrum-scale-csi-driver/csiplugin/connectors"
	"github.com/IBM/ibm-spectrum-scale-csi-driver/csiplugin/settings"
	"github.com/kubernetes-csi/csi-test/pkg/sanity"
	"k8s.io/klog"
)

func TestScaleDriver(t *testing.T) {
	var address = os.TempDir() + "/csi.sock"
	/*
		var parameters = map[string]string{
			connectors.UserSpecifiedInodeLimit:   "1024",
			connectors.UserSpecifiedVolBackendFs: "fs1",
			connectors.UserSpecifiedClusterId:    "16482346744146153652",
		}
	*/

	/*
		var parameters = map[string]string{
			connectors.UserSpecifiedFilesetType:  "dependant",
			connectors.UserSpecifiedInodeLimit:   "1024",
			connectors.UserSpecifiedVolBackendFs: "fs1",
			connectors.UserSpecifiedParentFset:   "test_csi",
			connectors.UserSpecifiedClusterId:    "16482346744146153652",
		}
	*/

	/*
		var parameters = map[string]string{
			connectors.UserSpecifiedVolBackendFs: "fs1",
			connectors.UserSpecifiedVolDirPath:   "/ibm/fs1/test_dir/",
			connectors.UserSpecifiedClusterId:    "16482346744146153652",
		}
	*/

	var parameters = map[string]string{
		connectors.UserSpecifiedInodeLimit:   "1024",
		connectors.UserSpecifiedVolBackendFs: "fs1",
		connectors.UserSpecifiedClusterId:    "16482346744146153652",
		connectors.UserSpecifiedVolDirPath:   "/ibm/fs1/test_dir/",
	}

	config := &sanity.Config{
		Address:              address,
		TestVolumeParameters: parameters,
	}

	configMap := settings.ScaleSettingsConfigMap{
		Clusters: []settings.Clusters{
			{
				Primary: settings.Primary{
					PrimaryCid:      "16482346744146153652",
					PrimaryFS:       "fs1",
					PrimaryFset:     "csifset1",
					PrimaryFSMount:  "/ibm/fs1",
					PrimaryFsetLink: "/ibm/fs1/csifset1",
				},
				ID: "16482346744146153652",
				RestAPI: []settings.RestAPI{
					{
						GuiHost: "edunn-k8s-master.fyre.ibm.com",
						GuiPort: 443,
					},
				},
				Secrets:      "guisecret",
				MgmtUsername: "evan",
				MgmtPassword: "Passw0rd",
			},
		},
	}

	os.Setenv("SCALE_HOSTPATH", "/ibm/fs1")

	driver := scale.GetScaleDriver()
	err := driver.SetupScaleDriver(
		"csi-scale-test",
		"explodable",
		"edunn-k8s-master.fyre.ibm.com",
		configMap,
	)
	if err != nil {
		klog.Fatalf("Failed to initialize Scale CSI Driver: %v", err)
	}

	go driver.Run("unix://" + address)

	// Now call the test suite
	sanity.Test(t, config)
}
