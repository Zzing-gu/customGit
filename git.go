package git

import (
	"fmt"
	"io"
	"os"
	"io/ioutil"

	"bytes"
	//"reflect"
	"strconv"
	"time"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/filemode"
	"gopkg.in/src-d/go-git.v4/plumbing/format/index"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"


	"gopkg.in/src-d/go-git.v4/utils/merkletrie/noder"

	//"gopkg.in/src-d/go-git.v4/config"

	"gopkg.in/src-d/go-git.v4/plumbing/storer"
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


func MyTreeDiff(th1 string , th2 string) {
	
	s := filesystem.NewStorage(osfs.New("nexivil8/hosuk8"), cache.NewObjectLRUDefault())

	h1 := plumbing.NewHash(th1)

	

	// tree hash
	t1, err := object.GetTree(s, h1)
	if err != nil {
		panic(err)
	}


	h2 := plumbing.NewHash(th2)

	

	// tree hash
	t2, err := object.GetTree(s, h2)
	if err != nil {
		panic(err)
	}

	/////////////////

	changes , err := t1.Diff(t2)
	if err != nil {
		panic(err)
	}

	fmt.Println(changes)

	p, err := changes.Patch()
	if err != nil {
		panic(err)
	}
	//fmt.Println(p)
	//fmt.Println(p.Message())
	//fmt.Println(p.Stats())
	//fmt.Println(p.FilePatches()[0])
	str := p.String()
	fmt.Println(str)
}


func MyCheckOut( branch string , isCreate bool) {

	s := filesystem.NewStorage(osfs.New("nexivil8/hosuk8"), cache.NewObjectLRUDefault())

	repo, err := git.Open(s, osfs.New("nexivil8/hosuk8"))
	if err != nil {
		panic(err)
	}

	w, err := repo.Worktree()
	if err != nil {
		panic(err)
	}

	//var branchName plumbing.ReferenceName  = "refs/heads/"
	branchName := plumbing.ReferenceName(branch)

	err = w.Checkout(&git.CheckoutOptions{
		Create: isCreate,
		Branch: branchName,
	})
	if err != nil {
		// unstaged change 오류 .... 
		//panic(err)
	}

	head, err := repo.Head()
	if err != nil {
		panic(err)
	}
	fmt.Println(head)



	h := head.Hash()

	c, err := repo.CommitObject(h)
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
	
}




func GetBranches() []string {
	s := filesystem.NewStorage(osfs.New("nexivil8/hosuk8"), cache.NewObjectLRUDefault())

	repo, err := git.Open(s, osfs.New("nexivil8/hosuk8"))
	if err != nil {
		panic(err)
	}

	var branchArr []string

	branches, err := repo.Branches()
	if err != nil {
		panic(err)
	}
	//fmt.Println(refs)

	branches.ForEach( func(branch *plumbing.Reference) error {
		fmt.Println(branch.Name())
		branchArr = append(branchArr, string(branch.Name()))
		return nil
	})


	return branchArr
	
}



func GetLog(skip int, limit int) {
	s := filesystem.NewStorage(osfs.New("nexivil8/hosuk8"), cache.NewObjectLRUDefault())

	repo, err := git.Open(s, osfs.New("nexivil8/hosuk8"))
	if err != nil {
		panic(err)
	}
	// ... retrieves the commit history
	cIter, err := repo.Log(&git.LogOptions{All: false})
	if err != nil {
		panic(err)
	}

	//fmt.Println(cIter)
	// 여기서 어떤 조작으로 로그 문제 해결 .... 
	// ... just iterates over the commits, printing it
	count := 0
	err = cIter.ForEach(func(c *object.Commit) error {

		if count > skip+limit {
			fmt.Println("break it")
			return storer.ErrStop 
		}

		if count > skip {

			hash := c.Hash.String()
			line := strings.Split(c.Message, "\n")
			fmt.Println(hash[:7], line[0])
			fmt.Println(c)

		}

		count++
		return nil
	})



	if err != nil {
		panic(err)
	}
}


func GetCommitTree(hash string) {

}




func GetRepoTree() ( []object.TreeEntry, error) {
	fmt.Printf("GetHeadTree")

	s := filesystem.NewStorage(osfs.New("nexivil8/hosuk8"), cache.NewObjectLRUDefault())

	repo, err := git.Open(s, osfs.New("nexivil8/hosuk8"))
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

	s := filesystem.NewStorage(osfs.New("nexivil8/hosuk8"), cache.NewObjectLRUDefault())

	h := plumbing.NewHash(hash)

	t, err := object.GetTree(s, h)
	if err != nil {
		panic(err)
	}

	

	fmt.Println(t.Entries)
	return  t.Entries, err
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

	//idx.Remove("timetest.txt")

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







func AddOrUpdateFileTest() {
	//   path string, filedata []byte, filename string, filetype uint32

	var path string = "nexivil8/hosuk8"

	filedata, err := ioutil.ReadFile("test.txt")
	if err != nil {
		panic(err)
	}

	var filename string = "test.txt"

	var u uint32 = 100644
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

	//idx.Remove("timetest.txt")

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