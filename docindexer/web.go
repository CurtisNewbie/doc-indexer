package docindexer

import (
	"os"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/miso/miso"
	"github.com/gin-gonic/gin"
)

const (
	ResCodeBookmark = "manage-bookmarks"
	ResNameBookmark = "Manage Bookmarks"

	MsgUnknownErr  = "Unknown error, please try again"
	MsgUploadFiled = "Upload failed, please try again"
)

var (
	ManageBookmarkRes = goauth.Protected(ResNameBookmark, ResCodeBookmark)
)

type ListBookmarksReq struct {
	Name   *string
	Paging miso.Paging
}

func RegisterRoutes(rail miso.Rail) error {
	goauth.ReportPathsOnBootstrapped(rail)
	goauth.ReportResourcesOnBootstrapped(rail, []goauth.AddResourceReq{
		{Code: ResCodeBookmark, Name: ResNameBookmark},
	})

	miso.BaseRoute("/bookmark").Group(
		miso.Put("/file/upload", UploadBookmarkFileEp).
			Extra(ManageBookmarkRes),

		miso.IPost[ListBookmarksReq]("/list", ListBookmarksEp).
			Extra(ManageBookmarkRes),

		miso.IPost[RemoveBookmarkReq]("/remove", RemoveBookmarkEp).
			Extra(ManageBookmarkRes),
	)

	return nil
}

// Upload bookmark file endpoint.
func UploadBookmarkFileEp(c *gin.Context, rail miso.Rail) (any, error) {
	user := common.GetUser(rail)
	path, err := TransferTmpFile(rail, c.Request.Body)
	if err != nil {
		return nil, err
	}
	defer os.Remove(path)

	lock := miso.NewRLock(rail, "docindexer:bookmark:"+user.UserNo)
	if err := lock.Lock(); err != nil {
		rail.Errorf("failed to lock for bookmark upload, user: %v, %v", user.Username, err)
		return nil, miso.NewErr("Please try again later")
	}
	defer lock.Unlock()

	if err := ProcessUploadedBookmarkFile(rail, path, user); err != nil {
		rail.Errorf("ProcessUploadedBookmarkFile failed, user: %v, path: %v, %v", user.Username, path, err)
		return nil, miso.NewErr("Failed to parse bookmark file")
	}

	return nil, nil
}

// List bookmarks endpoint.
func ListBookmarksEp(c *gin.Context, rail miso.Rail, req ListBookmarksReq) (any, error) {
	user := common.GetUser(rail)
	return ListBookmarks(rail, miso.GetMySQL(), req, user.UserNo)
}

type RemoveBookmarkReq struct {
	Id int64
}

// Remove bookmark endpoint.
func RemoveBookmarkEp(c *gin.Context, rail miso.Rail, req RemoveBookmarkReq) (any, error) {
	user := common.GetUser(rail)
	return nil, RemoveBookmark(rail, miso.GetMySQL(), req.Id, user.UserNo)
}
