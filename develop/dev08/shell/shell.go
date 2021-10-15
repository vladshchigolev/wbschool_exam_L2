package shell

// Shell - объект-контекст.
// Вместо того, чтобы изначальный класс сам выполнял тот или иной алгоритм,
// он будет играть роль контекста, ссылаясь на одну из стратегий и делегируя ей выполнение работы.
type Shell struct {
	Args     []string
	Executor Executor // Используем паттерн "Стратегия"
}

// New ...
func New() *Shell {
	return &Shell{}
}

// SetArgs ...
func (s *Shell) SetArgs(args []string) {
	s.Args = args
}

// SetExecutor ...
func (s *Shell) SetExecutor(e Executor) {
	s.Executor = e
}

// Start ...
func (s *Shell) Start() (string, error) {
	return s.Executor.Execute(s)
}
