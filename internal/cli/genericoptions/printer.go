package genericoptions

import (
	"strings"

	"github.com/spf13/cobra"
)

// PrintFlags provides flags for configuring resource printing.
type PrinterFlags struct {
	// fields is a comma-separated string listing the desired fields to pluck
	// from a resource and print.
	fields string

	// sortBy contains the name of the field the resource list should be sorted
	// upon.
	sortBy string

	// descending determines whether or not the resource list is printed in
	// descending order.
	descending bool

	// noHeader determines whether or not the header should be omitted by the
	// resource printer.
	noHeader bool
}

// AddFlags recieves a *cobra.Command reference and binds flags for specifying
// fields to display and to sorty by.
func (f *PrinterFlags) AddFlags(c *cobra.Command) {
	if f == nil {
		return
	}

	c.Flags().StringVar(&f.fields, "fields", f.fields, "Comma-separated list of flags to output")
	c.Flags().StringVar(&f.sortBy, "sort-by", f.sortBy, "The field to sort by")
	c.Flags().BoolVar(&f.descending, "descending", false, "If specified, resources will be printed in descending order")
	c.Flags().BoolVar(&f.noHeader, "no-header", false, "If specified, the header will be omitted by the resource printer")
}

// Fields returns the fields for the resource printer to print.
func (f *PrinterFlags) Fields() []string {
	if f.fields == "" {
		return nil
	}
	return strings.Split(f.fields, ",")
}

// SortBy returns the field for the resource printer to sort upon.
func (f *PrinterFlags) SortBy() string {
	return f.sortBy
}

// Descending returns whether the resource printer should sort in descending order.
func (f *PrinterFlags) Descending() bool {
	return f.descending
}

// NoHeader returns whether the resource printer should omit the header.
func (f *PrinterFlags) NoHeader() bool {
	return f.noHeader
}
