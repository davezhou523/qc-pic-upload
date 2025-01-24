package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/syyongx/php2go"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"qc/common/code"
	"qc/common/helper"
	"qc/qc/internal/svc"
	"qc/qc/ossClient"
	"regexp"
	"strings"
	"time"
)

// 测试
//const url = "https://www.51vj.cn/form-pro/api/data-receive?appid=58&corpid=202009232&key=JfLsrodD940njH_Cvu0b6ZFTAH6IzBVWtT8fkJIWVtg="

// 正式
const url = "https://www.51vj.cn/form-pro/api/data-receive?appid=58&corpid=202009232&key=OKNJTvxgXlP-3Aox6REtFoaDppYPAH0ThMdtKY47oj4="

// 微加表单
type WeiJiaFormStruct struct {
	Sno              string   `json:"sno"`
	Xray             []string `json:"xray"`
	TopMarking       []string `json:"topMarking"`
	UndersideMarking []string `json:"undersideMarking"`
	Leads            []string `json:"leads"`
	SideView         []string `json:"sideView"`
}

type XrayStruct struct {
	Sno  string   `json:"sno"`
	Xray []string `json:"xray"`
}

type TopMarkingStruct struct {
	Sno        string   `json:"sno"`
	TopMarking []string `json:"topMarking"`
}
type UndersideMarkingStruct struct {
	Sno              string   `json:"sno"`
	UndersideMarking []string `json:"undersideMarking"`
}
type LeadsStruct struct {
	Sno   string   `json:"sno"`
	Leads []string `json:"leads"`
}
type SideViewStruct struct {
	Sno      string   `json:"sno"`
	SideView []string `json:"sideView"`
}

type UploadDirStruct struct {
	OssPathFile string //上传到OSS目录路径
	Sno         string //流水号目录名
	Category    string //上传目录分类名称
}

// 上传目录分类名称
var FormPicCategory = []string{"leads", "sideView", "topMarking", "undersideMarking", "xRAY"}

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) DealUpload() {
	//l.backupDir("202302070007CG")
	//autoUploadDir := l.svcCtx.Config.InboxDir
	////上传云表单
	//l.uploadFormFromOss(autoUploadDir)

	publicDateDir := l.initPublicDir()
	l.listenPublicDir(publicDateDir)
}
func (l *UploadLogic) CheckAndCreateDir(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			fmt.Println(err.Error())
			l.Logger.Errorf("MkdirAll:%v\n", err.Error())
			return
		}
	}
}
func (l *UploadLogic) initPublicDir() string {
	publicDateDir := l.svcCtx.Config.PublicDir + "/" + time.Now().Format("2006-1-2")
	l.CheckAndCreateDir(publicDateDir)
	fmt.Println("从共享目录(" + publicDateDir + ")以流水号为基准将图片归类QC目录结构,再存入Inbox目标进行自动上传到云表单")
	fmt.Println()
	return publicDateDir

}

func (l *UploadLogic) listenPublicDir(publicDateDir string) {

	go func() {

		ticker := time.NewTicker(time.Second * 7)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				//fmt.Printf("共享目录 %s goroutines:%d,curr_time:%s\n", publicDateDir, runtime.NumGoroutine(), time.Now().Format("2006-01-02 15:04:05"))
				err := l.publicDirFileSort(publicDateDir)
				if err != nil {
					fmt.Println(err)
					return
				} else {
					_ = l.deleteEmptypublicDir(publicDateDir)
					autoUploadDir := l.svcCtx.Config.InboxDir
					err := l.uploadFileToOss(autoUploadDir)
					if err == nil {
						//上传云表单
						err := l.uploadFormFromOss(autoUploadDir)
						if err != nil {
							fmt.Println(err)
							return
						}
					}

				}

			}
		}
	}()
}

// 公共目录Z:\QC\2023-2-11文件归类
func (l *UploadLogic) publicDirFileSort(publicDateDir string) error {
	err := filepath.Walk(publicDateDir, func(pathFile string, info fs.FileInfo, err error) error {
		if info.IsDir() == false && info.Size() > 0 {
			fileSlice := strings.Split(info.Name(), "_")
			picCategoryDir := l.svcCtx.Config.InboxDir
			for k, v := range fileSlice {
				if k == 0 {
					if v == "Thumbs.db" {
						//err = os.Remove(pathFile)
						//fmt.Println("删除:" + pathFile)
						//fmt.Println(err)
						continue
					}
					reg, _ := regexp.MatchString(`^2.*CG$`, v)
					if !reg {
						msg := "流水号不正确，文件名以流水号开头;文件名规则：流水号_型号_qc图片分类名_\n比如：202302010001CG_TM4C1230H6PMI_leads_"
						fmt.Println(msg)
						fmt.Println("当前流水号:", v, "不正确")
						l.Logger.Errorf("%v,当前流水号:%v\n", msg, v)
						continue
					}
					picCategoryDir = picCategoryDir + "/" + v

				} else if k == 2 {
					if !php2go.InArray(v, FormPicCategory) {
						msg := "qc图片分类名不正确，qc图片分类名包含:" + strings.Join(FormPicCategory, ",") + ";文件名规则：流水号_型号_qc图片分类名_\n,比如：202302010001CG_TM4C1230H6PMI_leads_"

						fmt.Println(msg)
						fmt.Println("qc图片分类名:", v)
						l.Logger.Errorf("%v,qc图片分类名:%v", msg, v)
						continue
					}
					picCategoryDir = picCategoryDir + "/" + v
					//fmt.Println("tempdir目录创建流水号目录", picCategoryDir)
					l.CheckAndCreateDir(picCategoryDir)
					destFile := picCategoryDir + "/" + info.Name()
					//文件归类移动到移动inboxDir
					err := helper.MoveFile(pathFile, destFile)
					if err != nil {
						fmt.Println(err)
						return err
					}
					err = os.Remove(pathFile)
					continue
				}

			}
		}

		return err
	})
	return err
}

func (l *UploadLogic) deleteEmptypublicDir(publicDateDir string) error {
	rootDir := time.Now().Format("2006-1-2")
	err := filepath.Walk(publicDateDir, func(pathFile string, info fs.FileInfo, err error) error {
		//fmt.Println(pathFile, info.IsDir(), info.Size(), info.Name())
		if info.IsDir() && rootDir != info.Name() && (info.Size() == 4096 || info.Size() == 0) {
			err := os.RemoveAll(pathFile)
			if err != nil {
				return err
			}
			return filepath.SkipDir
		}
		return err
	})
	return err
}

func (l *UploadLogic) initInboxDir() {
	for _, dir := range l.svcCtx.Config.ListenDir {
		l.CheckAndCreateDir(dir)
	}
	fmt.Println("QC目录结构：202211290004CG：流水号\n\t\t├─ leads：引脚\n\t\t├─ sideView：芯片侧面\n\t\t├─ topMarking：芯片正面\n\t\t├─ undersideMarking：芯片背面\n\t\t└─ xRAY：光射线\n")
	fmt.Println("需要上传的图片目录放入此目录：" + l.svcCtx.Config.InboxDir)
	fmt.Println("上传成功的图片会自动存入此目录：" + l.svcCtx.Config.BackupDir)
	fmt.Println()

}

// 上传oss
func (l *UploadLogic) uploadFileToOss(autoUploadDir string) error {
	client := ossClient.NewOssClient(l.ctx, l.svcCtx)
	//leads", "sideView", "topMarking", "undersideMarking", "xRAY"

	err := filepath.Walk(autoUploadDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() == false && info.Size() > 0 {
			uploadDirObj := GetOssUploadPathFileByLocalPath(path)
			//当有更新时，那要删除已经存在的图片
			client.DeleteObjectList(uploadDirObj.Sno + "/" + uploadDirObj.Category)
			return filepath.SkipDir
		}
		return err
	})
	err = filepath.Walk(autoUploadDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() == false && info.Size() > 0 {
			uploadDirObj := GetOssUploadPathFileByLocalPath(path)

			client.PutObject(path, uploadDirObj.OssPathFile)

		}
		return err
	})
	return err
}

// 上传云表单
func (l *UploadLogic) uploadFormFromOss(autoUploadDir string) error {
	ossClient := ossClient.NewOssClient(l.ctx, l.svcCtx)
	var formMapStruct map[string]WeiJiaFormStruct
	formMapStruct = make(map[string]WeiJiaFormStruct, 100)
	var lastSno string //上次流水号
	var formStruct = new(WeiJiaFormStruct)
	err := filepath.Walk(autoUploadDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() == false && info.Size() > 0 {
			uploadDirObj := GetOssUploadPathFileByLocalPath(path)
			sno := uploadDirObj.Sno
			if len(lastSno) > 0 && sno != lastSno {
				//实现不同的流水号图片分开存储地址
				formStruct = new(WeiJiaFormStruct)
			}
			lastSno = sno
			imageUrl := l.svcCtx.Config.Oss.BucketDns + "/" + uploadDirObj.OssPathFile
			//fmt.Println("imageUrl:" + imageUrl)
			switch uploadDirObj.Category {
			case "leads":
				formStruct.Leads = append(formStruct.Leads, imageUrl)
			case "sideView":
				formStruct.SideView = append(formStruct.SideView, imageUrl)
			case "topMarking":
				formStruct.TopMarking = append(formStruct.TopMarking, imageUrl)
			case "undersideMarking":
				formStruct.UndersideMarking = append(formStruct.UndersideMarking, imageUrl)
			case "xRAY":
				formStruct.Xray = append(formStruct.Xray, imageUrl)
			default:
				//log := "表单目录名错误:" + uploadDirObj.Category
				//l.Logger.Errorf("%v\n", log)
				//return fmt.Errorf("%v\n", log)
			}
			formStruct.Sno = sno
			formMapStruct[sno] = *formStruct
			if err != nil {
				return err
			}

		}
		return err
	})

	if err != nil {
		l.Logger.Errorf("%v\n", err.Error())
		return err
	}
	for _, saveFormStruct := range formMapStruct {
		//fmt.Println("saveFormStruct:", saveFormStruct)
		if len(saveFormStruct.Sno) <= 0 {
			continue
		}
		sno := saveFormStruct.Sno

		if len(saveFormStruct.Leads) == 0 {
			url, err := ossClient.GetObjectList(sno + "/leads")
			if len(url) > 0 {
				saveFormStruct.Leads = url
			}
			if err != nil {
				fmt.Println(err)
				l.Logger.Errorf("GetObjectList:%v\n", err.Error())
				continue
			}
		}

		if len(saveFormStruct.SideView) == 0 {
			url, err := ossClient.GetObjectList(sno + "/sideView")
			if len(url) > 0 {
				saveFormStruct.SideView = url
			}
			if err != nil {
				fmt.Println(err)
				l.Logger.Errorf("GetObjectList:%v\n", err.Error())
				continue
			}
		}
		if len(saveFormStruct.TopMarking) == 0 {
			url, err := ossClient.GetObjectList(sno + "/topMarking")
			if len(url) > 0 {
				saveFormStruct.TopMarking = url
			}
			if err != nil {
				fmt.Println(err)
				l.Logger.Errorf("GetObjectList:%v\n", err.Error())
				continue
			}
		}
		if len(saveFormStruct.UndersideMarking) == 0 {
			url, err := ossClient.GetObjectList(sno + "/undersideMarking")
			if len(url) > 0 {
				saveFormStruct.UndersideMarking = url
			}
			if err != nil {
				fmt.Println(err)
				l.Logger.Errorf("GetObjectList:%v\n", err.Error())
				continue
			}
		}
		if len(saveFormStruct.Xray) == 0 {
			url, err := ossClient.GetObjectList(sno + "/xRAY")
			if len(url) > 0 {
				saveFormStruct.Xray = url
			}
			if err != nil {
				fmt.Println(err)
				l.Logger.Errorf("GetObjectList:%v\n", err.Error())
				continue
			}
		}
		formByte, err := json.Marshal(saveFormStruct)
		if err != nil {
			fmt.Println(err.Error())
			l.Logger.Errorf("formJson:%v\n", err.Error())
			return err
		}
		err = l.postRequst(formByte)
		if err != nil {
			fmt.Println(err.Error())
			l.Logger.Errorf("postRequst:%v\n", err.Error())
			return err
		}
		fmt.Println(saveFormStruct.Sno + " 上传表单成功")
		err = l.backupDir(saveFormStruct.Sno, autoUploadDir)
		if err != nil {
			return err
		}

	}

	return nil
}

func (l *UploadLogic) postRequst(formByte []byte) error {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(formByte))
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if err != nil {
		fmt.Println(err.Error())
		l.Logger.Errorf("httpRequst:%v\n", err.Error())
		return err
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		l.Logger.Errorf("httpRequst:%v\n", err.Error())
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("request response:%v\n", resp.StatusCode)
		l.Logger.Errorf("request response:%v\n", resp.StatusCode)
		return fmt.Errorf("request response:%v\n", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)

	res := code.HttpResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}
	if res.Code == 1 {
		return nil
	} else {
		return fmt.Errorf("%v", res.Msg)
	}
}

// 备份目录
func (l *UploadLogic) backupDir(sno string, autoUploadDir string) error {
	snoDir := autoUploadDir + "/" + sno
	err := filepath.Walk(snoDir, func(pathFile string, info fs.FileInfo, err error) error {
		if info.IsDir() && info.Name() != sno {
			backupDir := l.svcCtx.Config.BackupDir + "/" + sno + "/" + info.Name()
			l.CheckAndCreateDir(backupDir)
		}
		if info.IsDir() == false && info.Size() > 0 {
			pathSlice := strings.Split(pathFile, "\\")
			picCategoryName := pathSlice[len(pathSlice)-2]
			backupFile := l.svcCtx.Config.BackupDir + "/" + sno + "/" + picCategoryName + "/" + info.Name()
			//fmt.Println(backupFile)
			err := helper.MoveFile(pathFile, backupFile)
			if err != nil {
				fmt.Println(err)
				return err
			}
			//return filepath.SkipDir
		}

		return err
	})
	if err == nil {
		err = os.RemoveAll(snoDir)
		if err != nil {
			return err
		}
	}
	return err
}

// OssPathFile string //上传到OSS目录路径
// Sno         string //流水号目录名
// Category    string //上传目录分类名称
func GetOssUploadPathFileByLocalPath(localPathFile string) UploadDirStruct {
	var uploadDir UploadDirStruct
	pathSlice := strings.Split(localPathFile, "\\")
	var findkey = 0
	var ossPathFile string
	for k, v := range pathSlice {
		reg, _ := regexp.MatchString(`^2.*CG$`, v)
		if reg {
			findkey = k
			ossPathFile = v
			uploadDir.Sno = v
		}
		if k > findkey {
			ossPathFile = ossPathFile + "/" + v
		}
		if php2go.InArray(v, FormPicCategory) {
			////oss 上传目录名
			uploadDir.Category = v
		}
	}
	//oss 上传路径
	uploadDir.OssPathFile = ossPathFile
	return uploadDir
}
