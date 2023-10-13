package docindexer

import (
	"io"
	"os"

	"github.com/curtisnewbie/miso/miso"
)

func GetFstoreTmpToken(rail miso.Rail, fileId string) (string /* tmpToken */, error) {
	var res miso.GnResp[string]
	err := miso.NewDynTClient(rail, "/file/key", "fstore").
		EnableTracing().
		AddQueryParams("fileId", fileId).
		Get().
		Json(res)

	if err != nil {
		return "", err
	}
	return res.Res()
}

func DownloadFstoreFile(rail miso.Rail, tmpToken string, absPath string) error {
	r := miso.NewDynTClient(rail, "/file/raw", "fstore").
		EnableTracing().
		AddQueryParams("key", tmpToken).
		Get()
	if r.Err != nil {
		return r.Err
	}
	defer r.Close()

	out, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r.Resp.Body)
	if err != nil {
		return err
	}
	return nil
}
