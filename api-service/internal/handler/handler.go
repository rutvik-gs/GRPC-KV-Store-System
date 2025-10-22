package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GRPC-KV-Store-System/api-service/internal/client"
)

type Handler struct {
	grpcClient *client.KVStoreClient
}

func NewHandler(grpcClient *client.KVStoreClient) *Handler {
	return &Handler{
		grpcClient: grpcClient,
	}
}

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func (h *Handler) SetHandler(w http.ResponseWriter, r *http.Request) {
	var req SetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Key == "" {
		h.respondError(w, http.StatusBadRequest, "Key cannot be empty")
		return
	}

	log.Printf("REST API: Setting key=%s", req.Key)

	err := h.grpcClient.Set(req.Key, req.Value)
	if err != nil {
		h.handleGRPCError(w, err)
		return
	}

	h.respondSuccess(w, http.StatusCreated, "Key-value pair stored successfully")
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		h.respondError(w, http.StatusBadRequest, "Key cannot be empty")
		return
	}

	log.Printf("REST API: Getting key=%s", key)

	value, err := h.grpcClient.Get(key)
	if err != nil {
		h.handleGRPCError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, GetResponse{
		Key:   key,
		Value: value,
	})
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		h.respondError(w, http.StatusBadRequest, "Key cannot be empty")
		return
	}

	log.Printf("REST API: Deleting key=%s", key)

	err := h.grpcClient.Delete(key)
	if err != nil {
		h.handleGRPCError(w, err)
		return
	}

	h.respondSuccess(w, http.StatusOK, "Key deleted successfully")
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	h.respondJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

func (h *Handler) respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) respondError(w http.ResponseWriter, statusCode int, message string) {
	h.respondJSON(w, statusCode, ErrorResponse{Error: message})
}

func (h *Handler) respondSuccess(w http.ResponseWriter, statusCode int, message string) {
	h.respondJSON(w, statusCode, SuccessResponse{Message: message})
}

func (h *Handler) handleGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		h.respondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	var httpStatus int
	switch st.Code() {
	case codes.NotFound:
		httpStatus = http.StatusNotFound
	case codes.InvalidArgument:
		httpStatus = http.StatusBadRequest
	case codes.AlreadyExists:
		httpStatus = http.StatusConflict
	case codes.PermissionDenied:
		httpStatus = http.StatusForbidden
	case codes.Unauthenticated:
		httpStatus = http.StatusUnauthorized
	default:
		httpStatus = http.StatusInternalServerError
	}

	h.respondError(w, httpStatus, st.Message())
}
