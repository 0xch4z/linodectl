package printer

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/Charliekenney23/linodectl/internal/resource"
	"github.com/Charliekenney23/linodectl/internal/strutil"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Printer struct {
	t table.Writer
}

func New(w io.Writer) Printer {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	return Printer{t}
}

func (p *Printer) setResourceHeader(columns []string) {
	cols := make([]interface{}, len(columns))
	for i, col := range columns {
		cols[i] = col
	}
	p.t.AppendHeader(cols)
}

func (p *Printer) normalizeColumns(r resource.Meta, columns []string) []string {
	requiredColumns := r.RequiredColumns()
	normalized := make([]string, 0, len(requiredColumns)+len(columns))
	seen := make(map[string]struct{}, len(requiredColumns))

	for _, column := range append(requiredColumns, columns...) {
		if _, ok := seen[column]; !ok {
			normalized = append(normalized, column)
			seen[column] = struct{}{}
		}
	}
	return normalized
}

type ResourcePrintOptions struct {
	// The columns to print
	Columns []string

	// SortBy is the column to sort by
	SortBy string

	// OmitHeader specifies whether the resource header should be omitted
	OmitHeader bool

	// DescendingOrder specifies whether the resource list should be printed in
	// descending order
	DescendingOrder bool
}

func (p *Printer) PrintResources(ctx context.Context, resList resource.List, options ResourcePrintOptions) error {
	items := resList.Items()
	if len(items) == 0 {
		return nil
	}

	if options.Columns == nil {
		options.Columns = resList.Meta().DefaultColumns()
	}

	options.Columns = p.normalizeColumns(resList.Meta(), options.Columns)

	if !options.OmitHeader {
		p.setResourceHeader(options.Columns)
	}

	for _, res := range items {
		if err := p.printResource(ctx, resList.Meta(), res, options.Columns); err != nil {
			return err
		}
	}

	_ = p.t.Render()
	return nil
}

func (p *Printer) printResource(ctx context.Context, meta resource.Meta, res resource.Resource, columns []string) error {
	values := make([]interface{}, len(columns))
	properties := res.Properties()
	for i, col := range columns {
		property, ok := properties[col]
		if !ok {
			return fmt.Errorf("property %q does not exist on resource %q", col, meta.Name())
		}

		var err error
		values[i], err = property.Getter(ctx)
		if err != nil {
			return fmt.Errorf("failed to resolve value %q from resource %q: %w", col, meta.Name(), err)
		}

		switch v := values[i].(type) {
		case []string:
			values[i] = strutil.Fallback(strings.Join(v, ", "), "<empty>")
		}
	}

	p.t.AppendRow(values)
	return nil
}
