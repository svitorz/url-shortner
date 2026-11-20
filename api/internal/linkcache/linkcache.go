// Package linkcache guarda no Redis o resultado da busca do redirect
// (cache-aside): chave link:<slug> com o destino e se o link está ativo, ou
// um marcador de "não existe" com TTL curto (cache negativo).
//
// Leituras preenchem com SET NX e escritas (criar, desativar, excluir dono)
// sobrescrevem com o estado novo. Assim um leitor que leu o banco antes de
// uma desativação não consegue gravar o estado velho por cima do novo.
//
// Erros do Redis são devolvidos para quem chama decidir; o redirect trata
// qualquer erro como "miss" e vai ao banco.
package linkcache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "link:"

// opTimeout limita cada operação: com o Redis travado, o redirect deve cair
// no banco em vez de esperar o timeout padrão do cliente (3 s).
const opTimeout = 100 * time.Millisecond

// Valores gravados (curtos, sem JSON): "a:<url>" ativo, "i" inativo, "n" não existe.
const (
	activePrefix = "a:"
	inactiveVal  = "i"
	notFoundVal  = "n"
)

// Entry é o que o redirect precisa saber sobre um slug.
type Entry struct {
	Found  bool
	Active bool
	Target string
}

type Cache struct {
	rdb         *redis.Client
	ttl         time.Duration
	negativeTTL time.Duration
}

func New(rdb *redis.Client, ttl, negativeTTL time.Duration) *Cache {
	return &Cache{rdb: rdb, ttl: ttl, negativeTTL: negativeTTL}
}

func Key(slug string) string { return keyPrefix + slug }

var errBadValue = errors.New("valor inválido no cache")

// Get devolve (entry, true, nil) num acerto e (_, false, nil) num miss.
func (c *Cache) Get(ctx context.Context, slug string) (Entry, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	v, err := c.rdb.Get(ctx, Key(slug)).Result()
	if errors.Is(err, redis.Nil) {
		return Entry{}, false, nil
	}
	if err != nil {
		return Entry{}, false, err
	}
	switch {
	case v == notFoundVal:
		return Entry{}, true, nil
	case v == inactiveVal:
		return Entry{Found: true}, true, nil
	case strings.HasPrefix(v, activePrefix):
		return Entry{Found: true, Active: true, Target: v[len(activePrefix):]}, true, nil
	}
	return Entry{}, false, errBadValue
}

// Fill grava só se a chave não existe (preenchimento após um miss).
func (c *Cache) Fill(ctx context.Context, slug string, e Entry) error {
	ctx, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	v, ttl := c.encode(e)
	return c.rdb.SetNX(ctx, Key(slug), v, ttl).Err()
}

// Put sobrescreve com o estado novo (após uma escrita no banco).
func (c *Cache) Put(ctx context.Context, slug string, e Entry) error {
	ctx, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	v, ttl := c.encode(e)
	return c.rdb.Set(ctx, Key(slug), v, ttl).Err()
}

// PutNotFound marca vários slugs como inexistentes (ex.: dono excluído).
func (c *Cache) PutNotFound(ctx context.Context, slugs []string) error {
	if len(slugs) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	pipe := c.rdb.Pipeline()
	for _, s := range slugs {
		pipe.Set(ctx, Key(s), notFoundVal, c.negativeTTL)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (c *Cache) encode(e Entry) (string, time.Duration) {
	switch {
	case !e.Found:
		return notFoundVal, c.negativeTTL
	case !e.Active:
		return inactiveVal, c.ttl
	}
	return activePrefix + e.Target, c.ttl
}
