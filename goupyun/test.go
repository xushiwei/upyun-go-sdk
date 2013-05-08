package main

import (
	"fmt"
	"log"
	"os"
	"upyun"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("\n\n\n\t\t\t", err)
		}
	}()

	/// 初始化空间
	u := upyun.NewUpYun("空间名称", "用户名", "密码")
	log.Printf("SDK 版本 %v\n", u.Version())

	/// 设置是否打印调试信息, 当Debug==false时所有文件操作错误调试都将跳过，不会中断当前程序执行
	u.Debug = true

	/// 切换 API 接口的域名
	/// {默认 v0.api.upyun.com 自动识别, v1.api.upyun.com 电信, v2.api.upyun.com 联通, v3.api.upyun.com 移动}
	u.SetApiDomain("v0.api.upyun.com")

	/// 获取空间占用大小
	v, err := u.GetBucketUsage()
	log.Printf("GetBucketUsage: %v %v\n", v, err)

	/// 获取某个目录的空间占用大小
	v, err = u.GetFolderUsage("/")
	log.Printf("GetFolderUsage: %v %v\n", v, err) // 必须斜杠结尾

	/// 上传文件
	// log.Printf("%v\n", u.WriteFile("/test.txt", "test test"))
	// 上传文件时可使用u.WriteFile("/a/b/c/test.txt", "test test", true)进行父级目录的自动创建（最深10级目录）

	/// 设置待上传文件的 Content-MD5 值（如又拍云服务端收到的文件MD5值与用户设置的不一致，将回报 406 Not Acceptable 错误）
	fileMd5, err := upyun.FileMd5("/tmp/test.jpg")
	if err != nil {
		log.Fatalln("file md5 failed: ", err)
	}
	log.Printf("FileMd5: %v\n", fileMd5)
	u.SetContentMD5(fileMd5)

	/// 设置待上传文件的 访问密钥（注意：仅支持图片空！，设置密钥后，无法根据原文件URL直接访问，需带 URL 后面加上 （缩略图间隔标志符+密钥） 进行访问）
	/// 如缩略图间隔标志符为 ! ，密钥为 bac，上传文件路径为 /folder/test.jpg ，那么该图片的对外访问地址为： http://空间域名/folder/test.jpg!bac
	// u.SetFileSecret("bac")

	// 采用数据流模式上传文件（可节省内存）
	fileName := "/tmp/test.jpg"
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(fmt.Sprintf("open file %s failed: %v", fileName, err))
	}
	log.Printf("WriteFile: %v\n", u.WriteFile("/test.jpg", fh, true))
	fh.Close()

	fileName = "/tmp/test.tif"
	fh, err = os.Open(fileName)
	if err != nil {
		log.Fatalln(fmt.Sprintf("open file %s failed: %v", fileName, err))
	}
	log.Printf("WriteFile: %v\n", u.WriteFile("/tmp/test.tif", fh, true))
	fh.Close()

	/// 获取上传后的图片信息（仅图片空间有返回数据）
	log.Printf("x-upyun-width: %v\n", u.GetWritedFileInfo("x-upyun-width"))    // 图片宽度
	log.Printf("x-upyun-height: %v\n", u.GetWritedFileInfo("x-upyun-height"))  // 图片高度
	log.Printf("x-upyun-frames: %v\n", u.GetWritedFileInfo("x-upyun-frames"))  // 图片帧数
	log.Printf("x-upyun-type: %v\n", u.GetWritedFileInfo("x-upyun-file-type")) // 图片类型

	/// 读取文件
	// log.Printf("%v\n", u.ReadFile("/test.txt")
	// 采用数据流模式下载文件（可节省内存）
	fh, err = os.Create("/tmp/test2.jpg")
	if err != nil {
		log.Println(err)
	}
	log.Printf("ReadFile: %v\n", u.ReadFile("/test.jpg", fh))
	fh.Close()

	/// 获取文件信息 return map("type": file | folder, "size": file size, "date": unix time) 或 nil 
	// info, err := u.GetFileInfo("/test.txt")
	// log.Printf("GetFileInfo: %v\n", info)
	info, err := u.GetFileInfo("/test.jpg")
	log.Printf("GetFileInfo: %v\n", info)

	/// 删除文件
	log.Printf("DeleteFile: %v\n", u.DeleteFile("/test.jpg"))
	// log.Printf("DeleteFile: %v\n", u.DeleteFile("/test.txt"))

	/// 创建目录
	log.Printf("MkDir: %v\n", u.MkDir("/A", true))
	// 创建目录时可使用 u.mkDir("/a/b/c", true) 进行父级目录的自动创建（最深10级目录）
	log.Printf("MkDir: %v\n", u.MkDir("/A/B/C", true))
	log.Printf("MkDir: %v\n", u.MkDir("/1/2/3", false))
	log.Printf("MkDir: %v\n", u.MkDir("/folder", false))

	/// 删除目录（目录必须为空）
	log.Printf("RmDir: %v\n", u.RmDir("/A"))

	/// 读取目录
	dirs, err := u.ReadDir("/")
	log.Printf("ReadDir: %v\n", err)
	for i, d := range dirs {
		log.Printf("\t%d: %v\n", i, d)
	}
	dirs, err = u.ReadDir("/folder/")
	log.Printf("ReadDir: %v\n", err) // 必须斜杠结尾
	for i, d := range dirs {
		log.Printf("\t%d: %v\n", i, d)
	}
}
