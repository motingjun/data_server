package reader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/axgle/mahonia"
	"github.com/extrame/xls"
)

func ReadXls(filenames []string) [][]string {
	pwd, _ := os.Getwd() // 获取到当前目录，相当于python里的os.getcwd()
	fmt.Println("当前的操作路径为:", pwd)
	total_data := make([][]string, 0)
	for _, filename := range filenames {
		if filename == "" {
			continue
		}
		//文件路径拼接
		xlsPath := filepath.Join(pwd, "upload", filename)
		fmt.Println("文件的路径为:", xlsPath)

		//通过os.Stat()函数返回的文件状态，如果有错误则根据错误状态来判断文件或者文件夹是否存在
		fileinfo, err := os.Stat(xlsPath)
		if err != nil {
			fmt.Println(err.Error())
			if os.IsNotExist(err) {
				fmt.Println("file:", xlsPath, " not exist！")
			}
		} else {
			fmt.Println(fileinfo.Name())
		}

		// 读取xls文件
		xlsFile, closer, xlsErr := OpenWithCloser(xlsPath, "utf-8")
		if xlsErr != nil {
			continue
		}
		// 获取xls文件的第一个sheet
		// sheet := xlsFile.GetSheet(0)
		// if sheet.MaxRow != 0 {
		// 	temp := make([][]string, sheet.MaxRow-1)
		// 	// 从第二行开始，遍历xls文件
		// 	for i := 1; i <= int(sheet.MaxRow); i++ {
		// 		row := sheet.Row(i)
		// 		data := make([]string, 0)
		// 		// var data [9]string
		// 		fmt.Println(row.LastCol)
		// 		if row.LastCol() > 14 {
		// 			for j := 0; j < row.LastCol(); j++ {
		// 				// if j == 5 || j == 10 || j == 11 || j == 12 {
		// 				// 	continue
		// 				// }
		// 				col := row.Col(j)
		// 				data = append(data, col)
		// 				fmt.Println(col)
		// 			}
		// 			temp[i] = data
		// 			// temp = append(temp, data)
		// 			// fmt.Println(temp)
		// 		}
		// 	}

		// 	fmt.Println("---------start----------------")
		// 	total_data = append(total_data, temp...)
		// 	fmt.Println("---------end----------------")
		// } else {
		// 	fmt.Println("open xls file err!")
		// }

		if sheet := xlsFile.GetSheet(0); sheet != nil {
			columns := [9]int{1, 2, 3, 4, 6, 7, 11, 12, 13}
			// for i := 0; i <= int(sheet.MaxRow); i++ {
			for i := 1; i <= 100; i++ {
				row := sheet.Row(i)
				data := make([]string, 0)
				// fmt.Println(row)
				for _, column := range columns {
					col := row.Col(column)
					data = append(data, col)
					// fmt.Println(col)
				}
				fmt.Println(data)
			}
		}

		closer.Close()

		remove_err := os.Remove(xlsPath)
		if remove_err != nil {
			fmt.Println("file remove Error!")
			fmt.Printf("%s", remove_err)
		} else {
			fmt.Println("file remove OK!")
		}
	}
	return total_data
}

//Open one xls file and return the closer
func OpenWithCloser(file string, charset string) (*xls.WorkBook, io.Closer, error) {
	if fi, err := os.Open(file); err == nil {
		wb, err := xls.OpenReader(fi, charset)
		return wb, fi, err
	} else {
		return nil, nil, err
	}
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
