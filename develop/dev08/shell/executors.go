package shell

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/process"
	"os"
	"os/user"
	"strings"
)

var (
	errEmptyEcho             = errors.New("shell: команде echo требуется передать какое-либо аргумент")
	errNeedProcess           = errors.New("shell: команде kill требуется имя процесса")
	errCanNotKillProcess     = errors.New("shell: этот процесс не может быть прерван")
)

// Executor ...
type Executor interface {
	Execute(s *Shell) (string, error)
}

// CDExecutor ...
type CDExecutor struct{}

// Execute ...
func (c *CDExecutor) Execute(s *Shell) (string, error) {
	// Если аргументом к команде "cd" не передан путь, по которому нужно перейти...
	if len(s.Args) == 1 {
		cUser, err := user.Current() //...то получаем текущего пользователя...
		if err != nil {
			return "", err
		}
		if err := os.Chdir(cUser.HomeDir); err != nil { // ...а затем меняем рабочую директорию на домашнюю для текущего пользователя.
			return "", err
		} // ...
	} else { // Если помимо команды "cd" пользователь указал путь до директории, которую нужно установить в качестве рабочей...
		if err := os.Chdir(s.Args[1]); err != nil { // ...пытаемся сделать это и возвращаем ненулевое значение ошибки если попытка оказалась неудачной...
			return "", err
		}

	}

	return "Directory changed successfully", nil // ... в притивном случае возвращаем сообщение об успехе
}

// EchoExecutor ...
type EchoExecutor struct{}

// Execute ...
func (e *EchoExecutor) Execute(s *Shell) (string, error) {
	// Если аргумент не передан...
	if len(s.Args) == 1 {
		return "", errEmptyEcho // ...возвращаем ошибку
	}
	return s.Args[1], nil
}


type PSExecutor struct{}

// Execute выводит информацию о работе запущенных процессов в системе
func (p *PSExecutor) Execute(s *Shell) (string, error) {
	processes, err := ps.Processes() // Создаёт моментальный снимок таблицы процессов (Возвращает слайс значений Process (интерфейсного типа))
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	builder.WriteString("PID\t|\tCOMMAND\n")
	builder.WriteString("---------------\n")
	for _, proc := range processes {
		builder.WriteString(
			fmt.Sprintf("%v\t|\t%v\n", proc.Pid(), proc.Executable()), // методы Pid и Executable возвращают id данного процесса и имя исполняемого файла, запускающего этот процесс соответственно
		)
	}
	builder.WriteString("---------------\n")
	return builder.String(), nil
}

// PWDExecutor ...
type PWDExecutor struct{}

// Execute выводит текущий рабочий каталог
func (p *PWDExecutor) Execute(s *Shell) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return "Current work directory:" + path, err
}

// KillProcessExecutor ...
type KillProcessExecutor struct{}

// Execute ...
func (k *KillProcessExecutor) Execute(s *Shell) (string, error) {
	// Проверим, ввёл ли пользователь что-либо кроме команды kill
	if len(s.Args) < 2 {
		return "", errNeedProcess
	}
	processes, err := process.Processes() // Создаёт моментальный снимок таблицы процессов (Возвращает слайс указателей на значения Process (уже конкретные))
	if err != nil {
		return "", err
	}
	for _, process := range processes { // Перебираем все запущенные процессы и сравниваем имя каждого с тем, что ввёл пользователь в кач-ве имени процесса, который нужно убить...
		name, err := process.Name()
		if err != nil {
			return "", err
		}
		if name == s.Args[1] {
			if err := process.Kill(); err != nil { // ...если есть совпадение, пытаемся убить указанный процесс
				return "", errCanNotKillProcess
			}
		}
	}
	return fmt.Sprintf("Process %v successfully killed", s.Args[1]), nil
}

