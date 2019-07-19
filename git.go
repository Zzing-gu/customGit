package git

import (
	"fmt"
	"io"
	"os"

	"bytes"
	//"reflect"
	"strconv"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/filemode"
	"gopkg.in/src-d/go-git.v4/plumbing/format/index"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"

	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
	mindex "gopkg.in/src-d/go-git.v4/utils/merkletrie/index"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie/noder"
)

var emptyNoderHash = make([]byte, 24)

func diffTreeIsEquals(a, b noder.Hasher) bool {
	hashA := a.Hash()
	hashB := b.Hash()

	if bytes.Equal(hashA, emptyNoderHash) || bytes.Equal(hashB, emptyNoderHash) {
		return false
	}

	return bytes.Equal(hashA, hashB)
}

func GetIndex() {
	s := filesystem.NewStorage(osfs.New("nexivil6/hosuk6"), cache.NewObjectLRUDefault())

	repo, err := git.Open(s, osfs.New("nexivil6/hosuk6"))
	if err != nil {
		panic(err)
	}

	idx, err := repo.Storer.Index()
	if err != nil {
		panic(err)
	}

	str := idx.String()

	fmt.Println(str)

	//fmt.Println(idx.Entries[0])
}

func Getindexoftree(hash string) error {

	s := filesystem.NewStorage(osfs.New("nexivil6/hosuk6"), cache.NewObjectLRUDefault())

	repo, err := git.Open(s, osfs.New("nexivil6/hosuk6"))
	if err != nil {
		panic(err)
	}

	h := plumbing.NewHash(hash)

	idx, err := repo.Storer.Index()
	if err != nil {
		return err
	}

	// commit hash
	// c, err := repo.CommitObject(h)
	// if err != nil {
	// 	panic(err)
	// }

	// t, err := c.Tree()
	// if err != nil {
	// 	panic(err)
	// }

	// tree hash
	t, err := object.GetTree(s, h)
	if err != nil {
		panic(err)
	}

	var from noder.Noder
	if t != nil {
		from = object.NewTreeRootNode(t)
	}

	to := mindex.NewRootNode(idx)

	// if reverse {
	// 	return merkletrie.DiffTree(to, from, diffTreeIsEquals)
	// }
	// // reverse 없는 관계로 이게 맞음 ....
	// return merkletrie.DiffTree(to, from, diffTreeIsEquals)

	//merkletrie.Changes, error
	//return merkletrie.DiffTree(from, to, diffTreeIsEquals)

	changes, err := merkletrie.DiffTree(to, from, diffTreeIsEquals)
	if err != nil {
		return err
	}

	for _, ch := range changes {
		a, err := ch.Action()
		if err != nil {
			return err
		}

		var name string
		var e *object.TreeEntry

		switch a {
		case merkletrie.Modify, merkletrie.Insert:
			name = ch.To.String()
			e, err = t.FindEntry(name)
			if err != nil {
				return err
			}
		case merkletrie.Delete:
			name = ch.From.String()
		}

		_, _ = idx.Remove(name)
		if e == nil {
			continue
		}

		idx.Entries = append(idx.Entries, &index.Entry{
			Name: name,
			Hash: e.Hash,
			Mode: e.Mode,
		})

	}

	return repo.Storer.SetIndex(idx)
}

func GetRepoTree() ([]object.TreeEntry, error) {
	fmt.Printf("GetHeadTree")

	s := filesystem.NewStorage(osfs.New("nexivil6/hosuk6"), cache.NewObjectLRUDefault())

	repo, err := git.Open(s, osfs.New("nexivil6/hosuk6"))
	if err != nil {
		panic(err)
	}

	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}

	h := ref.Hash()

	c, err := repo.CommitObject(h)
	if err != nil {
		panic(err)
	}

	t, err := c.Tree()
	if err != nil {
		panic(err)
	}
	fmt.Println(t)
	return t.Entries, nil
}

func RenderTree(hash string) ([]object.TreeEntry, error) {
	// 트리를 인자로 받고 트리를 렌더링 한다 .... 그안에 깊숙한 폴더를 들어갈때도 다시 이함수를 호출해준다.
	fmt.Printf("RenderTree")

	s := filesystem.NewStorage(osfs.New("nexivil6/hosuk6"), cache.NewObjectLRUDefault())

	h := plumbing.NewHash(hash)

	t, err := object.GetTree(s, h)
	if err != nil {
		panic(err)
	}
	fmt.Println(t.Entries)
	return t.Entries, err
}

func CreateAndInitDirectory(path string) {

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// folder exists
		fmt.Printf("already exist")
	} else {

		os.Mkdir(path, os.ModeDir)
		git.PlainInit(path, true)
	}
}

//path string, filedata []byte, filename string, filemode uint32
func AddOrUpdateFile(path string, filedata []byte, filename string, filetype uint32) {
	var u uint32 = filetype
	var ssf = strconv.FormatUint(uint64(u), 10)
	filemodenum, err := filemode.New(ssf)
	fmt.Println(filemodenum)
	if err != nil {
		panic(err)
	}
	s := filesystem.NewStorage(osfs.New(path), cache.NewObjectLRUDefault())

	repo, err := git.Open(s, osfs.New(path))
	if err != nil {
		panic(err)
	}

	w, err := repo.Worktree()
	if err != nil {
		panic(err)
	}

	idx, err := repo.Storer.Index()
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(filedata)

	obj := repo.Storer.NewEncodedObject()
	obj.SetType(plumbing.BlobObject)
	fmt.Println(reader.Len())
	fmt.Println(int64(reader.Len()))
	fmt.Println(uint32(int64(reader.Len())))
	obj.SetSize(int64(reader.Len()))

	tmpSize := uint32(int64(reader.Len()))
	writer, err := obj.Writer()
	if err != nil {
		panic(err)
	}
	// 이 프로세스 이후 reader size  가 날라가더라 ....
	if _, err := io.Copy(writer, reader); err != nil {
		panic(err)
	}

	h, err := repo.Storer.SetEncodedObject(obj)
	if err != nil {
		panic(err)
	}

	e, err := idx.Entry(filename)

	//add
	if err == index.ErrEntryNotFound {
		e = idx.Add(filename)
	}

	e.Hash = h
	e.ModifiedAt = time.Now()
	e.Mode = filemodenum

	//# 나중에 십진법 팔진법 꼬일수 있으니 조심하자 ...
	if e.Mode.IsRegular() {
		fmt.Println("it is regular 100644")
		e.Size = tmpSize
		fmt.Println(tmpSize)

		fmt.Println(e.Size)
	}

	repo.Storer.SetIndex(idx)
	//end

	commit, err := w.Commit("this is a last file test", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "hosuk",
			Email: "kirklayer@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		panic(err)
	}

	commitObj, err := repo.CommitObject(commit)
	if err != nil {
		panic(err)
	}
	fmt.Println(commitObj)

}
