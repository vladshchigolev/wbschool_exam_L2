package shell

// Shell - объект-контекст.
// Вместо того, чтобы изначальный класс сам выполнял тот или иной алгоритм,
// он будет играть роль контекста, ссылаясь на одну из стратегий и делегируя ей выполнение работы.
type Shell struct {
	Args     []string
	Executor Executor // Будем устанавливать Executor конкретного типа, в зависимости от введённой пользователем команды
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
