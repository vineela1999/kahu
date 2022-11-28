/*
Copyright 2022 The SODA Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package deployment

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"

	k8s "github.com/soda-cdm/kahu/test/e2e/util/k8s"
	kahu "github.com/soda-cdm/kahu/test/e2e/util/kahu"
)

//testcase for E2E deployment backup and restore
var _ = Describe("DaemonsetBackup", Label("Daemonset"), func() {
	Context("Create backup of Daemonset and restore", func() {
		It("Daemonset", func() {

			kubeClient, kahuClient := kahu.Clients()
			//Create Daemonset to test
			ns := kahu.BackupNameSpace
			labels := make(map[string]string)
			image := "nginx"
			replicas := 2
			labels["daemonset"] = "daemonset"

			UUIDgen, err := uuid.NewRandom()
			Expect(err).To(BeNil())

			name := "daemonset" + "-" + UUIDgen.String()
			daemonSet, err := k8s.NewDaemonset(name, ns, int32(replicas), labels, image)
			log.Infof("daemonset:%v\n", daemonSet)
			daemonSet1, err := k8s.CreateDaemonset(kubeClient, ns, daemonSet)
			log.Debugf("daemonset:%v\n", daemonSet1)
			Expect(err).To(BeNil())
			err = k8s.WaitForDaemonsetComplete(kubeClient, ns, name)
			Expect(err).To(BeNil())

			//create backup for the daemonset
			backupName := "backup" + "daemonset" + "-" + UUIDgen.String()
			includeNs := kahu.BackupNameSpace
			resourceType := "DaemonSet"
			_, err = kahu.CreateBackup(kahuClient, backupName, includeNs, resourceType)
			Expect(err).To(BeNil())
			err = kahu.WaitForBackupCreate(kahuClient, backupName)
			Expect(err).To(BeNil())
			log.Infof("backup  of daemonset is done\n")

			// create restore for the backup
			restoreName := "restore" + "daemonset" + "-" + UUIDgen.String()
			nsRestore := kahu.RestoreNameSpace
			restore, err := kahu.CreateRestore(kahuClient, restoreName, backupName, includeNs, nsRestore)
			log.Debugf("restore is %v\n", restore)
			Expect(err).To(BeNil())
			err = kahu.WaitForRestoreCreate(kahuClient, restoreName)
			Expect(err).To(BeNil())
			log.Infof("restore of daemonset is created\n")

			//check if the restored daemonset is up
			daemonSet, err = k8s.GetDaemonset(kubeClient, nsRestore, name)
			log.Debugf("daemonset is %v\n", daemonSet)
			Expect(err).To(BeNil())
			err = k8s.WaitForDaemonsetComplete(kubeClient, nsRestore, name)
			Expect(err).To(BeNil())
			log.Infof("daemonset restored  is up\n")

			//Delete the restore
			err = kahu.DeleteRestore(kahuClient, restoreName)
			Expect(err).To(BeNil())
			err = kahu.WaitForRestoreDelete(kahuClient, restoreName)
			Expect(err).To(BeNil())
			log.Infof("restore of %v is deleted\n", name)

			//Delete the backup
			err = kahu.DeleteBackup(kahuClient, backupName)
			Expect(err).To(BeNil())
			err = kahu.WaitForBackupDelete(kahuClient, backupName)
			Expect(err).To(BeNil())
			log.Infof("backup of %v is deleted\n", name)
		})
	})
})
