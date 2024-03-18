package process

import (
	"time"
)

type ProcessOpt = func(p *Process)

// WithEnv добавляет переменные окружения вида KEY=VALUE процессу.
func WithEnv(env ...string) ProcessOpt {
	return func(p *Process) {
		p.cmd.Env = append(p.cmd.Env, env...)
	}
}

// WithArgs добавляет процессу аргументы командной строки.
func WithArgs(args ...string) ProcessOpt {
	return func(p *Process) {
		p.cmd.Args = append(p.cmd.Args, args...)
	}
}

// WaitPortConnTimeout устанавливает таймаут на поключение к порту.
func WaitPortConnTimeout(d time.Duration) ProcessOpt {
	return func(p *Process) {
		p.waitPortConnTimeout = d
	}
}

// WaitPortInterval устанавливает таймаут на ожидание порта.
func WaitPortInterval(d time.Duration) ProcessOpt {
	return func(p *Process) {
		p.waitPortInterval = d
	}
}
