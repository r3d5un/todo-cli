package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// item struct represents a ToDo item
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// List represents a list of ToDo items
type List []item

func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

// Add creates a new ToDo item and appends it to the list
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Complete method marks a ToDo item as completed by
// setting Done = true and CompletedAt to the current time
func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}

	*l = append(ls[:i-1], ls[i:]...)

	return nil
}

// Delete method delets a ToDo itme from the list
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

// Save method encodes the List as JSON and saves it
// using the provided file name
func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, l)
}

func (l *List) String() string {
	formatted := ""

	for k, t := range *l {
		prefix := " "
		if t.Done {
			prefix = "X "
		}
		formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)
	}
	return formatted
}

func (l *List) List(incomplete bool, verbose bool) {
	formatted := ""

	for k, t := range *l {
		if incomplete {
			if t.Done {
				continue
			}
		}
		prefix := " "
		if t.Done {
			prefix = "X "
		}
		if verbose {
			switch t.Done {
			case true:
				formatted += fmt.Sprintf(
					"%s%d: %s - created %v - completed %v\n",
					prefix, k+1, t.Task, t.CreatedAt.Format("2020-01-01 23:59"), t.CompletedAt.Format("2020-01-01 23:59"),
				)
			default:
				formatted += fmt.Sprintf(
					"%s%d: %s - created %v\n",
					prefix, k+1, t.Task, t.CreatedAt.Format("2020-01-01 23:59"),
				)
			}
		} else {
			formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)
		}
	}
	fmt.Printf(formatted)
}
