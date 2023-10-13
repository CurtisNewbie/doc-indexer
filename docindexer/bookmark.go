package docindexer

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	TagA        = "a" // bookmark file <A>
	AttrHref    = "href"
	AttrAddDate = "add_date"
	AttrIcon    = "icon"
)

type NetscapeBookmarkFile struct {
	Bookmarks []Bookmark
}

func (n *NetscapeBookmarkFile) Add(b Bookmark) {
	n.Bookmarks = append(n.Bookmarks, b)
}

type Bookmark struct {
	Name    string
	Href    string
	Icon    string
	AddDate string
}

func (b Bookmark) String() string {
	return fmt.Sprintf("Bookmark{\n\tName: %v,\n\tHref: %v,\n\tIcon: %v,\n\tAddDate: %v\n}", b.Name, b.Href, b.Icon, b.AddDate)
}

func ParseNetscapeBookmark(rail miso.Rail, body io.Reader) (NetscapeBookmarkFile, error) {
	bookmarkFile := NetscapeBookmarkFile{Bookmarks: []Bookmark{}}

	z := html.NewTokenizer(body)
	curr := Bookmark{}

	for {
		ttype := z.Next()
		bname, isAttr := z.TagName()
		name := string(bname)

		textB := z.Text()
		text := string(textB)

		attr := map[string]string{}

		for {
			attrKeyB, attrValB, more := z.TagAttr()
			k := string(attrKeyB)
			v := string(attrValB)
			attr[k] = v
			if !more {
				break
			}
		}

		rail.Debugf("tokenType: %v, text: %v, name: %v, isAttr: %v, attr: %v",
			ttype, text, name, isAttr, attr)

		switch ttype {
		case html.ErrorToken:
			err := z.Err()
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return bookmarkFile, err
		case html.TextToken:
			if !isAttr {
				curr.Name = text
			}
		case html.StartTagToken:
			if name == TagA && isAttr {
				if v, ok := attr[AttrHref]; ok {
					curr.Href = v
				}
				if v, ok := attr[AttrAddDate]; ok {
					curr.AddDate = v
				}
				if v, ok := attr[AttrIcon]; ok {
					curr.Icon = v
				}
			}
		case html.EndTagToken:
			if name == TagA {
				bookmarkFile.Add(curr)
				curr = Bookmark{}
			}
		}
	}
}

func DownloadAsTmpFile(rail miso.Rail, req ParseBookmarkFileReq) (string, error) {
	tkn, err := GetFstoreTmpToken(rail, req.TempFileId)
	if err != nil {
		return "", miso.NewErr(MsgUploadExpired, "GetFstoreTmpToken failed, tempFileId: %v, %v", req.TempFileId, err)
	}

	path := TempFilePath(tkn)
	if err := DownloadFstoreFile(rail, tkn, path); err != nil {
		rail.Warnf("DownloadFstoreFile failed, tempFileId: %v, path: %v, %v", req.TempFileId, path, err)
		return "", miso.NewErr(MsgUnknownErr, "downloaded fstore file to %v, tempFileId: %v", path, req.TempFileId)
	}
	rail.Infof("Downloaded fstore file, tempFileId: %v, path: %v", req.TempFileId, path)
	return path, nil
}

func ProcessUploadedBookmarkFile(rail miso.Rail, req ParseBookmarkFileReq, path string, user common.User) error {
	rail.Infof("User '%v' parse bookmark file, fileId: '%v'", user.Username, req.TempFileId)
	file, err := os.Open(path)
	if err != nil {
		return miso.NewErr(MsgUnknownErr, "open downloaded fstore file failed, path: %v, tempFileId: %v",
			path, req.TempFileId)
	}

	bookmarkFile, err := ParseNetscapeBookmark(rail, file)
	if err != nil {
		return miso.NewErr(MsgUnknownErr, "open downloaded fstore file failed, path: %v, tempFileId: %v",
			path, req.TempFileId)
	}

	return SaveBookmarks(rail, miso.GetMySQL(), bookmarkFile, user)
}

type NewBookmark struct {
	UserNo string
	Icon   string
	Name   string
	Href   string
	Md5    string
}

func SaveBookmarks(rail miso.Rail, tx *gorm.DB, bookmarkFile NetscapeBookmarkFile, user common.User) error {

	bookmarks := make([]NewBookmark, 0, len(bookmarkFile.Bookmarks))
	for i := range bookmarkFile.Bookmarks {
		bm := bookmarkFile.Bookmarks[i]
		md5 := BookmarkMd5(bm)
		bookmarks = append(bookmarks, NewBookmark{
			UserNo: user.UserNo,
			Icon:   bm.Icon,
			Name:   bm.Name,
			Href:   bm.Href,
			Md5:    md5,
		})
	}
	// rail.Debugf("bookmarks: %+v", bookmarks)

	err := tx.Table("bookmark").Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(bookmarks, 100).Error
	if err != nil {
		return fmt.Errorf("failed to insert bookmark, %v", err)
	}

	return nil
}

func BookmarkMd5(bm Bookmark) string {
	s := fmt.Sprintf("NA%vHR%vIC%v", bm.Name, bm.Href, bm.Icon)
	chksum := md5.Sum([]byte(s))
	return hex.EncodeToString(chksum[:])
}
