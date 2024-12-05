package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// *http.Server型をラップする
type server struct {
	srv *http.Server
}

// 独自のserver.Server型を返すコンストラクタ
func NewServer(port string, mux *http.ServeMux) *server {
	return &server{
		&http.Server{
			Addr:    port,
			Handler: mux,
		},
	}
}

// サーバーを起動する
func (s *server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	// シグナルの監視をやめてリソースを開放する
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// サーバー起動中にエラーが起きると、ctxに伝達される（Shutdownによるエラーは正常なので無視）
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("http server on %s failed : %+v", s.srv.Addr, err)
			return err
		}
		return nil
	})
	// サーバーからのエラーとシグナルを待機する
	<-ctx.Done()
	// シャットダウンをするまでのタイムアウト時間を設定
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// シャットダウン
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown http server on %s : %+v", s.srv.Addr, err)
	}
	// ゴルーチンのクロージャ内のエラーを返す
	// ゴルーチンを待機する
	return eg.Wait()
}