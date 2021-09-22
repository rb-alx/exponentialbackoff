// Simple exponential backoff
// Простейшая реализация экспоненциальной задержки

package exponentialbackoff

import (
	"context"
	"sync"
	"time"
)

type Config struct {
	Max    int `json:"max" yaml:"max"`       // Максимальное значение экспоненциальной задержки
	Factor int `json:"factor" yaml:"factor"` // Коэффициент увеличения задержки
}

type Delay struct {
	sync.RWMutex
	isInit        bool // will be false by default and changed in New func
	i             int  // will be zero by default
	max           int
	factor        int
	durationUnits time.Duration
}

// New ...
// Возвращает инициализированный объект
// экспоненциальной задержки
func New(c *Config) *Delay {

	if c.Max < 0 {
		c.Max = 0
	}

	if c.Factor < 1 {
		c.Factor = 1
	}

	return &Delay{
		isInit:        true,
		max:           c.Max,
		factor:        c.Factor,
		durationUnits: time.Second,
	}
}

// Incr ...
// Увеличение задержки
func (d *Delay) Incr() *Delay {

	if !d.isInit {
		return d
	}

	d.Lock()
	defer d.Unlock()

	if d.i == d.max {
		return d
	}

	d.i = d.i*d.factor + d.factor

	if d.i > d.max {
		d.i = d.max
	}

	return d
}

// Decr ...
// Уменьшение задержки
func (d *Delay) Decr() *Delay {

	if !d.isInit {
		return d
	}

	d.Lock()
	defer d.Unlock()

	if d.i == 0 {
		return d
	}

	d.i--

	if d.i < 0 {
		d.i = 0
	}

	return d
}

// Reset ...
func (d *Delay) Reset() *Delay {

	if !d.isInit {
		return d
	}

	d.Lock()
	defer d.Unlock()

	if d.i != 0 {
		d.i = 0
	}

	return d
}

// GetDelay ...
func (d *Delay) GetDelay() int {
	return d.i
}

// SetDelay ...
// Установить значение задежки
func (d *Delay) SetDelay(v int) *Delay {
	d.i = v
	return d
}

// SetDurationUnits
// Установить единицу времени, в которой будет измеряться задержка
//
func (d *Delay) SetDurationUnits(du time.Duration) *Delay {
	if d.isInit {
		d.durationUnits = du
	}

	return d
}

// IssetDelay ...
// Установлена ли задержка
func (d *Delay) IssetDelay() bool {

	if !d.isInit {
		return false
	}

	return d.i > 0
}

// Backoff ...
// Выполнить задержку, если возможно
//
// Принимает:
// 	context.Context - для отмены операции задержки
//
// Возвращает:
// 	bool - была ли задержка
// 	error - ошибка, если задержка была прервана
// 	time.Duration - фактическое время задержки
func (d *Delay) Backoff(ctx context.Context) (bool, error, time.Duration) {

	if !d.isInit {
		return false, nil, 0
	}

	ts := time.Now()
	isd := d.IssetDelay()

	if isd {

		select {
		case <-time.After(time.Duration(d.GetDelay()) * d.durationUnits):
		case <-ctx.Done():
		}
	}
	return isd, ctx.Err(), time.Since(ts)
}
