package main

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/go-flac/flacpicture/v2"
	"github.com/go-flac/flacvorbis/v2"
	"github.com/go-flac/go-flac/v2"
	"github.com/h2non/filetype"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OpenMusicFile() string {
	filepath, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a music file",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Audio (*.flac)",
				Pattern:     "*.flac",
			},
		},
	})
	return filepath
}

type MusicInfo struct {
	Size       string
	SampleRate string
	Comments   map[string]string
	Picture    string
}

func extractComments(f *flac.File) (*flacvorbis.MetaDataBlockVorbisComment, int) {
	var cmts *flacvorbis.MetaDataBlockVorbisComment
	var cmtIdx int
	var err error
	for idx, meta := range f.Meta {
		if meta.Type == flac.VorbisComment {
			cmts, err = flacvorbis.ParseFromMetaDataBlock(*meta)
			cmtIdx = idx
			if err != nil {
				return nil, -1
			}
			break
		}
	}
	return cmts, cmtIdx
}

func (a *App) CheckMusicType(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	header := make([]byte, 261)
	_, err = file.Read(header)
	if err != nil {
		return false
	}

	kind, err := filetype.Match(header)
	if err != nil {
		return false
	}
	return kind.MIME.Value == "audio/x-flac"
}

func (a *App) ParseMusicFile(path string) *MusicInfo {
	var musicInfo MusicInfo

	// Size
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil
	}
	musicInfo.Size = humanize.Bytes(uint64(fileInfo.Size()))

	f, err := flac.ParseFile(path)
	if err != nil {
		return nil
	}

	// SampleRate
	data, err := f.GetStreamInfo()
	if err != nil {
		return nil
	}
	musicInfo.SampleRate = humanize.SI(float64(data.SampleRate), "Hz")

	// Comments
	cmts, _ := extractComments(f)
	musicInfo.Comments = make(map[string]string)
	keys := []string{
		"TITLE", "VERSION", "ALBUM", "ALBUMARTIST", "ARTIST",
		"LYRICIST", "COMPOSER", "PERFORMER", "TRACKNUMBER", "COPYRIGHT",
		"LICENSE", "ORGANIZATION", "DESCRIPTION", "GENRE", "DATE",
		"LOCATION", "CONTACT", "ISRC", "LYRICS",
	}
	for _, key := range keys {
		tmp, _ := cmts.Get(key)
		if len(tmp) > 0 {
			musicInfo.Comments[key] = tmp[0]
		} else {
			musicInfo.Comments[key] = ""
		}
	}

	// Picture
	var pic *flacpicture.MetadataBlockPicture
	for _, meta := range f.Meta {
		if meta.Type == flac.Picture {
			pic, _ = flacpicture.ParseFromMetaDataBlock(*meta)
			musicInfo.Picture = "data:" + pic.MIME + ";base64," + base64.StdEncoding.EncodeToString(pic.ImageData)
			break
		}
	}
	return &musicInfo
}

func (a *App) SaveMusicFile(path string, comments map[string]string) bool {
	f, err := flac.ParseFile(path)
	if err != nil {
		return false
	}
	_, cmtIdx := extractComments(f)
	newCmts := flacvorbis.New()
	for k, v := range comments {
		if v != "" {
			newCmts.Add(k, v)
		}
	}
	newCmtsMeta := newCmts.Marshal()
	if cmtIdx > 0 {
		f.Meta[cmtIdx] = &newCmtsMeta
	} else {
		f.Meta = append(f.Meta, &newCmtsMeta)
	}
	os.Mkdir("./tmp", 0755)
	err = f.Save(filepath.Join("./tmp", filepath.Base(path)))
	if err != nil {
		return false
	}
	os.Rename(filepath.Join("./tmp", filepath.Base(path)), path)
	return err == nil
}
