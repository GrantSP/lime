package primitives

import (
	"sort"
	"sync"
)

type RegionSet struct {
	regions []Region
	lock    sync.Mutex
}

func (r *RegionSet) Adjust(position, delta int) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for i := range r.regions {
		if r.regions[i].A >= position {
			r.regions[i].A += delta
		} else if diff := position + delta - r.regions[i].A; diff < 0 {
			r.regions[i].A += diff
		}
		if r.regions[i].B >= position {
			r.regions[i].B += delta
		} else if diff := position + delta - r.regions[i].B; diff < 0 {
			r.regions[i].B += diff
		}
	}
	r.flush()
}

// TODO(q): There should be a on modified callback on the RegionSet
func (r *RegionSet) flush() {
	sort.Sort(r)
	for i := 1; i < len(r.regions); i++ {
		if r.regions[i-1] == r.regions[i] || r.regions[i-1].Intersects(r.regions[i]) || r.regions[i].Covers(r.regions[i-1]) {
			r.regions[i-1] = r.regions[i-1].Cover(r.regions[i])
			copy(r.regions[i:], r.regions[i+1:])
			r.regions = r.regions[:len(r.regions)-1]
		} else {
			i++
		}
	}
}

func (r *RegionSet) Substract(r2 Region) {
	r.Adjust(r2.Begin(), r2.Size())
}

func (r *RegionSet) Add(r2 Region) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.regions = append(r.regions, r2)
	r.flush()
}

func (r *RegionSet) Clear() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.regions = r.regions[0:0]
	r.flush()
}

func (r *RegionSet) Get(i int) Region {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.regions[i]
}

func (r *RegionSet) Less(i, j int) bool {
	if bi, bj := r.regions[i].Begin(), r.regions[j].Begin(); bi < bj {
		return true
	} else if bi == bj {
		return r.regions[i].End() < r.regions[j].End()
	}
	return false
}

func (r *RegionSet) Swap(i, j int) {
	r.regions[i], r.regions[j] = r.regions[j], r.regions[i]
}

func (r *RegionSet) Len() int {
	return len(r.regions)
}

func (r *RegionSet) AddAll(rs []Region) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.regions = append(r.regions, rs...)
	r.flush()
}

func (r *RegionSet) Contains(r2 Region) bool {
	for i := range r.regions {
		if r.regions[i] == r2 || (r.regions[i].Contains(r2.Begin()) && r.regions[i].Contains(r2.End())) {
			return true
		}
	}
	return false
}

func (r *RegionSet) Regions() (ret []Region) {
	r.lock.Lock()
	defer r.lock.Unlock()
	ret = make([]Region, len(r.regions))
	copy(ret, r.regions)
	return
}

func (r *RegionSet) HasNonEmpty() bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, r := range r.regions {
		if !r.Empty() {
			return true
		}
	}
	return false
}

func (r *RegionSet) HasEmpty() bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, r := range r.regions {
		if r.Empty() {
			return true
		}
	}
	return false
}
