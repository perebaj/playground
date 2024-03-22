package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type File struct {
	Name       string
	Size       int
	UploadedAt time.Time
	TTL        *int
}

type Files []File
type Files2 []File

func (files Files2) Len() int {
	return len(files)
}

func (files Files2) Less(i, j int) bool {
	if files[i].Size != files[j].Size {
		return files[i].Size < files[j].Size
	}
	return files[i].Name < files[j].Name

}

func (files Files2) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}

func (files Files) Len() int {
	return len(files)
}

func (files Files) Less(i, j int) bool {
	if files[i].Size != files[j].Size {
		return files[i].Size < files[j].Size
	}
	return files[i].Name < files[j].Name

}

func (files Files) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}

type FileSystem struct {
	fs map[string]File
}

func NewFileSystem() FileSystem {
	return FileSystem{
		fs: make(map[string]File),
	}
}

func (f *FileSystem) Upload(name string, size int) string {
	_, ok := f.fs[name]

	if !ok {
		newFile := File{
			Name: name,
			Size: size,
		}
		f.fs[name] = newFile
		return "file created"
	}
	return "file already exists"
}

func (f *FileSystem) UploadAt(uploadAt string, name string, size int, ttl *int) string {
	_, ok := f.fs[name]
	//convert string to time
	timeParsed, err := time.Parse(time.RFC3339, uploadAt)
	if err != nil {
		panic(err)
	}
	if !ok {
		newFile := File{
			Name:       name,
			Size:       size,
			UploadedAt: timeParsed,
			TTL:        ttl,
		}
		f.fs[name] = newFile
		return "file created"
	}
	return "file already exists"
}

func (f *FileSystem) GetAt(getAt string, name string) string {
	timeParsed, err := time.Parse(time.RFC3339, getAt)
	if err != nil {
		fmt.Println(err)
	}
	_, ok := f.fs[name]

	if ok {
		ttl := f.fs[name].TTL
		if ttl == nil {
			return fmt.Sprintf("file %s exists", name)
		} else {
			if timeParsed.Sub(f.fs[name].UploadedAt).Seconds() > float64(*ttl) {
				return "file doesn't exists. TTL expired"
			}
			return fmt.Sprintf("file %s exists", name)
		}
	}

	return fmt.Sprintf("file %s doesn't exists", name)
}

func (f *FileSystem) CopyAt(copiedAt string, source string, dest string) string {
	copiedAtTime, err := time.Parse(time.RFC3339, copiedAt)
	if err != nil {
		panic(err)
	}
	_, ok := f.fs[source]
	if !ok { // source not exists
		return "the source file doesn't exists"
	}

	sourceFile := f.fs[source]
	sourceFile.UploadedAt = copiedAtTime
	sourceFile.Name = dest
	f.fs[dest] = sourceFile

	return fmt.Sprintf("file %s copied to %s", source, dest)
}

func (f *FileSystem) Get(name string) int {
	_, ok := f.fs[name]

	if ok {
		return f.fs[name].Size
	}
	return 0
}

func (f *FileSystem) Copy(source string, dest string) string {
	_, ok := f.fs[source]
	if !ok { // source not exists
		return "the source file doesn't exists"
	}

	sourceFile := f.fs[source]
	sourceFile.Name = dest

	f.fs[dest] = sourceFile

	return "created"
}

func (f *FileSystem) SearchAt(timestamp string, prefix string) []string {
	var result Files
	for k, v := range f.fs {
		if strings.HasPrefix(k, prefix) {
			result = append(result, v)
		}
	}

	sort.Sort(result)
	var r []string

	timestampTime, err := time.Parse(time.RFC3339, timestamp)
	fmt.Println(timestampTime)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(result); i++ {
		if isValidTime(timestampTime, result[i]) {
			r = append(r, result[i].Name)
		}
	}

	if len(r) > 10 {
		return r[:10]
	}

	return r
}

func isValidTime(currentTime time.Time, file File) bool {
	if file.TTL == nil {
		return true
	} else if currentTime.Sub(file.UploadedAt).Seconds() <= float64(*file.TTL) {
		return true
	} else {
		return false
	}
}

func (f *FileSystem) Search(prefix string) []string {
	var result Files
	for k, v := range f.fs {
		if strings.HasPrefix(k, prefix) {
			result = append(result, v)
		}
	}

	sort.Sort(result)
	var r []string

	for i := 0; i < len(result); i++ {
		r = append(r, result[i].Name)
	}

	if len(r) > 10 {
		return r[:10]
	}

	return r
}

func main() {
	newObj := NewFileSystem()

	// fmt.Println(newObj.Upload("jojo.txt", 200))
	// fmt.Println(newObj.Upload("jojo.txt", 200))
	// fmt.Println(newObj.Get("jojo.txt"))
	// fmt.Println(newObj.Copy("jojo.txt", "mamarlon.txt"))
	// fmt.Println(newObj.Get("mamarlon.txt"))
	// fmt.Println(newObj.Upload("jojo.txt", 200))
	// fmt.Println(newObj.Upload("jojoa.txt", 300))
	// fmt.Println(newObj.Upload("jojob.txt", 300))
	// fmt.Println(newObj.Upload("jojo12.txt", 500))
	// fmt.Println(newObj.Upload("jojo5.txt", 600))
	// new uplaod at using rfc3339
	ttl := 3600
	fmt.Println(newObj.UploadAt("2023-05-02T09:34:01Z", "jojo5.txt", 600, &ttl))
	fmt.Println(newObj.GetAt("2023-05-02T09:55:01Z", "jojo5.txt"))
	fmt.Println(newObj.CopyAt("2023-05-03T09:56:01Z", "jojo5.txt", "jojo6.txt"))
	fmt.Println(newObj.GetAt("2023-05-02T09:55:01Z", "jojo6.txt"))

	fmt.Println(newObj.fs)
	// fmt.Println(newObj.GetAt("2022-01-01T00:00:00Z", "jojo5.txt"))
	// fmt.Println(newObj.Search("jojo"))
	fmt.Println(newObj.SearchAt("2023-05-02T09:55:01Z", "jojo"))

	// n := time.Now()
	// // add 1 hours
	// n2 := n.Add(time.Hour * 1)

	// fmt.Println(n2.Sub(n))
	// ttl := 7600 // 1 hour
	// // verify if n2.Sub(n) is greater than ttl
	// fmt.Println(n2.Sub(n).Seconds())
	// if n2.Sub(n).Seconds() >= float64(ttl) {
	// 	fmt.Println("greater")
	// }
}
