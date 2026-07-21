package git

import (
	"strings"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
)

func Status(pth string) (uncommitted StatusDiff, err error) {
	req := cmd.Request{
		Tool:    bin,
		Args:    []string{"status"},
		WorkDir: pth,
	}

	executeOut, err := cmd.Execute(req)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting git status")
	}

	out := make([]Changes, 0)

	untracked, ok := parseUntrackedFiles(executeOut)
	if ok {
		out = append(out, untracked)
	}

	out = append(out, parseCommitChanges(executeOut)...)

	return out, nil
}

func parseUntrackedFiles(executeOut string) (Changes, bool) {
	const messageForUntrackedFiles = "Untracked files"

	startIdx := strings.Index(executeOut, messageForUntrackedFiles)
	if startIdx == -1 {
		return Changes{}, false
	}

	startIdx += len(messageForUntrackedFiles)

	changeList := strings.Split(executeOut[startIdx:], "\n")

	gitChanges := Changes{
		Type:       ChangesTypeNotStaged,
		Changelist: make([]string, 0, len(changeList)),
	}
	for _, item := range changeList {
		if len(item) == 0 {
			continue
		}

		if strings.HasPrefix(item, "\t") {
			gitChanges.Changelist = append(gitChanges.Changelist, item[1:])
		}
	}

	return gitChanges, true
}

func parseCommitChanges(executeOut string) []Changes {
	var keyWords = []string{"deleted", "modified", "new file"}

	out := make([]Changes, 0)

	for _, message := range []string{"Changes to be committed", "Changes not staged for commit"} {
		startIdx := strings.Index(executeOut, message)
		if startIdx == -1 {
			continue
		}

		startIdx += len(message)

		changeList := strings.Split(executeOut[startIdx:], "\n")

		gitChanges := Changes{
			Type:       ChangesTypeNotCommitted,
			Changelist: make([]string, 0, len(changeList)),
		}

		for _, item := range changeList {
			if len(item) == 0 {
				continue
			}

			item = strings.ReplaceAll(item, "\t", "")
			for _, keyWord := range keyWords {
				if strings.HasPrefix(item, keyWord) {
					gitChanges.Changelist = append(gitChanges.Changelist, item)

					break
				}
			}
		}

		out = append(out, gitChanges)
	}

	return out
}

type Changes struct {
	Type       gitChangesType
	Changelist []string
}

type StatusDiff []Changes

func (s StatusDiff) GetFilesListed() string {
	const splittedParamsCount = 2

	sb := strings.Builder{}

	for _, item := range s {
		if len(item.Changelist) == 0 {
			continue
		}

		sb.WriteString(item.Type.Msg())
		sb.WriteString("\n")

		changeTypeToFile := map[string][]string{}

		for _, line := range item.Changelist {
			splited := strings.Split(line, ":")
			if len(splited) != splittedParamsCount {
				continue
			}

			changeTypeToFile[splited[0]] = append(changeTypeToFile[splited[0]], splited[1])
		}

		for k, v := range changeTypeToFile {
			sb.WriteString(k + ": \n\t")
			sb.WriteString(strings.Join(v, "; "))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
func (s StatusDiff) String() string {
	sb := strings.Builder{}

	for _, item := range s {
		if len(item.Changelist) == 0 {
			continue
		}

		sb.WriteString(item.Type.Msg())
		sb.WriteString("\n\t")
		sb.WriteString(strings.Join(item.Changelist, "\n\t"))
	}

	return sb.String()
}

type gitChangesType int

func (g gitChangesType) Msg() string {
	switch g {
	case ChangesTypeNotStaged:
		return "Not Staged changes:"
	case ChangesTypeNotCommitted:
		return "Not committed changes:"
	default:
		return "Unknown git changes type!"
	}
}
