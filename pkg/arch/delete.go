package arch

import (
	"fmt"
	"os"
)

func MyDelete(name string) error {
	return os.Remove(name)
}

func AllDelete(names []string) error {
	for _, v := range names {
		fmt.Println(v)
		err := os.Remove(v)
		if err != nil {
			return err
		}
	}
	return nil
}
