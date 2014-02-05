package set

type Set interface {
  Add(value string)
  Contains(value string) (bool)
  Length() (int)
	Set() ([]string)
  RemoveDuplicates()
}

type HashSet struct {
  data map[string]bool
}

func (this *HashSet) Set() ([]string) {
	var l []string
	for k, _ := range this.data {
		l = append(l, k)
	}
	return l  
}

func (this *HashSet) Add(value string) {
  this.data[value] = true
}

func (this *HashSet) Contains(value string) (exists bool) {
  _, exists = this.data[value]
  return
}

func (this *HashSet) Length() (int) {
  return len(this.data)
}
func (this *HashSet) RemoveDuplicates() {}


func NewSet() (Set) {
  return &HashSet{make(map[string] bool)}
}
