package reader

import (
	"fmt"
	"os"
	"path/filepath"

	"data_server/logger"

	emoji "github.com/Andrew-M-C/go.emoji"
	"github.com/xuri/excelize/v2"
)

func ReadXlsx(filenames []string, log *logger.Log) [][]string {
	pwd, _ := os.Getwd()
	// fmt.Println("当前的操作路径为:", pwd)
	total_data := make([][]string, 0)

	for _, filename := range filenames {
		if filename == "" {
			continue
		}
		//文件路径拼接
		xlsxPath := filepath.Join(pwd, "upload", filename)

		//通过os.Stat()函数返回的文件状态，如果有错误则根据错误状态来判断文件或者文件夹是否存在
		fileinfo, err := os.Stat(xlsxPath)
		if err != nil {
			log.Logger.Errorln(fmt.Sprintf("file stat is wrong, error : %s", err.Error()))
			if os.IsNotExist(err) {
				log.Logger.Errorln(fmt.Sprintf("file not exist : %s", xlsxPath))
			}
		} else {
			log.Logger.Infoln(fmt.Sprintf("file name is %s", fileinfo.Name()))
		}

		// 读取xlsx文件
		xlsxFile, err := excelize.OpenFile(xlsxPath)
		if err != nil {
			log.Logger.Errorln(fmt.Sprintf("excelize open xlsxFile fail, the error : %s", err))
		}
		// 遍历明为'Sheet'的工作簿
		rows, err := xlsxFile.GetRows("Sheet")
		if err != nil {
			log.Logger.Errorln(fmt.Sprintf("excelize read 'Sheet' book fail, the error : %s", err))
		}
		sheet_data := make([][]string, len(rows))
		for i, row := range rows {
			data := make([]string, 0)
			data = append(data, row...) // 追加列数据

			sheet_data[i] = data
		}

		// 删除本地文件
		remove_err := os.Remove(xlsxPath)
		if remove_err != nil {
			log.Logger.Errorln(fmt.Sprintf("file remove fail, the error : %s", remove_err))
		} else {
			log.Logger.Infoln(fmt.Sprintf("file remove ok, file name is %s", filename))
		}

		total_data = append(total_data, sheet_data...)
	}
	return total_data
}

func SyncReadXlsx(filename string, log *logger.Log) [][]string {
	if filename != "" {
		pwd, _ := os.Getwd()
		xlsxPath := filepath.Join(pwd, "upload", filename)
		fileinfo, err := os.Stat(xlsxPath)
		if err != nil {
			log.Logger.Errorln(fmt.Sprintf("file stat is wrong, error : %s", err.Error()))
			if os.IsNotExist(err) {
				log.Logger.Errorln(fmt.Sprintf("file not exist : %s", xlsxPath))
			}
		} else {
			log.Logger.Infoln(fmt.Sprintf("file name is %s", fileinfo.Name()))
		}

		// 读取xlsx文件
		xlsxFile, err := excelize.OpenFile(xlsxPath)
		if err != nil {
			log.Logger.Errorln(fmt.Sprintf("excelize open xlsxFile fail, the error : %s", err))
		}
		// 遍历明为'Sheet'的工作簿
		rows, err := xlsxFile.GetRows("Sheet")
		if err != nil {
			log.Logger.Errorln(fmt.Sprintf("excelize read 'Sheet' book fail, the error : %s", err))
		}
		sheet_data := make([][]string, len(rows))
		for i, row := range rows {
			data := make([]string, 0)
			data = append(data, row...)           // 追加列数据
			data[1] = handleTitleByEmoji(data[1]) // 处理标题中的颜文字(一个字符独占4个字节)

			sheet_data[i] = data
		}

		// 删除本地文件
		remove_err := os.Remove(xlsxPath)
		if remove_err != nil {
			log.Logger.Errorln(fmt.Sprintf("file remove fail, the error : %s", remove_err))
		} else {
			log.Logger.Infoln(fmt.Sprintf("file remove ok, file name is %s", filename))
		}
		return sheet_data
	}
	return nil
}

// 颜文字处理
func handleTitleByEmoji(title string) string {
	if len(title) != 0 {

		new_title := emoji.ReplaceAllEmojiFunc(title, func(emoji string) string {
			return fmt.Sprintf("%s", "")
		})
		return new_title
	}
	return title
}
