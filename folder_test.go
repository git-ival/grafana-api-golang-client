package gapi

import (
	"strings"
	"testing"

	"github.com/gobs/pretty"
)

const (
	getFoldersJSON = `{
    "id":1,
    "uid": "nErXDvCkzz",
    "title": "Departmenet ABC",
    "url": "/dashboards/f/nErXDvCkzz/department-abc",
    "hasAcl": false,
    "canSave": true,
    "canEdit": true,
    "canAdmin": true,
    "createdBy": "admin",
    "created": "2018-01-31T17:43:12+01:00",
    "updatedBy": "admin",
    "updated": "2018-01-31T17:43:12+01:00",
    "version": 1
	}`
	getFolderJSON = `
{
  "id":1,
  "uid": "nErXDvCkzz",
  "title": "Departmenet ABC",
  "url": "/dashboards/f/nErXDvCkzz/department-abc",
  "hasAcl": false,
  "canSave": true,
  "canEdit": true,
  "canAdmin": true,
  "createdBy": "admin",
  "created": "2018-01-31T17:43:12+01:00",
  "updatedBy": "admin",
  "updated": "2018-01-31T17:43:12+01:00",
  "version": 1
}
`
	createdFolderJSON = `
{
  "id":1,
  "uid": "nErXDvCkzz",
  "title": "Departmenet ABC",
  "url": "/dashboards/f/nErXDvCkzz/department-abc",
  "hasAcl": false,
  "canSave": true,
  "canEdit": true,
  "canAdmin": true,
  "createdBy": "admin",
  "created": "2018-01-31T17:43:12+01:00",
  "updatedBy": "admin",
  "updated": "2018-01-31T17:43:12+01:00",
  "version": 1
}
`
	createdSubfolderJSON = `
{
    "id": 2,
    "uid": "b7adb34f-08d9-4812-ae23-0b66612e44cb",
    "title": "Department ABC subfolder",
    "url": "/dashboards/f/b7adb34f-08d9-4812-ae23-0b66612e44cb/department-abc-subfolder",
    "parentUid": "nErXDvCkzz"
}
`
	updatedFolderJSON = `
{
  "id":1,
  "uid": "nErXDvCkzz",
  "title": "Departmenet DEF",
  "url": "/dashboards/f/nErXDvCkzz/department-def",
  "hasAcl": false,
  "canSave": true,
  "canEdit": true,
  "canAdmin": true,
  "createdBy": "admin",
  "created": "2018-01-31T17:43:12+01:00",
  "updatedBy": "admin",
  "updated": "2018-01-31T17:43:12+01:00",
  "version": 1
}
`
	deletedFolderJSON = `
{
  "message":"Folder deleted"
}
`
)

func TestFolders(t *testing.T) {
	mockData := strings.Repeat(getFoldersJSON+",", 1000) // make 1000 folders.
	mockData = "[" + mockData[:len(mockData)-1] + "]"    // remove trailing comma; make a json list.

	// This creates 1000 + 1000 + 1 (2001, 3 calls) worth of folders.
	client := gapiTestToolsFromCalls(t, []mockServerCall{
		{200, mockData},
		{200, mockData},
		{200, "[" + getFolderJSON + "]"},
	})

	const dashCount = 2001

	folders, err := client.Folders()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pretty.PrettyFormat(folders))

	if len(folders) != dashCount {
		t.Errorf("Length of returned folders should be %d", dashCount)
	}
	if folders[0].ID != 1 || folders[0].Title != "Departmenet ABC" {
		t.Error("Not correctly parsing returned folders.")
	}
	if folders[dashCount-1].ID != 1 || folders[dashCount-1].Title != "Departmenet ABC" {
		t.Error("Not correctly parsing returned folders.")
	}
}

func TestFolder(t *testing.T) {
	client := gapiTestTools(t, 200, getFolderJSON)

	folder := int64(1)
	resp, err := client.Folder(folder)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pretty.PrettyFormat(resp))

	if resp.ID != folder || resp.Title != "Departmenet ABC" {
		t.Error("Not correctly parsing returned folder.")
	}
}

func TestFolderByUid(t *testing.T) {
	client := gapiTestTools(t, 200, getFolderJSON)

	folder := "nErXDvCkzz"
	resp, err := client.FolderByUID(folder)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pretty.PrettyFormat(resp))

	if resp.UID != folder || resp.Title != "Departmenet ABC" {
		t.Error("Not correctly parsing returned folder.")
	}
}

func TestNewFolder(t *testing.T) {
	testCases := []struct {
		desc      string
		subfolder bool
	}{
		{
			desc:      "creating a parent folder only",
			subfolder: false,
		},
		{
			desc:      "creating a subfolder that contains the ParentUID",
			subfolder: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			client := gapiTestTools(t, 200, createdFolderJSON)

			resp, err := client.NewFolder("test-folder", "")
			if err != nil {
				t.Fatal(err)
			}

			parentUID := "nErXDvCkzz"
			if resp.UID != parentUID {
				t.Error("Not correctly parsing returned creation message.")
			}

			if tc.subfolder {
				client := gapiTestTools(t, 200, createdSubfolderJSON)

				resp, err := client.NewFolder("subfolder", parentUID)
				if err != nil {
					t.Fatal(err)
				}

				if resp.UID != "b7adb34f-08d9-4812-ae23-0b66612e44cb" {
					t.Error("Subfolder's UID does not match expected UID")
				}

				if resp.ParentUID != parentUID {
					t.Error("Failed to create subfolder with ParentUID")
				}
			}
		})
	}
}

func TestUpdateFolder(t *testing.T) {
	client := gapiTestTools(t, 200, updatedFolderJSON)

	err := client.UpdateFolder("nErXDvCkzz", "test-folder")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteFolder(t *testing.T) {
	client := gapiTestTools(t, 200, deletedFolderJSON)

	err := client.DeleteFolder("nErXDvCkzz")
	if err != nil {
		t.Fatal(err)
	}
}
