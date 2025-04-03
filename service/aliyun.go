// This file is auto-generated, don't edit it. Thanks.
package service

import (
	"encoding/json"
	"fmt"
	"generate-steam-ai-excel/code"
	"generate-steam-ai-excel/global"
	"generate-steam-ai-excel/util"
	bailian20231229 "github.com/alibabacloud-go/bailian-20231229/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliUtil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Description:
//
// 使用AK&SK初始化账号Client
//
// @return Client
//
// @throws Exception
func CreateClient() {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/bailian
	config.Endpoint = tea.String("bailian.cn-beijing.aliyuncs.com")
	_result, _err := bailian20231229.NewClient(config)
	if _err != nil {
		global.Logger.Error("创建客户端失败", code.ERROR, _err)
		os.Exit(1)
	}
	global.ALiYunClient = _result
	return
}

func ListBaiLianFile() (x []*bailian20231229.ListFileResponseBodyDataFileList, _err error) {
	result := make([]*bailian20231229.ListFileResponseBodyDataFileList, 0)
	//client, _err := CreateClient()
	//if _err != nil {
	//	return result, _err
	//}

	listFileRequest := &bailian20231229.ListFileRequest{
		CategoryId: global.ALiYunBaiLianCateID,
	}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		resp, _err := global.ALiYunClient.ListFileWithOptions(global.ALiYunBaiLianWorkspaceID, listFileRequest, global.AliYunHeaders, global.ALiYunRuntime)
		if _err != nil {
			return _err
		}
		result = resp.Body.Data.FileList

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 此处仅做打印展示，请谨慎对待异常处理，在工程项目中切勿直接忽略异常。
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}
		_, _err = aliUtil.AssertAsString(error.Message)
		if _err != nil {
			return result, _err
		}
	}
	return result, _err
}

func ApplyFileUploadLease(fileName string) (*bailian20231229.ApplyFileUploadLeaseResponseBodyData, error) {
	fileMd5, err := util.GetFileMD5(fileName)
	if err != nil {
		global.Logger.Error("获取文件MD5失败", code.ERROR, err)
		return nil, err
	}
	fileSize, err := util.GetFileSize(fileName)
	if err != nil {
		global.Logger.Error("获取文件大小失败", code.ERROR, err)
		return nil, err
	}
	fmt.Println(fileSize)
	applyFileUploadLeaseRequest := &bailian20231229.ApplyFileUploadLeaseRequest{
		FileName:    tea.String(fileName),
		SizeInBytes: tea.String(strconv.Itoa(int(fileSize))),
		Md5:         tea.String(fileMd5),
	}

	resp, _err := global.ALiYunClient.ApplyFileUploadLeaseWithOptions(global.ALiYunBaiLianCateID,
		global.ALiYunBaiLianWorkspaceID, applyFileUploadLeaseRequest, global.AliYunHeaders, global.ALiYunRuntime)

	if _err != nil {
		global.Logger.Error("申请文件上传失败", code.ERROR, _err)
		return nil, _err
	}
	fmt.Println(resp.Body.String())
	return resp.Body.Data, nil
}

func UploadFileToALiYunBaiLian(fileName string) (string, error) {
	lease, err := ApplyFileUploadLease(fileName)
	if err != nil {
		return "", err
	}

	x := lease.Param.Headers.(map[string]any)
	b, _ := os.ReadFile(fileName)
	req, _ := http.NewRequest(http.MethodPut, *lease.Param.Url, tea.ToReader(b))
	for k, v := range x {
		_v := v.(string)
		req.Header.Add(k, _v)
	}

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		global.Logger.Error("上传文件失败", code.ERROR, err)
		return "", err
	}
	if resp.StatusCode != 200 {
		respBytes, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		errStr := util.Bytes2str(respBytes)
		global.Logger.Error("上传文件失败", code.ERROR, errStr)
		return "", fmt.Errorf(errStr)
	}
	addFileRequest := &bailian20231229.AddFileRequest{
		Parser:     tea.String("DASHSCOPE_DOCMIND"),
		CategoryId: global.ALiYunBaiLianCateID,
		LeaseId:    lease.FileUploadLeaseId,
	}

	res, _err := global.ALiYunClient.AddFileWithOptions(global.ALiYunBaiLianWorkspaceID, addFileRequest,
		global.AliYunHeaders, global.ALiYunRuntime)
	if _err != nil {
		global.Logger.Error("上传文件失败", code.ERROR, _err)
		return "", _err
	}

	return *res.Body.Data.FileId, nil
}

func DescribeBaiLianFile(fileID string) (*bailian20231229.DescribeFileResponseBodyData, error) {
	s, _err := global.ALiYunClient.DescribeFileWithOptions(global.ALiYunBaiLianWorkspaceID, tea.String(fileID),
		global.AliYunHeaders, global.ALiYunRuntime)
	if _err != nil {
		global.Logger.Error("获取文件信息失败", code.ERROR, _err)
		return nil, _err
	}
	return s.Body.Data, nil
}

func SubmitIndexAddDocumentsJob(documentIds []*string) error {
	submitIndexAddDocumentsJobRequest := &bailian20231229.SubmitIndexAddDocumentsJobRequest{
		IndexId:     global.ALiYunIndexID,
		SourceType:  tea.String("DATA_CENTER_FILE"),
		DocumentIds: documentIds,
		CategoryIds: []*string{global.ALiYunBaiLianCateID},
	}
	_, _err := global.ALiYunClient.SubmitIndexAddDocumentsJobWithOptions(global.ALiYunBaiLianWorkspaceID,
		submitIndexAddDocumentsJobRequest, global.AliYunHeaders, global.ALiYunRuntime)
	if _err != nil {
		global.Logger.Error("提交索引添加文档任务失败", code.ERROR, _err)
		return _err
	}
	return nil
}

func ListIndexDocuments() []*bailian20231229.ListIndexDocumentsResponseBodyDataDocuments {
	listIndexDocumentsRequest := &bailian20231229.ListIndexDocumentsRequest{
		IndexId: global.ALiYunIndexID,
	}
	resp, _err := global.ALiYunClient.ListIndexDocumentsWithOptions(global.ALiYunBaiLianWorkspaceID,
		listIndexDocumentsRequest, global.AliYunHeaders, global.ALiYunRuntime)
	if _err != nil {
		global.Logger.Error("查询索引文档失败", code.ERROR, _err)
		return nil
	}
	return resp.Body.Data.Documents
}

func DeleteIndexDocument(docs []*string) error {
	deleteIndexDocumentRequest := &bailian20231229.DeleteIndexDocumentRequest{
		IndexId:     global.ALiYunIndexID,
		DocumentIds: docs,
	}
	_, _err := global.ALiYunClient.DeleteIndexDocumentWithOptions(global.ALiYunBaiLianWorkspaceID,
		deleteIndexDocumentRequest, global.AliYunHeaders, global.ALiYunRuntime)
	if _err != nil {
		return _err
	}
	return nil
}
