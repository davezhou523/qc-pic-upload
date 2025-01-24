package main

import (
	"flag"
	"fmt"
	"qc/qcClient/internal/config"
	"qc/qcClient/internal/handler"
	"qc/qcClient/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/qcClient.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	c.Oss.Endpoint = "https://oss-cn-shenzhen.aliyuncs.com"
	c.Oss.BucketName = "qcClient-pic"
	c.Oss.BucketDns = "https://qc-pic.oss-cn-shenzhen.aliyuncs.com"
	c.Oss.AccessKeyId = "LTAI5tDuxMaSsupd49ErNEUo"
	c.Oss.AccessKeySecret = ""

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.UploadHandler(ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
