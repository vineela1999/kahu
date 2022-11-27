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
var _ = Describe("statefulsetBackup", Label("statefulset"), func() {
	Context("Create backup of statefulset and restore", func() {
		It("statefulset", func() {
			kubeClient, kahuClient := kahu.Clients()
			//Create statefulset to test
			ns := kahu.BackupNameSpace
			labels := make(map[string]string)
			image := "nginx"
			replicas := 2
			labels["statefulset"] = "statefulset"

			UUIDgen, err := uuid.NewRandom()
			Expect(err).To(BeNil())
			name := "statefulset" + "-" + UUIDgen.String()
			statefulSet, err := k8s.NewStatefulset(name, ns, int32(replicas), labels, image)
			log.Infof("statefulset:%v\n", statefulSet)
			statefulSet1, err := k8s.CreateStatefulset(kubeClient, ns, statefulSet)
			log.Debugf("statefulset:%v\n", statefulSet1)
			Expect(err).To(BeNil())
			err = k8s.WaitForStatefulsetComplete(kubeClient, ns, name)
			Expect(err).To(BeNil())

			//create backup for the statefulset
			backupName := "backup" + "statefulset" + "-" + UUIDgen.String()
			includeNs := kahu.BackupNameSpace
			resourceType := "StatefulSet"
			_, err = kahu.CreateBackup(kahuClient, backupName, includeNs, resourceType)
			Expect(err).To(BeNil())
			err = kahu.WaitForBackupCreate(kahuClient, backupName)
			Expect(err).To(BeNil())
			log.Infof("backup of statefulset is done\n")

			// create restore for the backup
			restoreName := "restore" + "statefulset" + "-" + UUIDgen.String()
			nsRestore := kahu.RestoreNameSpace
			restore, err := kahu.CreateRestore(kahuClient, restoreName, backupName, includeNs, nsRestore)
			log.Debugf("restore is %v\n", restore)
			Expect(err).To(BeNil())
			err = kahu.WaitForRestoreCreate(kahuClient, restoreName)
			Expect(err).To(BeNil())
			log.Infof("restore of statefulset is created\n")

			//check if the restored deployment is up
			statefulSet, err = k8s.GetStatefulset(kubeClient, nsRestore, name)
			log.Debugf("statefulset is %v\n", statefulSet)
			Expect(err).To(BeNil())
			err = k8s.WaitForStatefulsetComplete(kubeClient, nsRestore, name)
			Expect(err).To(BeNil())
			log.Infof("statefulset restored is up\n")

			//Delete the. restore
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
			log.Infof("backup of  %v is deleted\n", name)
		})
	})
})
