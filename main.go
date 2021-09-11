package main

import (
	"context"
	"fmt"
	"github.com/CookieNyanCloud/driveApi/api"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b,
		//drive.DriveScope,
		//drive.DriveAppdataScope,
		//drive.DriveFileScope,
		//drive.DriveMetadataScope,
		//drive.DriveMetadataReadonlyScope,
		//drive.DrivePhotosReadonlyScope,
		drive.DriveReadonlyScope,
		//drive.DriveScriptsScope,
	)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := api.GetClient(config)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	name := "рашкин"
	if err = GetPhoto(srv, name); err != nil {

	}

}

func GetPhoto(srv *drive.Service, name string) error {
	query := `name contains '` + name + "'"
	r, err := srv.Files.
		List().
		PageSize(20).
		Fields("nextPageToken, files(id, name)").
		IncludeItemsFromAllDrives(true).
		//Corpora("drive").
		SupportsAllDrives(true).
		Q(query).
		IncludePermissionsForView("published").
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s))\n", i.Name, i.Id)
		}
	}

	res, err := srv.Files.Get(r.Files[0].Id).Download()

	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return err
	}
	fileNew, err := os.Create("new.jpeg")
	if err != nil {
		return err
	}
	defer fileNew.Close()
	_, err = io.Copy(fileNew, res.Body)
	if err != nil {
		return err
	}
	return nil
}
