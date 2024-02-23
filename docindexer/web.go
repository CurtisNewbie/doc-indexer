package docindexer

import (
	"os"

	"github.com/curtisnewbie/gocommon/auth"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

const (
	ResourceManageBookmark = "manage-bookmarks"

	MsgUnknownErr  = "Unknown error, please try again"
	MsgUploadFiled = "Upload failed, please try again"
)

var (
	ErrUnknown     = miso.NewErrf("Unknown error, please try again")
	ErrUploadFiled = miso.NewErrf("Upload failed, please try again")
)

type ListBookmarksReq struct {
	Name   *string
	Paging miso.Paging
}

func RegisterRoutes(rail miso.Rail) error {
	auth.ExposeResourceInfo([]auth.Resource{
		{Code: ResourceManageBookmark, Name: "Manage Bookmarks"},
	})

	miso.BaseRoute("/bookmark").Group(
		miso.Put("/file/upload", UploadBookmarkFileEp).
			Desc("Upload bookmark file").
			Resource(ResourceManageBookmark),

		miso.IPost[ListBookmarksReq]("/list", ListBookmarksEp).
			Desc("List bookmarks").
			Resource(ResourceManageBookmark),

		miso.IPost[RemoveBookmarkReq]("/remove", RemoveBookmarkEp).
			Desc("Remove bookmark").
			Resource(ResourceManageBookmark),
	)

	return nil
}

// Upload bookmark file endpoint.
func UploadBookmarkFileEp(inb *miso.Inbound) (any, error) {
	rail := inb.Rail()
	_, r := inb.Unwrap()
	user := common.GetUser(rail)
	path, err := TransferTmpFile(rail, r.Body)
	if err != nil {
		return nil, err
	}
	defer os.Remove(path)

	lock := miso.NewRLock(rail, "docindexer:bookmark:"+user.UserNo)
	if err := lock.Lock(); err != nil {
		rail.Errorf("failed to lock for bookmark upload, user: %v, %v", user.Username, err)
		return nil, miso.NewErrf("Please try again later")
	}
	defer lock.Unlock()

	if err := ProcessUploadedBookmarkFile(rail, path, user); err != nil {
		rail.Errorf("ProcessUploadedBookmarkFile failed, user: %v, path: %v, %v", user.Username, path, err)
		return nil, miso.NewErrf("Failed to parse bookmark file")
	}

	return nil, nil
}

// List bookmarks endpoint.
func ListBookmarksEp(inb *miso.Inbound, req ListBookmarksReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return ListBookmarks(rail, miso.GetMySQL(), req, user.UserNo)
}

type RemoveBookmarkReq struct {
	Id int64
}

// Remove bookmark endpoint.
func RemoveBookmarkEp(inb *miso.Inbound, req RemoveBookmarkReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, RemoveBookmark(rail, miso.GetMySQL(), req.Id, user.UserNo)
}
