package ossClient

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"qc/qc/internal/svc"
)

type OssClient struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	client *oss.Client
}

func NewOssClient(ctx context.Context, svcCtx *svc.ServiceContext) OssClient {
	client, err := oss.New(svcCtx.Config.Oss.Endpoint, svcCtx.Config.Oss.AccessKeyId, svcCtx.Config.Oss.AccessKeySecret)
	if err != nil {
		log := "连接阿里云oss失败:"
		fmt.Printf("%v,%v\n", log, err.Error())
		logx.Errorf("%v,%v\n", log, err.Error())
		return OssClient{
			Logger: logx.WithContext(ctx),
			ctx:    ctx,
			svcCtx: svcCtx,
			client: nil,
		}
	}
	return OssClient{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		client: client,
	}

}
func (l *OssClient) PutObject(localFile string, fileName string) error {
	bucket, err := l.client.Bucket(l.svcCtx.Config.Oss.BucketName)
	if err != nil {
		log := "阿里云oss上传文件Bucket失败:"
		logx.Errorf("%v,%v\n", log, err.Error())
		return fmt.Errorf("%v,%v\n", log, err.Error())
	}
	fd, err := os.Open(localFile)
	if err != nil {
		log := "打开文件失败:"
		logx.Errorf("%v,%v\n", log, err.Error())
		return fmt.Errorf("%v,%v\n", log, err.Error())
	}
	defer fd.Close()
	err = bucket.PutObject(fileName, fd)
	if err != nil {
		log := "阿里云oss上传文件失败:"
		logx.Errorf("%v,%v\n", log, err.Error())
		return fmt.Errorf("%v,%v\n", log, err.Error())
	}
	return nil
}
func (l *OssClient) GetObjectList(matchData string) ([]string, error) {
	var res []string
	bucket, err := l.client.Bucket(l.svcCtx.Config.Oss.BucketName)
	if err != nil {
		log := "阿里云oss获取文件Bucket失败:"
		logx.Errorf("%v,%v\n", log, err.Error())
		return res, fmt.Errorf("%v,%v\n", log, err.Error())
	}
	lsRes, err := bucket.ListObjects(oss.Prefix(matchData))
	if err != nil {
		log := "阿里云oss获取文件列表失败:"
		logx.Errorf("%v,%v\n", log, err.Error())
		return res, fmt.Errorf("%v,%v\n", log, err.Error())
	}

	for _, object := range lsRes.Objects {
		url := l.svcCtx.Config.Oss.BucketDns + "/" + object.Key
		res = append(res, url)
	}
	return res, nil
}
func (l *OssClient) DeleteObjectList(matchData string) ([]string, error) {
	var res []string
	bucket, err := l.client.Bucket(l.svcCtx.Config.Oss.BucketName)
	if err != nil {
		log := "阿里云oss获取文件Bucket失败:"
		logx.Errorf("%v,%v\n", log, err.Error())
		return res, fmt.Errorf("%v,%v\n", log, err.Error())
	}
	lsRes, err := bucket.ListObjects(oss.Prefix(matchData))
	if err != nil {
		log := "阿里云oss获取文件列表失败:"
		logx.Errorf("%v,%v\n", log, err.Error())
		return res, fmt.Errorf("%v,%v\n", log, err.Error())
	}

	for _, object := range lsRes.Objects {
		//url := l.svcCtx.Config.Oss.BucketDns + "/" + object.Key
		err = bucket.DeleteObject(object.Key)
	}
	return res, nil
}
