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
		Fields("nextPageToken, files(id, name, parents)").
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

func SendPhoto(srv *drive.Service, name,author, dirType,drivePeople, driveZag string) error {
	src, err := os.Open(name)
	if err != nil {
		return err
	}
	defer src.Close()
	fl:= &drive.File{}
	fl.Name = name
	fl.MimeType = "image/jpeg"
	fl.Description = "автор:"+author
	var folder string
	switch dirType {
	case "з":
		folder =driveZag
	case "л":
		folder = drivePeople
	default:
		folder = drivePeople
	}
	fs:=make([]string,1)
	fs[0] = folder
	fl.Parents = fs
	file, err:= srv.
		Files.
		Create(fl).
		SupportsAllDrives(true).
		SupportsTeamDrives(true).
		Media(src).
		Do()
	fmt.Println(file.Name)
	return nil
}


func load(srv *drive.Service, r *drive.File, wg *sync.WaitGroup) error {
	println(r.Name,"start")
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
	println(r.Name,"done")
	return nil
}
