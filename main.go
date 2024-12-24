package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Yashsharma1911/sigNoz-assignment/metrics"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var (
	serviceName = "items-app"
)

func main() {
	// Initialize OpenTelemetry metrics
	cleanup := metrics.InitTracer()
	defer cleanup(context.Background())

	meterProvider := metrics.InitMeter()
	meter := meterProvider.Meter(serviceName)

	router := gin.Default()

	// This will ensure on each API call we are calling out metrics function
	// to generate metrics like histogram, gauge
	// In production application you might not do this for each metrics
	router.Use(otelgin.Middleware(serviceName))
	metrics.MetricsGenerator(meter)
	tracer := otel.Tracer(serviceName)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hurray! welcome to application"})
	})

	router.PUT("/update", func(c *gin.Context) {
		time.Sleep(1000 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"message": "Data updated!"})
	})

	router.POST("/create", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Data created!"})
	})

	router.DELETE("/delete", func(c *gin.Context) {
		// Creating a intentional nil pointer difference
		// This is created to record errors, however not in use currently
		ctx, span := tracer.Start(c.Request.Context(), "DeleteHandler")
		defer span.End()
		var someStruct *struct {
			Message string
		}

		message := someStruct.Message
		metrics.RecordSpanError(ctx, errors.New("intentional nil pointer dereference"))
		span.SetStatus(codes.Error, "Delete operation failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete data!",
			"message": message,
		})
	})

	// Start the Gin server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting Gin server: ", err)
	}
}
