package clicks

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"svitorz/url-shortner/internal/models"
	"svitorz/url-shortner/internal/testdb"
)

var ctx = context.Background()

type env struct {
	c  *Counter
	db *gorm.DB
	mr *miniredis.Miniredis
}

func setup(t *testing.T) env {
	t.Helper()
	db := testdb.Open(t)
	if err := models.LoadMigrations(db); err != nil {
		t.Fatal(err)
	}
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr(), ContextTimeoutEnabled: true})
	t.Cleanup(func() { rdb.Close() })
	return env{c: New(rdb, db), db: db, mr: mr}
}

// link cria um link e devolve o slug.
func (e env) link(t *testing.T) models.Link {
	t.Helper()
	email := fmt.Sprintf("c-%d@test.example", time.Now().UnixNano())
	u := models.User{Name: "C", Email: &email, Password: "x"}
	if err := e.db.Create(&u).Error; err != nil {
		t.Fatal(err)
	}
	l := models.Link{Slug: fmt.Sprintf("c%d", time.Now().UnixNano()), TargetURL: "https://e.com", IsActive: true, OwnerID: u.ID}
	if err := e.db.Create(&l).Error; err != nil {
		t.Fatal(err)
	}
	return l
}

func (e env) dbClicks(t *testing.T, slug string) models.Link {
	t.Helper()
	var l models.Link
	if err := e.db.Where("slug = ?", slug).First(&l).Error; err != nil {
		t.Fatal(err)
	}
	return l
}

func TestIncrThenFlushMovesCountToDB(t *testing.T) {
	e := setup(t)
	l := e.link(t)
	before := e.dbClicks(t, l.Slug) // updated_at como o banco guarda (microssegundos)
	for i := 0; i < 3; i++ {
		if err := e.c.Incr(ctx, l.Slug); err != nil {
			t.Fatal(err)
		}
	}
	n, err := e.c.FlushOnce(ctx)
	if err != nil || n != 3 {
		t.Fatalf("flush: n=%d err=%v", n, err)
	}
	got := e.dbClicks(t, l.Slug)
	if got.Clicks != 3 {
		t.Fatalf("clicks no banco = %d, esperado 3", got.Clicks)
	}
	if !got.UpdatedAt.Equal(before.UpdatedAt) {
		t.Fatalf("flush alterou updated_at (%v -> %v): clique não é alteração do link", before.UpdatedAt, got.UpdatedAt)
	}
	if e.mr.Exists(Key(l.Slug)) {
		t.Fatal("contador continuou no Redis depois do flush")
	}
	// Segundo flush sem cliques novos não soma nada.
	if n, _ := e.c.FlushOnce(ctx); n != 0 || e.dbClicks(t, l.Slug).Clicks != 3 {
		t.Fatalf("flush vazio somou %d", n)
	}
}

func TestFlushRestoresCountWhenDBFails(t *testing.T) {
	e := setup(t)
	l := e.link(t)
	_ = e.c.Incr(ctx, l.Slug)
	_ = e.c.Incr(ctx, l.Slug)

	sqlDB, _ := e.db.DB()
	sqlDB.Close() // banco indisponível
	if _, err := e.c.FlushOnce(ctx); err == nil {
		t.Fatal("flush deveria falhar com o banco fechado")
	}
	if v, err := e.mr.Get(Key(l.Slug)); err != nil || v != "2" {
		t.Fatalf("contador não voltou ao Redis: %q err=%v", v, err)
	}
}

func TestFlushDropsClicksOfDeletedLink(t *testing.T) {
	e := setup(t)
	_ = e.c.Incr(ctx, "nao-existe")
	n, err := e.c.FlushOnce(ctx)
	if err != nil || n != 0 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	if e.mr.Exists(Key("nao-existe")) {
		t.Fatal("contador de link inexistente ficou no Redis para sempre")
	}
}

// Dois flushes ao mesmo tempo (ex.: duas instâncias da API) não podem contar
// o mesmo clique duas vezes.
func TestConcurrentFlushesDoNotDoubleCount(t *testing.T) {
	e := setup(t)
	links := []models.Link{e.link(t), e.link(t), e.link(t)}
	for i := 0; i < 50; i++ {
		for _, l := range links {
			_ = e.c.Incr(ctx, l.Slug)
		}
	}
	var wg sync.WaitGroup
	totals := make([]int64, 4)
	for i := range totals {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			n, err := e.c.FlushOnce(ctx)
			if err != nil {
				t.Error(err)
			}
			totals[i] = n
		}(i)
	}
	wg.Wait()
	var sum int64
	for _, n := range totals {
		sum += n
	}
	if sum != 150 {
		t.Fatalf("flushes somaram %d, esperado 150", sum)
	}
	for _, l := range links {
		if c := e.dbClicks(t, l.Slug).Clicks; c != 50 {
			t.Fatalf("%s: %d cliques no banco, esperado 50", l.Slug, c)
		}
	}
}

func TestRunFlushesOnShutdown(t *testing.T) {
	e := setup(t)
	l := e.link(t)
	_ = e.c.Incr(ctx, l.Slug)

	runCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() { e.c.Run(runCtx, time.Hour); close(done) }() // ticker nunca dispara
	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Run não terminou após o cancelamento")
	}
	if c := e.dbClicks(t, l.Slug).Clicks; c != 1 {
		t.Fatalf("flush final não gravou: clicks=%d", c)
	}
}

func TestIncrFailsFastWithRedisDown(t *testing.T) {
	e := setup(t)
	e.mr.Close()
	start := time.Now()
	if err := e.c.Incr(ctx, "x"); err == nil {
		t.Fatal("esperava erro com Redis fora")
	}
	if d := time.Since(start); d > time.Second {
		t.Fatalf("Incr demorou %v com Redis fora", d)
	}
}
