package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	//r, err := srv.Files.List().PageSize(10).
	//	Fields("nextPageToken, files(id, name)").Do()
	//if err != nil {
	//	log.Fatalf("Unable to retrieve files: %v", err)
	//}
	//fmt.Println("Files:")
	//if len(r.Files) == 0 {
	//	fmt.Println("No files found.")
	//} else {
	//	for _, i := range r.Files {
	//		fmt.Printf("%s (%s)\n", i.Name, i.Id)
	//	}
	//}

	// Step 1: Open  file
	//f, err := os.Open("sample.txt")
	//
	//if err != nil {
	//	panic(fmt.Sprintf("cannot open file: %v", err))
	//}
	//defer f.Close()

	// Step 2: Get the Google Drive service
	//srv, err := getDriveService()

	// Step 3: Create directory
	// dir, err := createFolder(srv, "New Folder", "root")

	// if err != nil {
	// 	panic(fmt.Sprintf("Could not create dir: %v\n", err))
	// }

	//give your drive folder id here in which you want to upload or create a new directory
	//folderId := "19Ghk_vKLc8BtfJJTOpiF9dzp1IB1gvYm"
	folderPeople := "19Ghk_vKLc8BtfJJTOpiF9dzp1IB1gvYm"

	// Step 4: create the file and upload
	//file, err := createFile(srv, f.Name(), "application/octet-stream", f, folderId)
	var name string
	name = "Юлия Галямина"
	err = getFile(srv,name,folderPeople)
	if err != nil {
		panic(fmt.Sprintf("Could not create file: %v\n", err))
	}

	//fmt.Printf("File '%s' uploaded successfully", file.Name)
	//fmt.Printf("\nFile Id: '%s' ", file.Id)


}

func getDriveService() (*drive.Service, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		fmt.Printf("Unable to read credentials.json file. Err: %v\n", err)
		return nil, err
	}

	// If you want to modifyt this scope, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)

	if err != nil {
		return nil, err
	}

	client := getClient(config)

	service, err := drive.New(client)

	if err != nil {
		fmt.Printf("Cannot create the Google Drive service: %v\n", err)
		return nil, err
	}

	return service, err
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}

func getFile(service *drive.Service, name string, parentId string) error {
	var pageToken string
	sumOwner:= parentId + "in parents"
	file, err := service.Files.List().
		Q("name contains 'Галямина'").
		Q(sumOwner).
		Fields("nextPageToken, files(id, name)").
		Spaces("drive").
		PageToken(pageToken).Corpora()
		Do()
	if err != nil {
		log.Println("Could not: " + err.Error())
		return err
	}

	res, err:= service.Files.
		Export(file.Files[0].Id,"image/jpeg").
		Download()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	fileNew, err := os.Create(name)
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

