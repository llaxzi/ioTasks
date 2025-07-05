package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"workmate/internal/models"
	"workmate/internal/storage"
)

var (
	ErrWrongBody = errors.New("wrong JSON body")
	ErrServer    = errors.New("server error")
)

type Handler struct {
	st storage.StorageI
}

func New(st storage.StorageI) *Handler {
	return &Handler{st: st}
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrMap(ErrWrongBody))
		return
	}

	var id models.TaskID
	if err = json.Unmarshal(body, &id); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrMap(ErrWrongBody))
		return
	}

	info, err := h.st.GetInfo(id.ID)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			WriteJSON(w, http.StatusBadRequest, ErrMap(err))
			return
		}
		WriteJSON(w, http.StatusInternalServerError, ErrMap(ErrServer))
		return
	}
	WriteJSON(w, http.StatusOK, info)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.st.GetTasks()
	if tasks == nil {
		tasks = make([]models.TaskID, 0)
	}
	WriteJSON(w, http.StatusOK, tasks)
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	id := h.st.Add()
	WriteJSON(w, http.StatusOK, map[string]string{"id": id})

}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrMap(ErrWrongBody))
		return
	}

	var id models.TaskID
	if err = json.Unmarshal(body, &id); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrMap(ErrWrongBody))
		return
	}

	if err = h.st.Delete(id.ID); err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) || errors.Is(err, storage.ErrTaskRunning) {
			WriteJSON(w, http.StatusBadRequest, ErrMap(err))
			return
		}

		WriteJSON(w, http.StatusInternalServerError, ErrMap(ErrServer))
		return
	}
	WriteJSON(w, http.StatusOK, "deleted")
}

func WriteJSON(w http.ResponseWriter, status int, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		log.Printf("Error writing json: %v", err)
	}
}

func ErrMap(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}
