package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	ListenDir map[string]string
	InboxDir  string
	BackupDir string
	FailedDir string
	PublicDir string
	TempDir   string
	Oss       struct {
		Endpoint        string
		BucketName      string
		BucketDns       string
		AccessKeyId     string
		AccessKeySecret string
	}
}
