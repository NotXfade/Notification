package methods

import "errors"

//CheckFileSize function used for check file size
func CheckFileSize(size int) error {

	var kilobytes int
	kilobytes = (size / 1024)

	//if file less then 1 kb return nil
	if kilobytes < 0 {
		return nil
	}
	var megabytes float64
	megabytes = (float64)(kilobytes / 1024)
	//if file less then 5 mb return nil
	if megabytes < 50 {
		return nil
	}
	return errors.New("File size must be less then 50 mb")
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
