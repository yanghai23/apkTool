package analysis

// 根据读取的配置信息，可以处理打包流程

import (
	"customConfig"
	"fmt"
	"github.com/yanghai23/GoLib/atfile"
	"merge"
	"model"
	"parse"
	"replace"
	"rjar"
	"sdkcfg"
	"utils"
)

func ExplainChannels(apkToolsPath, workPath string, game *model.Game, channels []*model.GameChannel) {
	for _, gameChannel := range channels {
		fmt.Println("开始打包，渠道【" + gameChannel.Id + "】")
		//读取sdk的配置信息
		sdkPath := atfile.GetCurrentDirectory() + "/config/sdk/" + gameChannel.Id
		fmt.Println("sdkPath:", sdkPath)
		sdkConfig, _ := parse.ReadSdkConfig(sdkPath)
		fmt.Println(sdkConfig)

		fmt.Println("清空temp目录")
		tempPath := utils.CreateNewFolder(atfile.GetCurrentDirectory() + "/work/temp")

		//复制原始的文件
		fmt.Println("复制原始文件到新的目录:", tempPath)
		fmt.Println(tempPath)
		utils.CopyBakToTemp(atfile.GetCurrentDirectory()+"/bak", tempPath, func() string {
			return "../../bak"
		})
		fmt.Println("修改包名")
		newPackageVal := merge.RenamePackage(gameChannel, tempPath)
		fmt.Println("NewPageName:", newPackageVal)

		fmt.Println("合并资源")
		// 将配置的jar，res，等资源进行合并
		fmt.Println("将配置的jar，res，等资源进行合并")
		fmt.Println(sdkConfig.Config.Operations)
		merge.MergeSource(sdkPath, tempPath, sdkConfig.Config.Operations)

		//是否添加闪屏图片
		fmt.Println("添加闪屏图片")
		merge.AddSplashImg(sdkPath, tempPath, gameChannel, game)
		//创建配置文件
		customConfig.CreateCustomConfig(tempPath, gameChannel, &sdkConfig.Config)

		//处理Icon图标
		merge.MergeIcon(sdkPath, tempPath, gameChannel)

		fmt.Println("生成R文件")
		rjar.ComplieR(apkToolsPath, tempPath, workPath, newPackageVal, &sdkConfig.Config)

		fmt.Println("jar2smali")
		//jar2smali.Jar2Smali(apkToolsPath, tempPath)

		fmt.Println("合并meta-data")
		merge.MergeMetaData(tempPath, gameChannel)

		fmt.Println("添加启动页面")
		maniActivityName := merge.AddSplashActivity(tempPath, gameChannel)

		merge.MergeAndroidManifest(sdkPath, tempPath, game)
		replace.ReplacePkgManifest(tempPath, newPackageVal, gameChannel)
		sdkcfg.CreateSdkConfigXml(gameChannel, sdkConfig, game, maniActivityName, tempPath)
		utils.CreateOutDir()
	}
}
