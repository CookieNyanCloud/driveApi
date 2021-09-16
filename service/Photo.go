package service

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
	"log"
	"os"
	"sync"
)

func GetPhoto(srv *drive.Service, name string) ([]string, error) {
	query := `name contains '` + name + "'"
	r, err := srv.Files.
		List().
		PageSize(20).
		Fields("nextPageToken, files(id, name)").
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		Q(query).
		IncludePermissionsForView("published").
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fileslist := make([]string, len(r.Files))
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		var wg sync.WaitGroup
		for j, i := range r.Files {
			wg.Add(1)
			fileslist[j] = i.Name
			go func(srv *drive.Service, i *drive.File, wg *sync.WaitGroup) {
				err = load(srv, i, wg)
				if err != nil {
					log.Fatalf("Unable to retrieve files2: %v", err)
				}
			}(srv, i, &wg)
		}
		wg.Wait()
	}
	return fileslist, nil
}

//func loadMult(srv *drive.Service, r *drive.FileList) ([][]byte, error) {
//	println("mult")
//	done := make(chan []byte, len(r.Files))
//	errch := make(chan error, len(r.Files))
//	for _, file := range r.Files {
//		go func(file *drive.File) {
//			b,err := load2(srv, file)
//			if err != nil {
//				errch <- err
//				done <- nil
//				return
//			}
//			done <- b
//			errch <- nil
//			//fileslist[i] = file.Name
//			//fmt.Println("121212",fileslist[i])
//		}(file)
//	}
//	bytesArray := make([][]byte, 0)
//	var errStr string
//	for i := 0; i < len(r.Files); i++ {
//		bytesArray = append(bytesArray, <-done)
//		if err := <-errch; err != nil {
//			errStr = errStr + " " + err.Error()
//		}
//	}
//	var err error
//	if errStr!=""{
//		err = errors.New(errStr)
//	}
//	return bytesArray, err
//	//return fileslist, nil
//}
//
//func load2(srv *drive.Service, r *drive.File) ([]byte,error) {
//	fmt.Println(r.Name)
//	println("ld2")
//	res, err := srv.Files.Get(r.Id).Download()
//
//	if err != nil {
//		return nil,err
//	}
//	defer res.Body.Close()
//	if res.StatusCode != 200 {
//		return nil,err
//	}
//	var data bytes.Buffer
//	_, err = io.Copy(&data, res.Body)
//	if err != nil {
//		return nil, err
//	}
//	println("asasas")
//	return data.Bytes(), nil
//
//	//fileNew, err := os.Create(r.Name)
//	//if err != nil {
//	//	return err
//	//}
//	//defer fileNew.Close()
//	//_, err = io.Copy(fileNew, res.Body)
//	//if err != nil {
//	//	return err
//	//}
//	//return nil
//}

func load(srv *drive.Service, r *drive.File, wg *sync.WaitGroup) (error) {
	res, err := srv.Files.Get(r.Id).Download()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return err
	}
	fileNew, err := os.Create(r.Name)
	if err != nil {
		return err
	}
	defer fileNew.Close()
	_, err = io.Copy(fileNew, res.Body)
	if err != nil {
		return err
	}
	wg.Done()
	return nil
}

//
//
//func load(srv *drive.Service, r *drive.File) error {
//	res, err := srv.Files.Get(r.Id).Download()
//
//	if err != nil {
//		return err
//	}
//	defer res.Body.Close()
//	if res.StatusCode != 200 {
//		return err
//	}
//	fileNew, err := os.Create(r.Name)
//	if err != nil {
//		return err
//	}
//	defer fileNew.Close()
//	_, err = io.Copy(fileNew, res.Body)
//	if err != nil {
//		return err
//	}
//	println("asasas")
//	return nil
//}
