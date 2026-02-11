package term

import (
	"io"
	"runtime"
)

// readPasswordLine 从读取器读取直到找到 \n 或 io.EOF。
// 返回的切片不包含 \n。
// readPasswordLine 还会忽略它找到的任何 \r。
// Windows 使用 \r 作为行尾。因此，在 Windows 上，readPasswordLine
// 读取直到找到 \r，并忽略处理过程中找到的任何 \n。
func readPasswordLine(reader io.Reader) ([]byte, error) {
	var buf [1]byte
	var ret []byte

	for {
		n, err := reader.Read(buf[:])
		if n > 0 {
			switch buf[0] {
			case '\b':
				if len(ret) > 0 {
					ret = ret[:len(ret)-1]
				}
			case '\n':
				if runtime.GOOS != "windows" {
					return ret, nil
				}
				// otherwise ignore \n
			case '\r':
				if runtime.GOOS == "windows" {
					return ret, nil
				}
				// otherwise ignore \r
			default:
				ret = append(ret, buf[0])
			}
			continue
		}
		if err != nil {
			if err == io.EOF && len(ret) > 0 {
				return ret, nil
			}
			return ret, err //nolint:wrapcheck
		}
	}
}
