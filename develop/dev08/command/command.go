package command

import (
	"bufio"
	"dev08/shell"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

var (

	ShellTerminalPrefix = `vlad@shell:`

	ShellExitCommand   = `\exit`
	errConsoleInput    = errors.New("shell: чтение не удалось")
	errPromptBuildFail = errors.New("shell: не удалось определить текущую директорию")
	errBadUser         = errors.New("shell: не удалось получить информацию о текущем пользователе")
	successExitMessage = "shell: exit successful"
)

// ShellTerminal ...
type ShellTerminal struct {
	shell  *shell.Shell
	reader io.Reader
	writer io.Writer
}

// New ...
func New(shell *shell.Shell, reader io.Reader, writer io.Writer) *ShellTerminal {
	return &ShellTerminal{
		shell:  shell,
		reader: reader,
		writer: writer,
	}
}

// Start
func (s *ShellTerminal) Start() error {
	fmt.Fprintln(s.writer, `ПРОСТЕЙШАЯ ОБОЛОЧКА ИНТЕРПРЕТАТОРА КОМАНДНОЙ СТОКИ. ДЛЯ ВЫХОДА ИСПОЛЬЗУЙТЕ \exit.`)
	scanner := bufio.NewScanner(s.reader)
	for {
		prompt, err := s.buildPrompt()
		if err != nil {
			return errPromptBuildFail
		}
		fmt.Fprint(s.writer, prompt)
		scanner.Scan()
		text := scanner.Text()
		if text == ShellExitCommand {
			break
		}
		args := strings.Fields(text) // Разделяет строку на подстроки по следующим разделителям: '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP)
		s.shell.SetArgs(args) // В зависимости от того, какую команду ввёл пользователь, выполняем соответствующий блок case
		switch args[0] {
		case "cd":
			cdExecutor := &shell.CDExecutor{}
			s.shell.SetExecutor(cdExecutor)
		case "echo":
			echoExecutor := &shell.EchoExecutor{}
			s.shell.SetExecutor(echoExecutor)
		case "ps":
			psExecutor := &shell.PSExecutor{}
			s.shell.SetExecutor(psExecutor)
		case "pwd":
			pwdExecutor := &shell.PWDExecutor{}
			s.shell.SetExecutor(pwdExecutor)
		case "kill":
			killExecutor := &shell.KillProcessExecutor{}
			s.shell.SetExecutor(killExecutor)

		default:
			fmt.Fprintln(s.writer, "shell: неизвестная команда")
		}
		res, err := s.shell.Start()
		if err != nil {
			fmt.Fprintln(s.writer, err.Error())
			continue
		}
		fmt.Fprintln(s.writer, res)
	}
	if scanner.Err() != nil {
		return errConsoleInput
	}
	if _, err := fmt.Fprintln(s.writer, successExitMessage); err != nil {
		return err
	}
	return nil
}
// buildPrompt конструирует приглашение командной строки, поскольку в зависимости от текущих рабочей директории или пользователя оно будет отличаться
func (s *ShellTerminal) buildPrompt() (string, error) {
	path, err := os.Getwd() // Getwd возвращает абсолютный путь до текущего (рабочего) каталога
	if err != nil {
		return "", err
	}
	var postfix string
	userName, err := user.Current()
	if err != nil {
		return "", errBadUser
	}
	if path == "/home/"+userName.Name { // Поскольку /home/$USERNAME эквивалентно "~"
		postfix = "~$ "
	} else {
		postfix = path + " "
	}
	return ShellTerminalPrefix + postfix, nil
}
