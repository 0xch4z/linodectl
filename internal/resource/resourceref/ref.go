package resourceref

import (
	"errors"
	"strconv"
)

// Ref is a reference to a resource. This can either be an ID or Label.
type Ref struct {
	val any
}

type Meta struct {
	Label string
	ID    int
}

// ID returns the reference as a numeric ID or 0.
func (r Ref) ID() int {
	id, ok := r.val.(int)
	if !ok {
		return 0
	}
	return id
}

// Label returns the reference as a string.
func (r Ref) Label() string {
	s, ok := r.val.(string)
	if !ok {
		return ""
	}
	return s
}

type List []Ref

// Label returns a label if the List contains only a single label.
// Otherwise, it's an empty string.
func (l List) Label() string {
	if len(l) != 1 {
		return ""
	}
	return l[0].Label()
}

// ID returns a numeric ID if the List contains only an ID.
// Otherwise, 0 is returned.
func (l List) ID() int {
	if len(l) != 1 {
		return 0
	}
	return l[0].ID()
}

// Identifiers gets all the label and id identifiers from the list.
func (l List) Identifiers() (map[string]struct{}, map[int]struct{}) {
	labels := make(map[string]struct{})
	ids := make(map[int]struct{})
	for _, identifier := range l {
		if id := identifier.ID(); id != 0 {
			ids[id] = struct{}{}
		}
		if label := identifier.Label(); label != "" {
			labels[label] = struct{}{}
		}
	}
	return labels, ids
}

func newRef(v any) (*Ref, error) {
	switch val := v.(type) {
	case int:
		if val < 0 {
			return nil, errors.New("numeric identifiers must be positive integers")
		}
	case string:
		if val == "" {
			return nil, errors.New("labels can not be empty strings")
		}
	}
	return &Ref{val: v}, nil
}

// ListFromArgs parses a List of resource references from a list of arguments. If something other
// than a resource label or numeric ID is specified, an error will be thrown.
func ListFromArgs(args []string) (List, error) {
	refs := make(List, len(args))
	for i, arg := range args {
		// If successful, we have a valid integer that is maybe a valid numeric ID.
		n, err := strconv.Atoi(arg)
		if err == nil {
			ref, err := newRef(n)
			if err != nil {
				return nil, err
			}
			refs[i] = *ref
			continue
		}

		// This has to be a resource label.
		ref, err := newRef(arg)
		if err != nil {
			return nil, err
		}
		refs[i] = *ref
	}

	return refs, nil
}
