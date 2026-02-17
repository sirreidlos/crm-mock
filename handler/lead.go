package handler

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "crm-mock/model"
)

type LeadHandler struct { db *sql.DB }

func NewLeadHandler(db *sql.DB) *LeadHandler {
    return &LeadHandler{db: db}
}

func (h *LeadHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var l model.Lead
    err := h.db.QueryRow(
        `SELECT id, first_name, last_name, email, company, status, source, created_at
         FROM leads WHERE id = $1`, id,
    ).Scan(&l.ID, &l.FirstName, &l.LastName, &l.Email,
           &l.Company, &l.Status, &l.Source, &l.CreatedAt)

    if err == sql.ErrNoRows {
        writeJSON(w, http.StatusNotFound,
            model.ErrorResponse{Error: "not_found", Message: "Lead not found", Code: 404})
        return
    }
    writeJSON(w, http.StatusOK, l)
}

func (h *LeadHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req model.CreateLeadRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest,
            model.ErrorResponse{Error: "invalid_body", Message: err.Error(), Code: 400})
        return
    }

    _, err := h.db.Exec(
        `INSERT INTO leads (id, first_name, last_name, email, company, source)
         VALUES ($1, $2, $3, $4, $5, $6)`,
        req.ID, req.FirstName, req.LastName, req.Email, req.Company, req.Source,
    )
    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            model.ErrorResponse{Error: "db_error", Message: err.Error(), Code: 500})
        return
    }

    writeJSON(w, http.StatusCreated, req)
}

func (h *LeadHandler) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var req model.UpdateLeadRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest,
            model.ErrorResponse{Error: "invalid_body", Message: err.Error(), Code: 400})
        return
    }

    res, err := h.db.Exec(
        `UPDATE leads SET status = $1, source = COALESCE(NULLIF($2,''), source) WHERE id = $3`,
        req.Status, req.Source, id,
    )
    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            model.ErrorResponse{Error: "db_error", Message: err.Error(), Code: 500})
        return
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        writeJSON(w, http.StatusNotFound,
            model.ErrorResponse{Error: "not_found", Message: "Lead not found", Code: 404})
        return
    }

    writeJSON(w, http.StatusOK, map[string]interface{}{
        "id":      id,
        "status":  req.Status,
        "updated": true,
    })
}

func (h *LeadHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res, err := h.db.Exec(`DELETE FROM leads WHERE id = $1`, id)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            model.ErrorResponse{Error: "db_error", Message: err.Error(), Code: 500})
        return
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        writeJSON(w, http.StatusNotFound,
            model.ErrorResponse{Error: "not_found", Message: "Lead not found", Code: 404})
        return
    }

    writeJSON(w, http.StatusOK, map[string]interface{}{"id": id, "deleted": true})
}