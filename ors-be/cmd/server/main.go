package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpsrv "ors-be/internal/api/http"
	"ors-be/internal/api/http/handler"
	"ors-be/internal/auth"
	"ors-be/internal/config"
	"ors-be/internal/repository/postgres"
	"ors-be/internal/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer pool.Close()
	log.Println("数据库连接成功")

	hasher := auth.NewHasher()
	tokenGen := auth.NewTokenGenerator(cfg.JWTSecret, cfg.JWTExpirationHours)

	userRepo := postgres.NewUserRepo(pool)
	providerRepo := postgres.NewServiceProviderRepo(pool)
	serviceRepo := postgres.NewServiceRepo(pool)
	tagRepo := postgres.NewTagRepo(pool)
	serviceTagRepo := postgres.NewServiceTagRepo(pool)
	categoryRepo := postgres.NewCategoryRepo(pool)
	interestRepo := postgres.NewUserInterestRepo(pool)
	reservationRepo := postgres.NewReservationRepo(pool)
	reviewRepo := postgres.NewReviewRepo(pool)
	notificationRepo := postgres.NewNotificationRepo(pool)
	authSvc := service.NewAuthService(userRepo, hasher, tokenGen)
	userSvc := service.NewUserService(userRepo, hasher)
	providerSvc := service.NewServiceProviderService(providerRepo)
	serviceSvc := service.NewServiceService(serviceRepo, providerRepo, tagRepo, serviceTagRepo)
	tagSvc := service.NewTagService(tagRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	interestSvc := service.NewUserInterestService(tagRepo, interestRepo)
	notificationSvc := service.NewNotificationService(notificationRepo)
	reservationSvc := service.NewReservationService(reservationRepo, serviceRepo, providerRepo, notificationSvc)
	reviewSvc := service.NewReviewService(reviewRepo, reservationRepo)
	authH := handler.NewAuthHandler(authSvc)
	userH := handler.NewUserHandler(userSvc)
	providerH := handler.NewServiceProviderHandler(providerSvc)
	serviceH := handler.NewServiceHandler(serviceSvc)
	tagH := handler.NewTagHandler(tagSvc)
	categoryH := handler.NewCategoryHandler(categorySvc)
	interestH := handler.NewUserInterestHandler(interestSvc)
	reservationH := handler.NewReservationHandler(reservationSvc)
	reviewH := handler.NewReviewHandler(reviewSvc)
	notificationH := handler.NewNotificationHandler(notificationSvc)

	srv := httpsrv.NewServer(authH, userH, providerH, serviceH, tagH, categoryH, interestH, reservationH, reviewH, notificationH, tokenGen, cfg.AllowedOrigins)

	httpServer := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      srv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		addr := fmt.Sprintf(":%s", cfg.HTTPPort)
		log.Printf("服务启动，监听 %s", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	schedulerCtx, stopScheduler := context.WithCancel(context.Background())
	defer stopScheduler()
	startReservationCompletionScheduler(schedulerCtx, reservationSvc, time.Minute)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务...")
	stopScheduler()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("服务关闭异常: %v", err)
	}
	log.Println("服务已关闭")
}

func startReservationCompletionScheduler(ctx context.Context, reservationSvc service.ReservationService, interval time.Duration) {
	run := func() {
		runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		count, err := reservationSvc.CompleteDue(runCtx, time.Now())
		if err != nil {
			log.Printf("自动完成到期预约失败: %v", err)
			return
		}
		if count > 0 {
			log.Printf("自动完成到期预约数量: %d", count)
		}
	}

	go func() {
		run()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				run()
			}
		}
	}()
}
