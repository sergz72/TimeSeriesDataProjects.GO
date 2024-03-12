package core

import "os"

type FileWithDate struct {
	FileName string
	Date     int
}

type DatedSource[T any] interface {
	GetFileDate(fileName string, folderName string) (int, error)
	Load(files []FileWithDate) (*T, error)
	GetFiles(date int, dataFolderPath string) ([]FileWithDate, error)
	Save(date int, data *T, dataFolderPath string) error
}

type TimeSeriesData[T any] struct {
	dataFolderPath string
	source         DatedSource[T]
	// calculates array index from file date
	IndexCalculator func(int) int
	// calculates file date from date yyyymmdd
	DateCalculator func(int) int
	data           []*LruItem[T]
	maxIndex       int
	maxActiveItems int
	lruManager     LruManager[T]
	modified       map[int]bool
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
		date := dateCalculator(v[0].Date)
		err = data.Add(k, date, item)
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
	indexes := make(map[int]int)
	for _, f := range files {
		idx := indexCalculator(f.Date)
		date := dateCalculator(f.Date)
		indexes[idx] = date
	}
	for idx, date := range indexes {
		data.set(idx, NewLruItem[T](idx, date))
	}
	return data, nil
}

func (t *TimeSeriesData[T]) set(k int, item *LruItem[T]) {
	if k > t.maxIndex {
		t.maxIndex = k
	}
	t.data[k] = item
}

func (t *TimeSeriesData[T]) Add(k, date int, item *T) error {
	err := t.cleanup()
	if err != nil {
		return err
	}
	t.set(k, t.lruManager.Add(k, date, item))
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

func (t *TimeSeriesData[T]) GetExact(date int) (*T, error) {
	idx := t.IndexCalculator(date)
	if idx < 0 || idx > t.maxIndex {
		return nil, nil
	}
	d := t.data[idx]
	if d != nil {
		return t.get(d)
	}
	return nil, nil
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
	files, err := t.source.GetFiles(item.Date, t.dataFolderPath)
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
		err := t.source.Save(tail.Date, tail.Data, t.dataFolderPath)
		if err != nil {
			return err
		}
		delete(t.modified, tail.Key)
	}
	tail.Data = nil
	t.lruManager.Detach(tail)
	return nil
}

func (t *TimeSeriesData[T]) saveIndex(index int, source DatedSource[T], dataFolderPath string) error {
	d := t.data[index]
	if d != nil && d.Data != nil {
		err := source.Save(d.Date, d.Data, dataFolderPath)
		if err != nil {
			return err
		}
		delete(t.modified, index)
	}
	return nil
}

func (t *TimeSeriesData[T]) SaveAll(source DatedSource[T], dataFolderPath string) error {
	for i := 0; i <= t.maxIndex; i++ {
		err := t.saveIndex(i, source, dataFolderPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TimeSeriesData[T]) Save() error {
	var toDelete []int
	for idx := range t.modified {
		toDelete = append(toDelete, idx)
	}
	for _, idx := range toDelete {
		err := t.saveIndex(idx, t.source, t.dataFolderPath)
		if err != nil {
			return err
		}
	}
	return nil
}

type TimeSeriesDataIterator[T any] struct {
	data        *TimeSeriesData[T]
	current     int
	to          int
	currentData *T
}

func (i *TimeSeriesDataIterator[T]) HasNext() bool {
	return i.current <= i.to
}

func (i *TimeSeriesDataIterator[T]) seekToNext() error {
	i.current++
	for i.current <= i.to {
		d := i.data.data[i.current]
		if d != nil {
			var err error
			i.currentData, err = i.data.get(d)
			return err
		}
		i.current++
	}
	return nil
}

func (i *TimeSeriesDataIterator[T]) Next() (int, *T, error) {
	next := i.current
	data := i.currentData
	return next, data, i.seekToNext()
}

func (t *TimeSeriesData[T]) Iterator(from, to int) (*TimeSeriesDataIterator[T], error) {
	idx1 := t.IndexCalculator(from)
	if idx1 < 0 {
		idx1 = 0
	}
	idx2 := t.IndexCalculator(to)
	if idx2 > t.maxIndex {
		idx2 = t.maxIndex
	}
	i := TimeSeriesDataIterator[T]{t, idx1 - 1, idx2, nil}
	err := i.seekToNext()
	return &i, err
}

func (t *TimeSeriesData[T]) MarkAsModified(idx int) {
	t.modified[idx] = true
}

func (t *TimeSeriesData[T]) GetDate(key int) int {
	return t.data[key].Date
}