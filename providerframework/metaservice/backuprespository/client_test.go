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

package backuprespository

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	pb "github.com/soda-cdm/kahu/providers/lib/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

var flag int = 0

type ClientTestSuite struct {
	suite.Suite
}

type FakeMetaBackup_DownloadClient struct {
	mock.Mock
	grpc.ClientStream
}

type FakeMetaBackup_UploadClient struct {
	mock.Mock
	grpc.ClientStream
}

func (f *FakeMetaBackup_DownloadClient) Recv() (*pb.DownloadResponse, error) {
	args := f.Called()
	return args.Get(0).(*pb.DownloadResponse), args.Error(1)
}

func (f *FakeMetaBackup_UploadClient) Send(ur *pb.UploadRequest) error {
	/*
		flag++
		fmt.Printf("flag value:%v\n", flag)
		if flag == 1 && (reflect.TypeOf(ur.Data) != reflect.TypeOf(&pb.UploadRequest_Info{})) {
			return fmt.Errorf("Data of type 'UploadRequest_Info' is expected from UploadRequest\n")
		}

		if flag == 2 && (reflect.TypeOf(ur.Data) != reflect.TypeOf(&pb.UploadRequest_ChunkData{})) {
			return fmt.Errorf("Data of type 'UploadRequest_ChunkData' is expected from UploadRequest\n")
		}
	*/
	args := f.Called(ur)
	return args.Error(0)
}

func (f *FakeMetaBackup_UploadClient) CloseAndRecv() (*pb.Empty, error) {
	args := f.Called()
	return args.Get(0).(*pb.Empty), args.Error(1)
}

type FakeMetaBackupClient struct {
	mock.Mock
}

func (f *FakeMetaBackupClient) Upload(ctx context.Context, opts ...grpc.CallOption) (pb.MetaBackup_UploadClient, error) {
	args := f.Called(ctx, opts)
	return args.Get(0).(pb.MetaBackup_UploadClient), args.Error(1)
}

func (f *FakeMetaBackupClient) ObjectExists(ctx context.Context, in *pb.ObjectExistsRequest, opts ...grpc.CallOption) (*pb.ObjectExistsResponse, error) {
	args := f.Called(ctx, in, opts)
	return args.Get(0).(*pb.ObjectExistsResponse), args.Error(1)
}

func (f *FakeMetaBackupClient) Download(ctx context.Context, in *pb.DownloadRequest, opts ...grpc.CallOption) (pb.MetaBackup_DownloadClient, error) {
	args := f.Called(ctx, in, opts)
	return args.Get(0).(pb.MetaBackup_DownloadClient), args.Error(1)

}

func (f *FakeMetaBackupClient) Delete(ctx context.Context, in *pb.DeleteRequest, opts ...grpc.CallOption) (*pb.Empty, error) {
	args := f.Called(ctx, in, opts)
	return args.Get(0).(*pb.Empty), args.Error(1)
}

func (suite *ClientTestSuite) BeforeTest() {}

func (suite *ClientTestSuite) TestUploadInvalidPath() {
	repo, _, _ := NewBackupRepository("/tmp/nfs_test.sock")
	filePath := "fakeFilePath/fakeFile"
	err := repo.Upload(filePath)
	assert.NotNil(suite.T(), err)
}

func (suite *ClientTestSuite) TestUploadCLientUploadErr() {
	client := new(FakeMetaBackupClient)
	fakeRespoClient := new(FakeMetaBackup_UploadClient)
	expErr := fmt.Errorf("this is error from mocked upload")
	client.On("Upload", mock.Anything, mock.Anything).Return(fakeRespoClient, expErr)
	repo := &backupRepository{
		backupRepositoryAddress: "127.0.0.1:8181",
		client:                  client,
	}
	filePath := "fakeFile"
	file, err := os.Create(filePath)
	assert.Nil(suite.T(), err)
	defer file.Close()
	str := "this is sample data"
	data := []byte(str)
	err = os.WriteFile(filePath, data, 0777)
	err = repo.Upload(filePath)
	assert.Equal(suite.T(), err, expErr)
}

func (suite *ClientTestSuite) TestUploadCLientSendErr() {
	client := new(FakeMetaBackupClient)
	fakeRespoClient := new(FakeMetaBackup_UploadClient)

	client.On("Upload", mock.Anything, mock.Anything).Return(fakeRespoClient, nil)
	repo := &backupRepository{
		backupRepositoryAddress: "127.0.0.1:8181",
		client:                  client,
	}
	expErr := fmt.Errorf("this is error from mocked send func")
	fakeRespoClient.On("Send", mock.Anything).Return(expErr)
	//fakeRespoClient.On("CloseAndRecv", mock.Anything).Return(&pb.Empty{}, nil)
	filePath := "fakeFile"
	file, err := os.Create(filePath)
	assert.Nil(suite.T(), err)
	defer file.Close()
	str := "this is sample data"
	data := []byte(str)
	err = os.WriteFile(filePath, data, 0777)
	err = repo.Upload(filePath)
	assert.Equal(suite.T(), err, expErr)
}

func (suite *ClientTestSuite) TestDownloadFail() {
	repo, _, _ := NewBackupRepository("/tmp/nfs_test.sock")
	fileId := "fakeFilePath/fakeFile"

	attributes := map[string]string{"key": "value"}
	_, err := repo.Download(fileId, attributes)
	assert.NotNil(suite.T(), err)
}

func (suite *ClientTestSuite) TestDownloadCLientDownloadErr() {
	client := new(FakeMetaBackupClient)
	fakeRespoClient := new(FakeMetaBackup_DownloadClient)
	expErr := fmt.Errorf("this is error from mocked download func")
	client.On("Download", mock.Anything, mock.Anything, mock.Anything).Return(fakeRespoClient, expErr)
	repo := &backupRepository{
		backupRepositoryAddress: "127.0.0.1:8181",
		client:                  client,
	}
	fakeRespoClient.On("Recv", mock.Anything).Return(&pb.DownloadResponse{}, nil)
	fakeRespoClient.On("CloseSend").Return(nil)
	fileID := "fakeFile"
	attributes := make(map[string]string)
	//out := "fakeFile"
	_, err := repo.Download(fileID, attributes)
	assert.Equal(suite.T(), err, expErr)
}

func (suite *ClientTestSuite) TestDownloadCLientRecvErr() {
	client := new(FakeMetaBackupClient)
	fakeRespoClient := new(FakeMetaBackup_DownloadClient)
	expErr := fmt.Errorf("this is error from mocked Recv func")
	client.On("Download", mock.Anything, mock.Anything, mock.Anything).Return(fakeRespoClient, nil)
	repo := &backupRepository{
		backupRepositoryAddress: "127.0.0.1:8181",
		client:                  client,
	}
	fakeRespoClient.On("Recv", mock.Anything).Return(&pb.DownloadResponse{}, expErr)
	fakeRespoClient.On("CloseSend").Return(nil)
	fileID := "fakeFile"
	attributes := make(map[string]string)
	//out := "fakeFile"
	_, err := repo.Download(fileID, attributes)
	assert.Equal(suite.T(), err, expErr)
}

func (suite *ClientTestSuite) TestUploadValid() {
	client := new(FakeMetaBackupClient)
	fakeRespoClient := new(FakeMetaBackup_UploadClient)
	client.On("Upload", mock.Anything, mock.Anything).Return(fakeRespoClient, nil)
	repo := &backupRepository{
		backupRepositoryAddress: "127.0.0.1:8181",
		client:                  client,
	}
	fakeRespoClient.On("Send", mock.Anything).Return(nil)
	fakeRespoClient.On("CloseAndRecv", mock.Anything).Return(&pb.Empty{}, nil)
	filePath := "fakeFile"
	file, err := os.Create(filePath)
	assert.Nil(suite.T(), err)
	defer file.Close()
	str := "this is sample data"
	data := []byte(str)
	err = os.WriteFile(filePath, data, 0777)
	err = repo.Upload(filePath)
	assert.Nil(suite.T(), err)

}

func retVal() (*pb.DownloadResponse, error) {
	flag++
	fmt.Printf("flag %d :\n", flag)
	downloadRes := &pb.DownloadResponse{
		Data: &pb.DownloadResponse_Info{
			Info: &pb.DownloadResponse_FileInfo{
				FileIdentifier: "fakefile",
			},
		},
	}

	downloadRes1 := &pb.DownloadResponse{
		Data: &pb.DownloadResponse_ChunkData{
			ChunkData: ([]byte)("fakefile"),
		},
	}
	if flag == 1 {
		fmt.Printf("flag x %d :\n", flag)
		return downloadRes, nil
	}

	if flag == 2 {
		fmt.Printf("flag y %d :\n", flag)
		return downloadRes1, io.EOF
	}
	return downloadRes, nil
}
func (suite *ClientTestSuite) TestDownloadValid() {
	client := new(FakeMetaBackupClient)
	fakeRespoClient := new(FakeMetaBackup_DownloadClient)
	client.On("Download", mock.Anything, mock.Anything, mock.Anything).Return(fakeRespoClient, nil)
	repo := &backupRepository{
		backupRepositoryAddress: "127.0.0.1:8181",
		client:                  client,
	}
	/*
		downloadRes := &pb.DownloadResponse{
			Data: &pb.DownloadResponse_Info{
				Info: &pb.DownloadResponse_FileInfo{
					FileIdentifier: "fakefile",
				},
			},
		}

		downloadRes1 := &pb.DownloadResponse{
			Data: &pb.DownloadResponse_ChunkData{
				ChunkData: ([]byte)("fakefile"),
			},
		}
	*/
	fakeRespoClient.On("Recv").Return(retVal())
	fakeRespoClient.On("CloseSend").Return(nil)
	fileID := "fakeFile"
	attributes := make(map[string]string)
	_, err := repo.Download(fileID, attributes)
	assert.Nil(suite.T(), err)
}

func TestTarTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
