package handler

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "crm-mock/model"
)

type ContactHandler struct { db *sql.DB }

func NewContactHandler(db *sql.DB) *ContactHandler {
    return &ContactHandler{db: db}
}

func (h *ContactHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var c model.Contact
    err := h.db.QueryRow(
        `SELECT id, account_id, first_name, last_name, email, phone, created_at
         FROM contacts WHERE id = $1`, id,
    ).Scan(&c.ID, &c.AccountID, &c.FirstName, &c.LastName,
           &c.Email, &c.Phone, &c.CreatedAt)

    if err == sql.ErrNoRows {
        writeJSON(w, http.StatusNotFound,
            model.ErrorResponse{Error: "not_found", Message: "Contact not found", Code: 404})
        return
    }
    writeJSON(w, http.StatusOK, c)
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ID        string `json:"id"`
        AccountID string `json:"account_id"`
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
        Email     string `json:"email"`
        Phone     string `json:"phone"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest,
            model.ErrorResponse{Error: "invalid_body", Message: err.Error(), Code: 400})
        return
    }

    _, err := h.db.Exec(
        `INSERT INTO contacts (id, account_id, first_name, last_name, email, phone)
         VALUES ($1, $2, $3, $4, $5, $6)`,
        req.ID, req.AccountID, req.FirstName, req.LastName, req.Email, req.Phone,
    )
    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            model.ErrorResponse{Error: "db_error", Message: err.Error(), Code: 500})
        return
    }

    writeJSON(w, http.StatusCreated, req)
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res, err := h.db.Exec(`DELETE FROM contacts WHERE id = $1`, id)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            model.ErrorResponse{Error: "db_error", Message: err.Error(), Code: 500})
        return
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        writeJSON(w, http.StatusNotFound,
            model.ErrorResponse{Error: "not_found", Message: "Contact not found", Code: 404})
        return
    }

    writeJSON(w, http.StatusOK, map[string]interface{}{"id": id, "deleted": true})
}