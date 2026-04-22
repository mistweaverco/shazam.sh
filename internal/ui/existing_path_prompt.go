package ui

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/term"
)

type ExistingPathResolution string

const (
	ResolutionSkip      ExistingPathResolution = "skip"
	ResolutionOverwrite ExistingPathResolution = "overwrite"
	ResolutionAbort     ExistingPathResolution = "abort"
)

type ExistingPathPromptInput struct {
	Source      string
	Destination string
}

var (
	ttyOnce   sync.Once
	cachedTTY bool
)

func IsTTY() bool {
	ttyOnce.Do(func() {
		cachedTTY = term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd()))
	})
	return cachedTTY
}

func PromptExistingPathResolution(in ExistingPathPromptInput) (ExistingPathResolution, error) {
	if strings.TrimSpace(in.Destination) == "" {
		return ResolutionAbort, errors.New("destination is empty")
	}

	m := newExistingPathModel(in)
	p := tea.NewProgram(&m, tea.WithOutput(os.Stdout), tea.WithInput(os.Stdin))
	final, err := p.Run()
	if err != nil {
		return ResolutionAbort, err
	}

	fm, ok := final.(*existingPathModel)
	if !ok {
		return ResolutionAbort, errors.New("unexpected prompt model")
	}
	if fm.err != nil {
		return ResolutionAbort, fm.err
	}
	return fm.resolution, nil
}

type existingPathState int

const (
	stateSelect existingPathState = iota
	stateDiff
)

type existingPathModel struct {
	in ExistingPathPromptInput

	state      existingPathState
	resolution ExistingPathResolution
	quitting   bool
	err        error

	list     list.Model
	viewport viewport.Model
	diffText string
}

type existingPathChoice struct {
	key   ExistingPathResolution
	title string
	desc  string
}

func (c existingPathChoice) Title() string       { return c.title }
func (c existingPathChoice) Description() string { return c.desc }
func (c existingPathChoice) FilterValue() string { return c.title }

func newExistingPathModel(in ExistingPathPromptInput) existingPathModel {
	items := []list.Item{
		existingPathChoice{
			key:   ResolutionSkip,
			title: "Skip",
			desc:  "Keep the existing destination and continue.",
		},
		existingPathChoice{
			key:   ResolutionOverwrite,
			title: "Overwrite",
			desc:  "Delete the existing destination (recursively if directory) and recreate.",
		},
		existingPathChoice{
			key:   ExistingPathResolution("diff"),
			title: "Show diff",
			desc:  "Show a best-effort diff between destination and source.",
		},
		existingPathChoice{
			key:   ResolutionAbort,
			title: "Abort",
			desc:  "Stop the run.",
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Destination already exists — how should shazam resolve this?"
	l.SetShowHelp(true)
	l.DisableQuitKeybindings()

	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().Padding(0, 1)

	return existingPathModel{
		in:       in,
		state:    stateSelect,
		list:     l,
		viewport: vp,
	}
}

func (m *existingPathModel) Init() tea.Cmd {
	return nil
}

func (m *existingPathModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Keep the prompt compact (avoid "full screen" feel).
		// Still adapt down for small terminals.
		w := clamp(msg.Width-2, 40, 84)
		// Height is fluid, but capped so it feels like a modal.
		// Target: up to ~2/3 of terminal height, with a hard max.
		maxH := min(22, max(10, (msg.Height*2)/3))
		h := clamp(msg.Height-4, 10, maxH)
		m.list.SetSize(w, h)
		m.viewport.Width = w
		m.viewport.Height = h
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			m.resolution = ResolutionAbort
			return m, tea.Quit
		}

		if m.state == stateDiff {
			switch msg.String() {
			case "esc", "backspace":
				m.state = stateSelect
				return m, nil
			case "q":
				m.quitting = true
				m.resolution = ResolutionAbort
				return m, tea.Quit
			}
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

		// stateSelect
		switch msg.String() {
		case "enter":
			it := m.list.SelectedItem()
			choice, ok := it.(existingPathChoice)
			if !ok {
				m.err = errors.New("invalid selection")
				m.resolution = ResolutionAbort
				return m, tea.Quit
			}
			if choice.key == ExistingPathResolution("diff") {
				m.diffText = computeDiff(m.in.Source, m.in.Destination)
				m.viewport.SetContent(m.diffText)
				m.viewport.GotoTop()
				m.state = stateDiff
				return m, nil
			}
			m.resolution = choice.key
			m.quitting = true
			return m, tea.Quit
		case "q":
			m.quitting = true
			m.resolution = ResolutionAbort
			return m, tea.Quit
		}
	}

	if m.state == stateDiff {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *existingPathModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Render("Destination: ") + m.in.Destination + "\n" +
		lipgloss.NewStyle().Bold(true).Render("Source:      ") + m.in.Source + "\n\n"

	if m.state == stateDiff {
		help := "\n" + lipgloss.NewStyle().Faint(true).Render("esc/backspace: back • ↑/↓/pgup/pgdn: scroll • q: abort")
		return header + m.viewport.View() + help
	}

	footer := "\n" + lipgloss.NewStyle().Faint(true).Render("enter: select • q/ctrl+c: abort")
	return header + m.list.View() + footer
}

func computeDiff(source, destination string) string {
	if strings.TrimSpace(source) == "" || strings.TrimSpace(destination) == "" {
		return "Diff unavailable: missing source/destination path.\n"
	}

	out, err := diffPaths(destination, source)
	if err != nil {
		return fmt.Sprintf("Diff unavailable: %v\n", err)
	}
	if strings.TrimSpace(out) == "" {
		return "No differences.\n"
	}
	return colorizeUnifiedDiff(out)
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func colorizeUnifiedDiff(diff string) string {
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	delStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	hunkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	fileStyle := lipgloss.NewStyle().Bold(true)
	metaStyle := lipgloss.NewStyle().Faint(true)
	binStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)

	lines := strings.SplitAfter(diff, "\n")
	var out strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSuffix(line, "\n")
		hasNL := strings.HasSuffix(line, "\n")

		colored := trimmed
		switch {
		case strings.HasPrefix(trimmed, "+++ ") || strings.HasPrefix(trimmed, "--- "):
			colored = fileStyle.Render(trimmed)
		case strings.HasPrefix(trimmed, "@@"):
			colored = hunkStyle.Render(trimmed)
		case strings.HasPrefix(trimmed, "+") && !strings.HasPrefix(trimmed, "+++"):
			colored = addStyle.Render(trimmed)
		case strings.HasPrefix(trimmed, "-") && !strings.HasPrefix(trimmed, "---"):
			colored = delStyle.Render(trimmed)
		case strings.HasPrefix(trimmed, "diff "):
			colored = metaStyle.Render(trimmed)
		case strings.HasPrefix(trimmed, "Binary "):
			colored = binStyle.Render(trimmed)
		}

		out.WriteString(colored)
		if hasNL {
			out.WriteString("\n")
		}
	}

	return out.String()
}

func diffPaths(aPath, bPath string) (string, error) {
	aPathAbs, err := filepath.Abs(aPath)
	if err != nil {
		return "", err
	}
	bPathAbs, err := filepath.Abs(bPath)
	if err != nil {
		return "", err
	}

	aIsDir := isDir(aPathAbs)
	bIsDir := isDir(bPathAbs)

	switch {
	case aIsDir && bIsDir:
		return diffDirs(aPathAbs, bPathAbs)
	case !aIsDir && !bIsDir:
		return diffFiles(aPathAbs, bPathAbs, aPathAbs, bPathAbs)
	default:
		return "", fmt.Errorf("cannot diff file vs directory (%s vs %s)", aPathAbs, bPathAbs)
	}
}

func diffDirs(aDir, bDir string) (string, error) {
	aFiles, err := listRegularFiles(aDir)
	if err != nil {
		return "", err
	}
	bFiles, err := listRegularFiles(bDir)
	if err != nil {
		return "", err
	}

	paths := make([]string, 0, len(aFiles)+len(bFiles))
	seen := make(map[string]struct{}, len(aFiles)+len(bFiles))
	for rel := range aFiles {
		if _, ok := seen[rel]; !ok {
			seen[rel] = struct{}{}
			paths = append(paths, rel)
		}
	}
	for rel := range bFiles {
		if _, ok := seen[rel]; !ok {
			seen[rel] = struct{}{}
			paths = append(paths, rel)
		}
	}
	sort.Strings(paths)

	var out strings.Builder
	for _, rel := range paths {
		aFull, aOk := aFiles[rel]
		bFull, bOk := bFiles[rel]

		switch {
		case aOk && bOk:
			d, err := diffFiles(aFull, bFull, filepath.Join(aDir, rel), filepath.Join(bDir, rel))
			if err != nil {
				fmt.Fprintf(&out, "diff %s\nerror: %v\n\n", rel, err)
				continue
			}
			if strings.TrimSpace(d) != "" {
				out.WriteString(d)
				if !strings.HasSuffix(d, "\n") {
					out.WriteString("\n")
				}
				out.WriteString("\n")
			}
		case aOk && !bOk:
			// Removed from bDir relative to aDir
			d, err := diffFiles(aFull, "", filepath.Join(aDir, rel), filepath.Join(bDir, rel))
			if err != nil {
				fmt.Fprintf(&out, "diff %s\nerror: %v\n\n", rel, err)
				continue
			}
			out.WriteString(d)
			if !strings.HasSuffix(d, "\n") {
				out.WriteString("\n")
			}
			out.WriteString("\n")
		case !aOk && bOk:
			// Added in bDir relative to aDir
			d, err := diffFiles("", bFull, filepath.Join(aDir, rel), filepath.Join(bDir, rel))
			if err != nil {
				fmt.Fprintf(&out, "diff %s\nerror: %v\n\n", rel, err)
				continue
			}
			out.WriteString(d)
			if !strings.HasSuffix(d, "\n") {
				out.WriteString("\n")
			}
			out.WriteString("\n")
		}
	}

	return out.String(), nil
}

func listRegularFiles(root string) (map[string]string, error) {
	files := make(map[string]string)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = path
		return nil
	})
	return files, err
}

func diffFiles(aFile, bFile, aLabel, bLabel string) (string, error) {
	aLines, aBinary, err := readTextLines(aFile)
	if err != nil {
		return "", err
	}
	bLines, bBinary, err := readTextLines(bFile)
	if err != nil {
		return "", err
	}

	// Best-effort: if either side looks binary, just say so (like `diff` does).
	if aBinary || bBinary {
		if aFile == "" {
			return fmt.Sprintf("Binary file added: %s\n", bLabel), nil
		}
		if bFile == "" {
			return fmt.Sprintf("Binary file removed: %s\n", aLabel), nil
		}
		return fmt.Sprintf("Binary files differ: %s and %s\n", aLabel, bLabel), nil
	}

	ud := difflib.UnifiedDiff{
		A:        aLines,
		B:        bLines,
		FromFile: aLabel,
		ToFile:   bLabel,
		Context:  3,
	}
	text, err := difflib.GetUnifiedDiffString(ud)
	if err != nil {
		return "", err
	}
	return text, nil
}

func readTextLines(path string) ([]string, bool, error) {
	if strings.TrimSpace(path) == "" {
		return []string{}, false, nil
	}

	b, err := os.ReadFile(path)
	if err != nil {
		// If one side doesn't exist (added/removed), treat as empty.
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, false, nil
		}
		return nil, false, err
	}

	// Basic binary detection: if there are NUL bytes, treat as binary.
	if bytes.IndexByte(b, 0) >= 0 {
		return nil, true, nil
	}

	// Normalize to '\n' lines. difflib expects a slice of strings, each line including '\n' is fine.
	s := string(b)
	if s == "" {
		return []string{}, false, nil
	}
	lines := strings.SplitAfter(s, "\n")
	// If file doesn't end with newline, SplitAfter returns a last line without '\n' — that's OK.
	return lines, false, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clamp(v, minV, maxV int) int {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
