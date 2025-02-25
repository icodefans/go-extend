package cloud

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	cdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v1/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v1/region"
)

// 华为云
type HuaWeiCloud struct {
	AccessKeyId     string `mapstructure:"AccessKeyId"`
	AccessKeySecret string `mapstructure:"AccessKeySecret"`
	EndPoint        string `mapstructure:"EndPoint"`
	ProjectId       string `mapstructure:"ProjectId"`
}

// cdn缓存刷新(文件)
func (huawei *HuaWeiCloud) CdnFileRefresh(url string) error {
	ak := huawei.AccessKeyId
	sk := huawei.AccessKeySecret

	auth := global.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := cdn.NewCdnClient(
		cdn.CdnClientBuilder().
			WithRegion(region.ValueOf("cn-north-1")).
			WithCredential(auth).
			Build())

	request := &model.CreateRefreshTasksRequest{}
	enterpriseProjectIdRequest := huawei.ProjectId
	request.EnterpriseProjectId = &enterpriseProjectIdRequest

	// 增加刷新的url
	var listUrlsRefreshTask = []string{}
	listUrlsRefreshTask = append(listUrlsRefreshTask, url)

	typeRefreshTask := model.GetRefreshTaskRequestBodyTypeEnum().FILE
	refreshTaskbody := &model.RefreshTaskRequestBody{
		Type: &typeRefreshTask,
		Urls: listUrlsRefreshTask,
	}
	request.Body = &model.RefreshTaskRequest{
		RefreshTask: refreshTaskbody,
	}
	response, err := client.CreateRefreshTasks(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
	return err
}

// cdn缓存刷新(目录)
func (huawei *HuaWeiCloud) CdnDirectoryRefresh(url string) error {
	ak := huawei.AccessKeyId
	sk := huawei.AccessKeySecret

	auth := global.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := cdn.NewCdnClient(
		cdn.CdnClientBuilder().
			WithRegion(region.ValueOf("cn-north-1")).
			WithCredential(auth).
			Build())

	request := &model.CreateRefreshTasksRequest{}
	enterpriseProjectIdRequest := huawei.ProjectId
	request.EnterpriseProjectId = &enterpriseProjectIdRequest

	// 增加刷新的url
	var listUrlsRefreshTask = []string{}
	listUrlsRefreshTask = append(listUrlsRefreshTask, url)

	typeRefreshTask := model.GetRefreshTaskRequestBodyTypeEnum().DIRECTORY
	refreshTaskbody := &model.RefreshTaskRequestBody{
		Type: &typeRefreshTask,
		Urls: listUrlsRefreshTask,
	}
	request.Body = &model.RefreshTaskRequest{
		RefreshTask: refreshTaskbody,
	}
	response, err := client.CreateRefreshTasks(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}
	return err
}
