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
	fmt.Println("Files:")
	println(len(r.Files))
	fileslist:= make([]string, len(r.Files))
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		var wg sync.WaitGroup
		var mu sync.Mutex

		for j, i := range r.Files {
			fmt.Printf("%s (%s))\n", i.Name, i.Id)
			//fileslist = append(fileslist,i.Name)
			fileslist[j]= i.Name
			wg.Add(1)

			go func() {
				mu.Lock()
				err = load(srv, i)
				mu.Unlock()
				wg.Done()
			}()

		}
		wg.Wait()
	}
	return fileslist,nil
}

func load(srv *drive.Service, r *drive.File) error {
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
	println("asasas")
	return nil
}