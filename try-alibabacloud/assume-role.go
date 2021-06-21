// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"fmt"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/client"
	"github.com/alibabacloud-go/tea/tea"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *sts20150401.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("sts.cn-beijing.aliyuncs.com")
	_result = &sts20150401.Client{}
	_result, _err = sts20150401.NewClient(config)
	return _result, _err
}

func _main(args []*string) (_err error) {
	var (
		ak = "xxx"
		sk = "xxx"
	)
	client, _err := CreateClient(&ak, &sk)
	if _err != nil {
		return _err
	}

	assumeRoleRequest := &sts20150401.AssumeRoleRequest{
		DurationSeconds: tea.Int64(3600),
		RoleArn:         tea.String("acs:ram::xxx-id-xxx:role/poc"),
		RoleSessionName: tea.String("poc"),
	}
	// 复制代码运行请自行打印 API 的返回值
	res, _err := client.AssumeRole(assumeRoleRequest)
	if _err != nil {
		return _err
	}
	fmt.Println(*res.Body.Credentials.AccessKeyId)
	fmt.Println(*res.Body.Credentials.AccessKeySecret)
	fmt.Println(*res.Body.Credentials.SecurityToken)
	return _err
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}
