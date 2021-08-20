package handlers

import (
	"net/http"

	"github.com/arschles/megaboom/pkg/k8s"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

func StartJob(lggr logr.Logger, jobCreator k8s.JobCreator) gin.HandlerFunc {
	type reqBody struct {
		NumRunners             uint   `json:"num_runners"`
		NumConcurrentPerRunner uint   `json:"num_concurrent_per_runner"`
		NumReqsPerRunner       uint   `json:"num_reqs_per_runner"`
		Endpoint               string `json:"endpoint"`
	}
	type resBody struct {
		JobID string `json:"job_id"`
	}

	return func(ctx *gin.Context) {
		req := new(reqBody)
		if err := ctx.BindJSON(req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uid := uuid.New()
		jobs := k8s.NewJobs(
			uid,
			req.NumRunners,
			req.Endpoint,
			req.NumReqsPerRunner,
			req.NumConcurrentPerRunner,
		)
		if err := k8s.CreateJobs(ctx, jobCreator, jobs); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, &resBody{JobID: uid.String()})
	}
}

func DeleteJob(lggr logr.Logger, jobDeleter k8s.JobDeleter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jobID := ctx.Param("id")
		jobUUID, err := uuid.Parse(jobID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := k8s.DeleteJobs(ctx, jobDeleter, jobUUID); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, nil)
	}
}
