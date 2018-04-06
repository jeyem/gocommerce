package file

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type File struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Owner     bson.ObjectId `bson:"owner,omitempty"`
	Name      string        `bson:"name"`
	Path      string        `bson:"path"`
	Format    string        `bson:"format"`
	CheckSum  string        `bson:"check_sum"`
	Keywords  []string      `bson:"keywords"`
	CreatedAt time.Time     `bson:"created_at"`
}

func (f *File) Ext() string {
	return filepath.Ext(f.Name)
}

func (f *File) NameNoExt() string {
	return strings.Replace(f.Name, f.Ext(), "", -1)
}

func (f *File) Orginal() string {
	return filepath.Join(config.ServingPrefix, f.Name)
}

func (f *File) Thumb() string {
	if f.Format != Image {
		return f.Orginal()
	}
	suffix := "*w150h150"
	name := f.NameNoExt() + suffix + f.Ext()
	return filepath.Join(config.ServingPrefix, name)
}

func (f *File) Micro() string {
	if f.Format != Image {
		return f.Orginal()
	}
	suffix := "*w50h50"
	name := f.NameNoExt() + suffix + f.Ext()
	return filepath.Join(config.ServingPrefix, name)
}

func (f *File) Mid() string {
	if f.Format != Image {
		return f.Orginal()
	}
	suffix := "*w500h500"
	name := f.NameNoExt() + suffix + f.Ext()
	return filepath.Join(config.ServingPrefix, name)
}

func (f *File) Big() string {
	if f.Format != Image {
		return f.Orginal()
	}
	suffix := "*w900h900"
	name := f.NameNoExt() + suffix + f.Ext()
	return filepath.Join(config.ServingPrefix, name)
}

func (f *File) Rest() map[string]interface{} {
	if f.Format != Image {
		return map[string]interface{}{
			"name":     f.NameNoExt(),
			"original": f.Orginal(),
		}
	}
	return map[string]interface{}{
		"name":     f.NameNoExt(),
		"original": f.Orginal(),
		"micro":    f.Micro(),
		"thumb":    f.Thumb(),
		"mid":      f.Mid(),
		"big":      f.Big(),
	}
}

func (f *File) Save() error {
	if f.ID.Valid() {
		return db.Update(f)
	}
	f.Keywords = []string{
		f.Owner.Hex(), f.Name, f.CheckSum,
	}
	duplicate := new(File)
	if err := duplicate.Load(f.CheckSum, f.Owner); err == nil {
		return errors.New("duplicate image entry")
	}
	f.CreatedAt = time.Now()
	return db.Create(f)
}

func (f *File) Load(checksum string, owner bson.ObjectId) error {
	return db.Where(bson.M{
		"owner":     owner,
		"check_sum": checksum,
	}).Find(f)
}

func (f *File) LoadByName(owner bson.ObjectId, name string) error {
	return db.Where(bson.M{
		"owner":     owner,
		"check_sum": name,
	}).Find(f)
}

func (File) Meta() []mgo.Index {
	return []mgo.Index{
		{Key: []string{"owner"}},
		{Key: []string{"owner", "check_sum"}},
		{Key: []string{"created_at"}},
		{Key: []string{"owner", "created_at"}},
	}
}

func LoadFilesByName(owner bson.ObjectId, filesPath ...string) (res []map[string]interface{}) {
	for _, path := range filesPath {
		file := new(File)
		if err := file.LoadByName(owner, path); err != nil {
			file.Name = path
		}
		res = append(res, file.Rest())

	}
	return res
}

func GetFile(owner bson.ObjectId, path string) string {
	name := getFileOrginalName(path)
	base := filepath.Base(path)
	checksum := strings.Replace(name, filepath.Ext(name), "", -1)
	file := new(File)
	var (
		fullpath    string
		orginalPath string
	)
	if err := db.Where(bson.M{
		"check_sum": checksum, "owner": owner}).Find(file); err == nil {
		fullpath = filepath.Join(config.ImagePath,
			filepath.Dir(file.Path), base)
		orginalPath = filepath.Join(config.ImagePath, file.Path)
	} else {
		fullpath = filepath.Join(config.ImagePath, owner.Hex(), base)
		orginalPath = filepath.Join(config.ImagePath, owner.Hex(), name)
		orginalPath = oldstylefiles(orginalPath)
	}
	if _, err := os.Stat(fullpath); err == nil {
		return fullpath
	}
	if err := makeFile(orginalPath, fullpath); err != nil {
		log.Println(err)
	}
	return fullpath
}

func Count(owner bson.ObjectId) int {
	count, _ := db.Where(bson.M{
		"owner": owner,
	}).Count(&File{})
	return count
}

func LoadFiles(owner bson.ObjectId, limit, page int) (files []File) {
	db.Where(bson.M{"owner": owner}).
		Sort("-created_at").Paginate(limit, page).Find(&files)
	return files
}

func LoadFilesForServe(owner bson.ObjectId, limit, page int) []map[string]interface{} {
	res := []map[string]interface{}{}
	files := LoadFiles(owner, limit, page)
	for _, f := range files {
		res = append(res, f.Rest())
	}
	return res
}

func Load(query bson.M, limit, page int) (files []File) {
	db.Where(query).Sort("-created_at").Paginate(limit, page).Find(&files)
	return files
}

func (f *File) Delete() error {
	if err := db.Collection(f).Remove(bson.M{"id": f.ID}); err != nil {
		return err
	}
	path := filepath.Join(config.ImagePath, f.Path)
	return os.Remove(path)
}

func getFileOrginalName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	base = strings.Replace(base, ext, "", -1)
	splited := strings.Split(base, "*")
	return splited[0] + ext
}

func oldstylefiles(path string) string {
	if _, err := os.Stat(path); err == nil {
		return path
	}
	ext := filepath.Ext(path)
	path = strings.Replace(path, ext, "", -1)
	path += "_big" + ext
	return path
}
