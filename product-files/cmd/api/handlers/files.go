package handlers

import (
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/critma/prodfiles/cmd/api/config"
	"github.com/critma/prodfiles/internal/store"
	"github.com/gin-gonic/gin"
)

type Files struct {
	log   *log.Logger
	store store.Storage
	cfg   *config.Config
}

func NewFiles(s store.Storage, l *log.Logger, cfg *config.Config) *Files {
	return &Files{store: s, log: l, cfg: cfg}
}

type InputParams struct {
	Id       string `form:"id" uri:"id"`
	FileName string `form:"filename" uri:"filename"`
}

func (f *Files) UploadREST(c *gin.Context) {
	var query InputParams
	err := c.ShouldBindUri(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not valid query params"})
		return
	}

	f.log.Println("Handle POST", "id:", query.Id, "filename:", query.FileName)

	f.saveFile(query.Id, query.FileName, c, c.Request.Body)
}

func (f *Files) GetFile(c *gin.Context) {
	var query InputParams
	err := c.ShouldBindUri(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not valid query params"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Date", "")
	c.File(f.cfg.BasePath + "/" + query.Id + "/" + query.FileName)
}

func (f *Files) UploadMultipart(c *gin.Context) {
	var form InputParams
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not valid query params"})
		return
	}

	formFile, err := c.FormFile("file")
	if err != nil {
		f.log.Printf("not valid form file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "not valid form file"})
		return
	}

	ff, err := formFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "not valid form file"})
		return
	}
	defer ff.Close()

	f.saveFile(form.Id, form.FileName, c, ff)

}

func (f *Files) saveFile(id, path string, c *gin.Context, r io.ReadCloser) {
	f.log.Println("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Println("Unable to save file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to save file"})
	}
}
