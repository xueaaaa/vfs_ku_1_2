package builtin

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"vfs/src/internal/flags"
	"vfs/src/internal/session"
	"vfs/src/internal/vfs"
)

const (
	lsName = "ls"

	flagLong      = 'l'
	flagAll       = 'a'
	flagHuman     = 'h'
	flagRecursive = 'R'
	lsFlags       = "lahR"

	hiddenPrefix = "."
	headerSuffix = ":"
	nameSep      = " "
	columnSep    = " "
	dirMark      = "d"
	fileMark     = "-"

	sizeBase      = 1024
	sizeUnits     = "BKMGT"
	sizePrecision = 1
	sizeFormat    = 'f'
	floatBits     = 64
	sizeZeroFrac  = ".0"

	msgEmptyDir  = "No files or directories in %s"
	msgNoSuchEnt = "no such file or directory: %s"
)

type entry struct {
	name  string
	isDir bool
	size  int
}

type lsContext struct {
	s       *session.Session
	out     Console
	opts    flags.Set
	headers bool
	printed bool
}

func newLs(s *session.Session, out Console) func([]string) (any, error) {
	return func(args []string) (any, error) {
		opts, operands, err := flags.Parse(args, lsFlags)
		if err != nil {
			return nil, err
		}
		if len(operands) == 0 {
			operands = []string{session.CurrentDir}
		}
		ctx := &lsContext{
			s:       s,
			out:     out,
			opts:    opts,
			headers: len(operands) > 1 || opts.Has(flagRecursive),
		}
		return nil, ctx.listAll(operands)
	}
}

func (c *lsContext) listAll(operands []string) error {
	var errs []error
	for _, operand := range operands {
		if err := c.listOperand(operand); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (c *lsContext) listOperand(operand string) error {
	if dir, err := c.s.ResolveDir(operand); err == nil {
		c.listTree(dir, operand)
		return nil
	}
	file, err := c.s.FindFile(operand)
	if err != nil {
		return fmt.Errorf(msgNoSuchEnt, operand)
	}
	c.out.Print(c.render([]entry{fileEntry(*file)})...)
	return nil
}

func (c *lsContext) listTree(dir *vfs.Directory, label string) {
	c.listDir(dir, label)
	if !c.opts.Has(flagRecursive) {
		return
	}
	for i := range dir.Subdirs {
		sub := &dir.Subdirs[i]
		if c.visible(sub.Name) {
			c.listTree(sub, joinLabel(label, sub.Name))
		}
	}
}

func (c *lsContext) listDir(dir *vfs.Directory, label string) {
	if c.headers {
		if c.printed {
			c.out.Print("")
		}
		c.out.Print(label + headerSuffix)
	}
	c.printed = true

	entries := c.entries(dir)
	if len(entries) == 0 {
		c.out.Print(fmt.Sprintf(msgEmptyDir, dir.Name))
		return
	}
	c.out.Print(c.render(entries)...)
}

func (c *lsContext) entries(dir *vfs.Directory) []entry {
	var list []entry
	if c.opts.Has(flagAll) {
		list = append(list,
			entry{name: session.CurrentDir, isDir: true},
			entry{name: session.ParentDir, isDir: true})
	}
	for _, d := range dir.Subdirs {
		if c.visible(d.Name) {
			list = append(list, entry{name: d.Name, isDir: true})
		}
	}
	for _, f := range dir.Files {
		if c.visible(f.Name) {
			list = append(list, fileEntry(f))
		}
	}
	return list
}

func (c *lsContext) visible(name string) bool {
	return c.opts.Has(flagAll) || !strings.HasPrefix(name, hiddenPrefix)
}

func fileEntry(f vfs.File) entry {
	return entry{name: f.Name, size: len(f.Content)}
}

func (c *lsContext) render(entries []entry) []string {
	if !c.opts.Has(flagLong) {
		return []string{strings.Join(names(entries), nameSep)}
	}
	return longLines(entries, c.opts.Has(flagHuman))
}

func names(entries []entry) []string {
	result := make([]string, 0, len(entries))
	for _, e := range entries {
		result = append(result, e.name)
	}
	return result
}

func longLines(entries []entry, human bool) []string {
	sizes := make([]string, len(entries))
	for i, e := range entries {
		sizes[i] = formatSize(e.size, human)
	}
	width := widest(sizes)

	lines := make([]string, len(entries))
	for i, e := range entries {
		cells := []string{mark(e), fmt.Sprintf("%*s", width, sizes[i]), e.name}
		lines[i] = strings.Join(cells, columnSep)
	}
	return lines
}

func mark(e entry) string {
	if e.isDir {
		return dirMark
	}
	return fileMark
}

func widest(cells []string) int {
	width := 0
	for _, cell := range cells {
		width = max(width, len(cell))
	}
	return width
}

func formatSize(size int, human bool) string {
	if !human {
		return strconv.Itoa(size)
	}
	return humanSize(size)
}

func humanSize(size int) string {
	value, unit := float64(size), 0
	for value >= sizeBase && unit < len(sizeUnits)-1 {
		value /= sizeBase
		unit++
	}
	if unit == 0 {
		return strconv.Itoa(size)
	}
	text := strconv.FormatFloat(value, sizeFormat, sizePrecision, floatBits)
	return strings.TrimSuffix(text, sizeZeroFrac) + string(sizeUnits[unit])
}

func joinLabel(label, name string) string {
	return strings.TrimSuffix(label, session.PathSep) + session.PathSep + name
}
