package core

import "os"

type LruItem[T any] struct {
	data *T
	next *LruItem[T]
	prev *LruItem[T]
}

type FileWithDate struct {
	FileName string
	Date     int
}

type DatedSource[T any] interface {
	GetFileDate(fileName string, folderName string) (int, error)
	Load(files []FileWithDate) (*T, error)
}

type TimeSeriesDataRange[T any] struct {
	Idx  int
	Data *T
}

type TimeSeriesData[T any] struct {
	dataFolderPath  string
	source          DatedSource[T]
	indexCalculator func(int) int
	data            []*LruItem[T]
	maxIndex        int
	head            *LruItem[T]
	tail            *LruItem[T]
}

func LoadTimeSeriesData[T any](dataFolderPath string, source DatedSource[T], capacity int,
	indexCalculator func(int) int) (TimeSeriesData[T], error) {
	data := TimeSeriesData[T]{dataFolderPath: dataFolderPath, source: source, maxIndex: -1, indexCalculator: indexCalculator,
		data: make([]*LruItem[T], capacity), head: nil, tail: nil}
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
		data.Add(k, item)
	}
	return data, nil
}

func (t *TimeSeriesData[T]) Add(k int, item *T) {
	if k > t.maxIndex {
		t.maxIndex = k
	}
	i := &LruItem[T]{item, nil, t.head}
	t.data[k] = i
	if t.head != nil {
		t.head.next = i
	} else {
		t.tail = i
	}
	t.head = i
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
	idx1 := t.indexCalculator(from)
	if idx1 < 0 {
		idx1 = 0
	}
	idx2 := t.indexCalculator(to)
	var result []TimeSeriesDataRange[T]
	for i := idx1; i <= idx2; i++ {
		if i > t.maxIndex {
			break
		}
		d := t.data[i]
		if d != nil {
			result = append(result, TimeSeriesDataRange[T]{i, d.data})
		}
	}
	return result, nil
}

func (t *TimeSeriesData[T]) Get(date int) (*T, error) {
	idx := t.indexCalculator(date)
	if idx < 0 || idx > t.maxIndex {
		return nil, nil
	}
	for i := idx; i >= 0; idx-- {
		d := t.data[i]
		if d != nil {
			return d.data, nil
		}
	}
	return nil, nil
}
