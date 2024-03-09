package core

import "os"

type FileWithDate struct {
	FileName string
	Date     int
}

type DatedSource[T any] interface {
	GetFileDate(fileName string, folderName string) (int, error)
	Load(files []FileWithDate) (*T, error)
	GetFiles(date int) ([]FileWithDate, error)
	Save(date int, data *T) error
}

type TimeSeriesDataRange[T any] struct {
	Idx  int
	Data *T
}

type TimeSeriesData[T any] struct {
	dataFolderPath  string
	source          DatedSource[T]
	IndexCalculator func(int) int
	DateCalculator  func(int) int
	data            []*LruItem[T]
	maxIndex        int
	maxActiveItems  int
	lruManager      LruManager[T]
	modified        map[int]bool
}

func NewTimeSeriesData[T any](
	dataFolderPath string,
	source DatedSource[T],
	capacity int,
	indexCalculator func(int) int,
	dateCalculator func(int) int,
	maxActiveItems int) TimeSeriesData[T] {
	return TimeSeriesData[T]{dataFolderPath: dataFolderPath, source: source, maxIndex: -1, IndexCalculator: indexCalculator,
		DateCalculator: dateCalculator, data: make([]*LruItem[T], capacity), lruManager: LruManager[T]{},
		maxActiveItems: maxActiveItems, modified: make(map[int]bool)}
}

func LoadTimeSeriesData[T any](
	dataFolderPath string,
	source DatedSource[T],
	capacity int,
	indexCalculator func(int) int,
	dateCalculator func(int) int,
	maxActiveItems int) (TimeSeriesData[T], error) {
	data := NewTimeSeriesData(dataFolderPath, source, capacity, indexCalculator, dateCalculator, maxActiveItems)
	files, err := data.getFileList("")
	if err != nil {
		return data, err
	}
	fileMap := make(map[int][]FileWithDate)
	for _, f := range files {
		idx := indexCalculator(f.Date)
		fileList, ok := fileMap[idx]
		if !ok {
			fileMap[idx] = []FileWithDate{f}
		} else {
			fileList = append(fileList, f)
			fileMap[idx] = fileList
		}
	}
	for k, v := range fileMap {
		var item *T
		item, err = source.Load(v)
		if err != nil {
			return data, err
		}
		err = data.Add(k, item)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}

func InitTimeSeriesData[T any](
	dataFolderPath string,
	source DatedSource[T],
	capacity int,
	indexCalculator func(int) int,
	dateCalculator func(int) int,
	maxActiveItems int) (TimeSeriesData[T], error) {
	data := NewTimeSeriesData(dataFolderPath, source, capacity, indexCalculator, dateCalculator, maxActiveItems)
	files, err := data.getFileList("")
	if err != nil {
		return data, err
	}
	indexes := make(map[int]bool)
	for _, f := range files {
		idx := indexCalculator(f.Date)
		indexes[idx] = true
	}
	for idx := range indexes {
		data.set(idx, &LruItem[T]{})
	}
	return data, nil
}

func (t *TimeSeriesData[T]) set(k int, item *LruItem[T]) {
	if k > t.maxIndex {
		t.maxIndex = k
	}
	t.data[k] = item
}

func (t *TimeSeriesData[T]) Add(k int, item *T) error {
	err := t.cleanup()
	if err != nil {
		return err
	}
	t.set(k, t.lruManager.Add(k, item))
	return nil
}

func (t *TimeSeriesData[T]) getFileList(folder string) ([]FileWithDate, error) {
	baseFolder := t.dataFolderPath + "/" + folder
	files, err := os.Open(baseFolder)
	if err != nil {
		return nil, err
	}
	defer func() { _ = files.Close() }()
	fi, err := files.Readdir(-1)
	if err != nil {
		return nil, err
	}
	var result []FileWithDate
	for _, file := range fi {
		if file.Mode().IsDir() {
			var info []FileWithDate
			info, err = t.getFileList(file.Name())
			if err != nil {
				return nil, err
			}
			result = append(result, info...)
		} else {
			var date int
			date, err = t.source.GetFileDate(file.Name(), folder)
			if err != nil {
				return nil, err
			}
			result = append(result, FileWithDate{FileName: baseFolder + "/" + file.Name(), Date: date})
		}
	}
	return result, nil
}

func (t *TimeSeriesData[T]) GetRange(from int, to int) ([]TimeSeriesDataRange[T], error) {
	idx1 := t.IndexCalculator(from)
	if idx1 < 0 {
		idx1 = 0
	}
	idx2 := t.IndexCalculator(to)
	var result []TimeSeriesDataRange[T]
	for i := idx1; i <= idx2; i++ {
		if i > t.maxIndex {
			break
		}
		d := t.data[i]
		if d != nil {
			item, err := t.get(d)
			if err != nil {
				return nil, err
			}
			result = append(result, TimeSeriesDataRange[T]{i, item})
		}
	}
	return result, nil
}

func (t *TimeSeriesData[T]) Get(date int) (int, *T, error) {
	idx := t.IndexCalculator(date)
	if idx < 0 {
		return idx, nil, nil
	}
	if idx > t.maxIndex {
		idx = t.maxIndex
	}
	for i := idx; i >= 0; idx-- {
		d := t.data[i]
		if d != nil {
			dd, err := t.get(d)
			return i, dd, err
		}
	}
	return idx, nil, nil
}

func (t *TimeSeriesData[T]) get(item *LruItem[T]) (*T, error) {
	if item.Data != nil {
		t.lruManager.MoveToFront(item)
		return item.Data, nil
	}
	err := t.cleanup()
	if err != nil {
		return nil, err
	}
	date := t.DateCalculator(item.Key)
	files, err := t.source.GetFiles(date)
	if err != nil {
		return nil, err
	}
	var data *T
	data, err = t.source.Load(files)
	if err != nil {
		return nil, err
	}
	item.Data = data
	t.lruManager.Attach(item)

	return item.Data, nil
}

func (t *TimeSeriesData[T]) cleanup() error {
	for t.lruManager.activeItems >= t.maxActiveItems {
		err := t.removeByLru()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TimeSeriesData[T]) removeByLru() error {
	tail := t.lruManager.GetTail()
	if t.modified[tail.Key] {
		date := t.DateCalculator(tail.Key)
		err := t.source.Save(date, tail.Data)
		if err != nil {
			return err
		}
		delete(t.modified, tail.Key)
	}
	tail.Data = nil
	t.lruManager.Detach(tail)
	return nil
}
