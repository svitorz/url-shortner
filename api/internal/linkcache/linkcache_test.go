package linkcache

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

const (
	ttl    = 10 * time.Minute
	negTTL = 30 * time.Second
)

func setup(t *testing.T) (*Cache, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { rdb.Close() })
	return New(rdb, ttl, negTTL), mr
}

var ctx = context.Background()

func TestMissThenFillThenHit(t *testing.T) {
	c, mr := setup(t)
	if _, hit, err := c.Get(ctx, "abc"); hit || err != nil {
		t.Fatalf("esperava miss: hit=%v err=%v", hit, err)
	}
	want := Entry{Found: true, Active: true, Target: "https://example.com/x?y=1"}
	if err := c.Fill(ctx, "abc", want); err != nil {
		t.Fatal(err)
	}
	got, hit, err := c.Get(ctx, "abc")
	if !hit || err != nil || got != want {
		t.Fatalf("esperava hit %+v, veio %+v hit=%v err=%v", want, got, hit, err)
	}
	if d := mr.TTL(Key("abc")); d != ttl {
		t.Fatalf("TTL = %v, esperado %v", d, ttl)
	}
}

func TestInactiveAndNotFound(t *testing.T) {
	c, mr := setup(t)
	_ = c.Fill(ctx, "off", Entry{Found: true})
	_ = c.Fill(ctx, "none", Entry{})
	if e, hit, _ := c.Get(ctx, "off"); !hit || !e.Found || e.Active {
		t.Fatalf("inativo: %+v hit=%v", e, hit)
	}
	if e, hit, _ := c.Get(ctx, "none"); !hit || e.Found {
		t.Fatalf("inexistente: %+v hit=%v", e, hit)
	}
	if d := mr.TTL(Key("none")); d != negTTL {
		t.Fatalf("TTL do cache negativo = %v, esperado %v", d, negTTL)
	}
}

func TestExpiresAfterTTL(t *testing.T) {
	c, mr := setup(t)
	_ = c.Fill(ctx, "abc", Entry{Found: true, Active: true, Target: "https://e.com"})
	_ = c.Fill(ctx, "none", Entry{})
	mr.FastForward(negTTL)
	if _, hit, _ := c.Get(ctx, "none"); hit {
		t.Fatal("cache negativo não expirou no TTL curto")
	}
	if _, hit, _ := c.Get(ctx, "abc"); !hit {
		t.Fatal("entrada positiva expirou cedo demais")
	}
	mr.FastForward(ttl)
	if _, hit, _ := c.Get(ctx, "abc"); hit {
		t.Fatal("entrada positiva não expirou no TTL")
	}
}

// Fill não sobrescreve: um leitor atrasado não desfaz uma desativação.
func TestFillDoesNotOverwriteWrite(t *testing.T) {
	c, _ := setup(t)
	if err := c.Put(ctx, "abc", Entry{Found: true}); err != nil { // desativou
		t.Fatal(err)
	}
	_ = c.Fill(ctx, "abc", Entry{Found: true, Active: true, Target: "https://velho"}) // leitor atrasado
	if e, _, _ := c.Get(ctx, "abc"); e.Active {
		t.Fatal("Fill sobrescreveu o estado gravado por Put")
	}
}

func TestPutOverwrites(t *testing.T) {
	c, _ := setup(t)
	_ = c.Fill(ctx, "abc", Entry{}) // cache negativo
	_ = c.Put(ctx, "abc", Entry{Found: true, Active: true, Target: "https://novo"})
	if e, _, _ := c.Get(ctx, "abc"); e.Target != "https://novo" {
		t.Fatalf("Put não sobrescreveu: %+v", e)
	}
	_ = c.PutNotFound(ctx, []string{"abc"})
	if e, hit, _ := c.Get(ctx, "abc"); !hit || e.Found {
		t.Fatalf("PutNotFound não marcou como inexistente: %+v", e)
	}
}

func TestRedisDownReturnsErrorFast(t *testing.T) {
	c, mr := setup(t)
	mr.Close()
	start := time.Now()
	if _, hit, err := c.Get(ctx, "abc"); err == nil || hit {
		t.Fatalf("esperava erro com Redis fora: hit=%v err=%v", hit, err)
	}
	if d := time.Since(start); d > time.Second {
		t.Fatalf("Get demorou %v com Redis fora", d)
	}
}

// Redis travado (aceita a conexão e não responde): o Get tem que desistir no
// opTimeout, não no timeout padrão do cliente (3 s).
func TestRedisHangingTimesOut(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			defer conn.Close() // segura a conexão aberta sem responder
		}
	}()
	// Mesma opção do cliente de produção (repository.redisOptions).
	rdb := redis.NewClient(&redis.Options{Addr: ln.Addr().String(), ContextTimeoutEnabled: true})
	defer rdb.Close()
	c := New(rdb, ttl, negTTL)

	start := time.Now()
	if _, _, err := c.Get(ctx, "abc"); err == nil {
		t.Fatal("esperava erro de timeout")
	}
	if d := time.Since(start); d > 5*opTimeout {
		t.Fatalf("Get levou %v com Redis travado (limite %v)", d, opTimeout)
	}
}
