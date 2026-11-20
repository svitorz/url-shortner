// Package clicks conta acessos aos links sem pôr uma escrita no banco no
// caminho do redirect: cada clique é um INCR em clicks:<slug> no Redis, e um
// flush periódico soma esses contadores na coluna links.clicks.
//
// Trade-offs (ver PROGRESS.md):
//   - o valor no banco fica até um intervalo de flush atrasado;
//   - Redis fora do ar no momento do clique = clique não contado (o redirect
//     segue funcionando);
//   - contadores ainda não descarregados se perdem se o Redis reiniciar sem
//     persistência.
package clicks

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"svitorz/url-shortner/internal/models"
)

const keyPrefix = "clicks:"

// opTimeout: o INCR está no caminho do redirect; com o Redis travado, desiste
// rápido em vez de segurar a resposta (mesmo limite do cache de links).
const opTimeout = 100 * time.Millisecond

func Key(slug string) string { return keyPrefix + slug }

type Counter struct {
	rdb *redis.Client
	db  *gorm.DB
}

func New(rdb *redis.Client, db *gorm.DB) *Counter {
	return &Counter{rdb: rdb, db: db}
}

// Incr conta um clique. Erro só é devolvido para quem chama logar.
func (c *Counter) Incr(ctx context.Context, slug string) error {
	ctx, cancel := context.WithTimeout(ctx, opTimeout)
	defer cancel()
	return c.rdb.Incr(ctx, Key(slug)).Err()
}

// FlushOnce descarrega todos os contadores pendentes no banco e devolve
// quantos cliques foram gravados. GETDEL lê e zera de forma atômica, então
// dois processos fazendo flush ao mesmo tempo não contam o mesmo clique duas
// vezes. Se o UPDATE falhar, o valor volta para o Redis (INCRBY) e entra no
// próximo flush.
func (c *Counter) FlushOnce(ctx context.Context) (int64, error) {
	var total int64
	iter := c.rdb.Scan(ctx, 0, keyPrefix+"*", 500).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		n, err := c.rdb.GetDel(ctx, key).Int64()
		if errors.Is(err, redis.Nil) {
			continue // outro processo já descarregou
		}
		if err != nil {
			return total, err
		}
		slug := strings.TrimPrefix(key, keyPrefix)
		// UpdateColumn: não mexe em updated_at (clique não é "alteração" do link).
		res := c.db.WithContext(ctx).Model(&models.Link{}).Where("slug = ?", slug).
			UpdateColumn("clicks", gorm.Expr("clicks + ?", n))
		if res.Error != nil {
			if rerr := c.rdb.IncrBy(context.WithoutCancel(ctx), key, n).Err(); rerr != nil {
				log.Printf("cliques: %d cliques de %s perdidos (banco: %v; redis: %v)", n, slug, res.Error, rerr)
			}
			return total, res.Error
		}
		if res.RowsAffected == 0 {
			continue // link apagado depois do clique: não há onde somar
		}
		total += n
	}
	return total, iter.Err()
}

// Run faz flush a cada interval até ctx terminar, e um último flush na saída
// (desligamento gracioso), para não deixar cliques parados no Redis.
func (c *Counter) Run(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			final, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
			if _, err := c.FlushOnce(final); err != nil {
				log.Printf("cliques: flush final falhou: %v", err)
			}
			cancel()
			return
		case <-t.C:
			if _, err := c.FlushOnce(ctx); err != nil {
				log.Printf("cliques: flush falhou: %v", err)
			}
		}
	}
}
