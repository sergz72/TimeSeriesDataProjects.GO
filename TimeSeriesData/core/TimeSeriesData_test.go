package core

import "testing"

type testData struct{}

type testDatedSource struct{}

func (t testDatedSource) GetFileDate(fileName string, folderName string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (t testDatedSource) Load(files []FileWithDate) (*testData, error) {
	//TODO implement me
	panic("implement me")
}

func (t testDatedSource) GetFiles(date int, dataFolderPath string) ([]FileWithDate, error) {
	//TODO implement me
	panic("implement me")
}

func (t testDatedSource) Save(date int, data *testData, dataFolderPath string) error {
	//TODO implement me
	panic("implement me")
}

func TestLruList(t *testing.T) {
	data := NewTimeSeriesData[testData]("", testDatedSource{}, 500,
		func(date int) int { return date }, 500)
	for i := 0; i < 3; i++ {
		_ = data.Add(i, i, &testData{})
	}
	if data.lruManager.head.Key != 2 {
		t.Fatal("head should be 2")
	}
	if data.lruManager.head.prev != nil {
		t.Fatal("head prev should be nil")
	}
	item := data.lruManager.head.next
	if item.Key != 1 {
		t.Fatal("head next key should be 1")
	}
	if data.lruManager.tail.Key != 0 {
		t.Fatal("tail should be 0")
	}
	if data.lruManager.tail.prev != item {
		t.Fatal("data.lruManager.tail.prev != data.lruManager.head.next")
	}
	if data.lruManager.tail.next != nil {
		t.Fatal("tail next should be nil")
	}
	if item.prev != data.lruManager.head {
		t.Fatal("item.prev != data.lruManager.head")
	}
	if item.next != data.lruManager.tail {
		t.Fatal("item.next != data.lruManager.tail")
	}
}

func TestLruExpireAndMoveToFront(t *testing.T) {
	data := NewTimeSeriesData[testData]("", testDatedSource{}, 2000,
		func(date int) int { return date }, 500)
	for i := 0; i < 1000; i++ {
		err := data.Add(i, i, &testData{})
		if err != nil {
			t.Fatal(err)
		}
	}
	if data.lruManager.activeItems != 500 {
		t.Fatal("data.lruManager.activeItems should be 500")
	}
	if data.lruManager.head.Key != 999 {
		t.Fatal("head should be 999")
	}
	if data.lruManager.head.prev != nil {
		t.Fatal("head prev should be nil")
	}
	if data.lruManager.head.next.Key != 998 {
		t.Fatal("head next key should be 998")
	}
	if data.lruManager.tail.Key != 500 {
		t.Fatal("tail should be 500")
	}
	if data.lruManager.tail.prev.Key != 501 {
		t.Fatal("tail prev key should be 501")
	}
	if data.lruManager.tail.next != nil {
		t.Fatal("tail next should be nil")
	}

	key, _, err := data.Get(501)
	if err != nil {
		t.Fatal(err)
	}
	if key != 501 {
		t.Fatal("key should be 501")
	}
	if data.lruManager.head.Key != 501 {
		t.Fatal("head should be 501")
	}
	if data.lruManager.head.prev != nil {
		t.Fatal("head prev should be nil")
	}
	if data.lruManager.head.next.Key != 999 {
		t.Fatal("head next key should be 999")
	}
}

func TestIterator(t *testing.T) {
	data := NewTimeSeriesData[testData]("", testDatedSource{}, 500,
		func(date int) int { return date }, 500)
	for i := 1; i <= 5; i += 2 {
		_ = data.Add(i, i, &testData{})
	}
	iter, err := data.Iterator(1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if !iter.HasNext() {
		t.Fatal("HasNext1")
	}
	idx, _, err := iter.Next()
	if err != nil {
		t.Fatal(err)
	}
	if idx != 1 {
		t.Fatal("Next1")
	}

	if !iter.HasNext() {
		t.Fatal("HasNext3")
	}
	idx, _, err = iter.Next()
	if err != nil {
		t.Fatal(err)
	}
	if idx != 3 {
		t.Fatal("Next3")
	}

	if !iter.HasNext() {
		t.Fatal("HasNext5")
	}
	idx, _, err = iter.Next()
	if err != nil {
		t.Fatal(err)
	}
	if idx != 5 {
		t.Fatal("Next5")
	}

	if iter.HasNext() {
		t.Fatal("HasNext7")
	}
}
